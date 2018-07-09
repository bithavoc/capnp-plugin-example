package main

import (
	"context"
	"crypto/sha1"
	"fmt"
	"hash"
	"io"
	"log"
	"os/exec"
	"sync"
	"time"

	"github.com/bithavoc/procplugin/common"
	"github.com/bithavoc/procplugin/hashes"
	"zombiezen.com/go/capnproto2/rpc"
)

// hashFactory is a local implementation of HashFactory.
type hashFactory struct{}

func (hf hashFactory) NewSha1(call hashes.HashFactory_newSha1) error {
	fmt.Println("NewSha1 called")
	// Create a new locally implemented Hash capability.
	hs := hashes.Hash_ServerToClient(hashServer{sha1.New()})
	// Notice that methods can return other interfaces.
	return call.Results.SetHash(hs)
}

// hashServer is a local implementation of Hash.
type hashServer struct {
	h hash.Hash
}

func (hs hashServer) Write(call hashes.Hash_write) error {
	data, err := call.Params.Data()
	if err != nil {
		return err
	}
	_, err = hs.h.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (hs hashServer) Sum(call hashes.Hash_sum) error {
	s := hs.h.Sum(nil)
	return call.Results.SetHash(s)
}

type Registry struct {
	plugins map[string]hashes.Plugin
}

func (r *Registry) Register(reg hashes.PluginRegistry_register) error {

	name, err := reg.Params.Name()
	if err != nil {
		return err
	}
	fmt.Println("registry registering", name)
	plugin := reg.Params.Plugin()
	r.plugins[name] = plugin
	callResult, err := plugin.Call(context.Background(), func(p hashes.Plugin_call_Params) error {
		return nil
	}).Struct()
	if err != nil {
		fmt.Println("registry plugin test error", err.Error())
		return err
	}
	msg, _ := callResult.Message()
	fmt.Println("registry registered", r.plugins, msg)
	return nil
}
func (r *Registry) Retrieve(ret hashes.PluginRegistry_retrieve) error {
	fmt.Println("registry retrieve")
	name, err := ret.Params.Name()
	if err != nil {
		return err
	}
	fmt.Println("Retrieve success", name)
	plugin := r.plugins[name]
	ret.Results.SetPlugin(plugin)
	return nil
}

var registry = &Registry{
	plugins: map[string]hashes.Plugin{},
}
var hashServerToClient = hashes.HashFactory_ServerToClient(hashFactory{})
var registryServerToClient = hashes.PluginRegistry_ServerToClient(registry)

func server(c io.ReadWriteCloser) error {
	// Listen for calls, using the HashFactory as the bootstrap interface.
	conn := rpc.NewConn(rpc.StreamTransport(c), rpc.MainInterface(hashServerToClient.Client), rpc.MainInterface(registryServerToClient.Client))

	// Wait for connection to abort.

	fmt.Println("Serving, now waiting")
	err := conn.Wait()
	if err != nil {
		fmt.Println("Serve error", err.Error())
		return err
	}
	fmt.Println("Serve finished")
	return nil
}

var waitGroup = sync.WaitGroup{}

func startPlugin(execPath string) {
	defer func() {
		waitGroup.Done()
	}()
	cmd := exec.Command(execPath)
	inPipeReader, inPipeWriter := io.Pipe()
	outPipeReader, outPipeWriter := io.Pipe()

	cmd.Stdin = inPipeReader
	cmd.Stdout = outPipeWriter

	pipe := common.NewStdStreamJoint(outPipeReader, inPipeWriter)

	go server(pipe)
	time.Sleep(500 * time.Millisecond)
	fmt.Println("Running plugin")
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	go startPlugin("../middle/middle")
	time.Sleep(500 * time.Millisecond)
	go startPlugin("../client/client")
	waitGroup.Add(2)
	waitGroup.Wait()
	fmt.Println("Pluggins stopped")
}
