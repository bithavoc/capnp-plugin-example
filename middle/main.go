package main

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	"github.com/bithavoc/procplugin/common"
	"github.com/bithavoc/procplugin/hashes"
	"zombiezen.com/go/capnproto2/rpc"
)

type middlePlugin struct{}

func (p *middlePlugin) Call(c hashes.Plugin_call) error {
	c.Results.SetMessage("message from middle")
	return nil
}

func client(ctx context.Context, c io.ReadWriteCloser) error {
	plugin := &middlePlugin{}
	//pluginServerToClient := hashes.Plugin_ServerToClient(plugin)
	logger.Println("Creating connection")
	// Create a connection that we can use to get the HashFactory.
	conn := rpc.NewConn(rpc.StreamTransport(c))
	defer conn.Close()

	logger.Println("connection open")
	pluginObject := hashes.Plugin_ServerToClient(plugin)
	// Get the "bootstrap" interface.  This is the capability set with
	// rpc.MainInterface on the remote side.
	hf := hashes.PluginRegistry{Client: conn.Bootstrap(ctx)}

	if _, err := hf.Register(ctx, func(p hashes.PluginRegistry_register_Params) error {
		p.SetName("middle")
		p.SetPlugin(pluginObject)
		return nil
	}).Struct(); err != nil {
		logger.Println("failed to register plugin", err.Error())
		panic(err)
	}
	logger.Println("plugin registered in server")

	return nil
}

var logger *log.Logger

func main() {
	debugFile, err := os.Create("middle-debug.log")
	if err != nil {
		panic(err)
	}
	logger = log.New(debugFile, "", log.LstdFlags)
	logger.Println("Debug started")
	time.Sleep(200 * time.Microsecond)
	pipe := common.NewStdStreamJoint(os.Stdin, os.Stdout)
	err = client(context.Background(), pipe)
	if err != nil {
		logger.Println("client error", err.Error())
	}
}
