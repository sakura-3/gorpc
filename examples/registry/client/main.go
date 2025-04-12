package main

import (
	"context"
	"fmt"
	"gorpc/xclient"
	"math/rand/v2"
)

const (
	// client should know registry addr
	registryAddr = "http://localhost:9999/_gorpc_/registry"
)

func main() {
	d := xclient.NewGorpcGeeRegistryDiscovery(registryAddr, 0)

	var (
		xc    = xclient.NewXClient(d, xclient.RandomSelect, nil)
		args  = struct{ Num1, Num2 int }{Num1: rand.Int(), Num2: rand.Int()}
		reply int
	)
	defer func() { _ = xc.Close() }()
	fmt.Printf("client send req %v\n", args)
	if err := xc.Call(context.Background(), "Foo.Sum", args, &reply); err != nil {
		panic(err)
	} else {
		fmt.Printf("%d + %d = %d\n", args.Num1, args.Num2, reply)
	}
}
