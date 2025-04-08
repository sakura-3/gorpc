package main

import (
	"gorpc"
	"log"
	"time"
)

func doOtherThings() {
	log.Println("do other things")
	time.Sleep(3 * time.Second)
}

func main() {
	// client
	client, err := gorpc.Dial("tcp", ":8080")
	if err != nil {
		log.Panic(err)
	}
	defer func() { _ = client.Close() }()

	args := "this is a async call args"
	var reply string
	call := client.Go("Foo.Sum", args, &reply)
	// Asynchronously wait call done or timeout
	go func() {
		select {
		case <-call.Done:
			if call.Error != nil {
				log.Printf("err=%v", call.Error.Error())
			} else {
				log.Println("reply=", reply)
			}
		case <-time.After(1 * time.Second):
			log.Println("time out")
		}
	}()

	doOtherThings()
}
