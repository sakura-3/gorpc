package main

import (
	"gorpc"
	"gorpc/registry"
	"log"
	"net"
	"time"
)

const (
	// server should know registry addr
	registryAddr = "http://localhost:9999/_gorpc_/registry"
)

type Foo int

type Args struct{ Num1, Num2 int }

func (f Foo) Sum(args Args, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}
func main() {
	var foo Foo
	if err := gorpc.Register(&foo); err != nil {
		log.Fatal("register error:", err)
	}

	// find avaliable addr
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Printf("start rpc server on %s", l.Addr().String())

	// helper func provide by registry
	registry.Heartbeat(registryAddr, "tcp@"+l.Addr().String(), 1*time.Minute)

	gorpc.Accept(l)
}
