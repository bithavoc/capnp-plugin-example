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

func client(ctx context.Context, c io.ReadWriteCloser) error {

	logger.Println("Creating connection")
	// Create a connection that we can use to get the HashFactory.
	conn := rpc.NewConn(rpc.StreamTransport(c))
	defer conn.Close()

	logger.Println("connection open")
	// Get the "bootstrap" interface.  This is the capability set with
	// rpc.MainInterface on the remote side.
	hf := hashes.HashFactory{Client: conn.Bootstrap(ctx)}

	// Now we can call methods on hf, and they will be sent over c.
	s := hf.NewSha1(ctx, func(p hashes.HashFactory_newSha1_Params) error {
		return nil
	}).Hash()
	logger.Println("hash client open")
	// s refers to a remote Hash.  Method calls are delivered in order.
	s.Write(ctx, func(p hashes.Hash_write_Params) error {
		err := p.SetData([]byte("Hello, "))
		return err
	})

	logger.Println("hello written")
	s.Write(ctx, func(p hashes.Hash_write_Params) error {
		err := p.SetData([]byte("World!"))
		return err
	})
	logger.Println("world written")
	logger.Println("will now call Sum")
	// Get the sum, waiting for the result.
	result, err := s.Sum(ctx, func(p hashes.Hash_sum_Params) error {
		return nil
	}).Struct()
	if err != nil {
		return err
	}

	// Display the result.
	sha1Val, err := result.Hash()
	if err != nil {
		return err
	}

	logger.Printf("sha1: %x\n", sha1Val)
	return nil
}

var logger *log.Logger

func main() {
	debugFile, err := os.Create("debug.log")
	if err != nil {
		panic(err)
	}
	logger = log.New(debugFile, "debug", log.LstdFlags)
	logger.Println("Debug started")
	time.Sleep(200 * time.Microsecond)
	pipe := common.NewStdStreamJoint(os.Stdin, os.Stdout)
	err = client(context.Background(), pipe)
	if err != nil {
		logger.Println("client error", err.Error())
	}
}
