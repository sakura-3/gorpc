package main

import (
	"context"
	"gorpc"
	"log"
	"sync"
	"time"
)

func main() {
	// client
	client, err := gorpc.Dial("tcp", ":8080")
	if err != nil {
		log.Panic(err)
	}
	defer func() { _ = client.Close() }()

	time.Sleep(time.Second)
	// send request & receive response
	var wg sync.WaitGroup
	for i := range 5 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := struct {
				Num1, Num2 int
			}{i, i * i}
			var reply int

			if err := client.Call(context.Background(), "Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error:", err)
			}
			log.Printf("%d + %d = %d", args.Num1, args.Num2, reply)
		}(i)
	}
	wg.Wait()
}
