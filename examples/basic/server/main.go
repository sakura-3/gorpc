package main

import (
	"gorpc"
	"log"
	"net"
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

	addr := ":8080"
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Printf("start rpc server on %s", addr)
	gorpc.Accept(l)

}
