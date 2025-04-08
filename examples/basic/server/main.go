package main

import (
	"gorpc"
	"log"
	"net"
)

func main() {
	addr := ":8080"
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Printf("start rpc server on %s", addr)
	gorpc.Accept(l)

}
