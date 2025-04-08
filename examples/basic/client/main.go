package main

import (
	"fmt"
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
			args := fmt.Sprintf("gorpc req %d", i)
			var reply string
			if err := client.Call("Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error:", err)
			}
			log.Println("reply:", reply)
		}(i)
	}
	wg.Wait()
}
