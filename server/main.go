package main

import (
	"crypto/sha1"
	"fmt"
	"hash"
	"io"
	"log"
	"os/exec"
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

func server(c io.ReadWriteCloser) error {
	// Create a new locally implemented HashFactory.
	main := hashes.HashFactory_ServerToClient(hashFactory{})

	// Listen for calls, using the HashFactory as the bootstrap interface.
	conn := rpc.NewConn(rpc.StreamTransport(c), rpc.MainInterface(main.Client))

	// Wait for connection to abort.

	fmt.Println("Serving, now waiting")
	err := conn.Wait()
	if err != nil {
		fmt.Println("Serve error", err.Error())
	}
	fmt.Println("Serve finished")
	return nil
}

func main() {
	cmd := exec.Command("../client/client")
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
	fmt.Println("Pluggin stopped")
}
