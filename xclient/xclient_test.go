package xclient

import (
	"context"
	"gorpc"
	"net"
	"sync"
	"testing"
	"time"
)

type Foo int

type Args struct{ Num1, Num2 int }

func (f Foo) Sum(args Args, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}

func (f Foo) Sleep(args Args, reply *int) error {
	time.Sleep(time.Second * time.Duration(args.Num1))
	*reply = args.Num1 + args.Num2
	return nil
}

func startServer(addrCh chan string) {
	var foo Foo
	l, _ := net.Listen("tcp", ":0")
	server := gorpc.NewServer()
	_ = server.Register(&foo)
	addrCh <- l.Addr().String()
	server.Accept(l)
}

func call(t *testing.T, addr1, addr2 string) {
	d := NewMultiServerDiscovery([]string{"tcp@" + addr1, "tcp@" + addr2})
	xc := NewXClient(d, RandomSelect, nil)
	defer func() { _ = xc.Close() }()
	// send request & receive response
	var wg sync.WaitGroup
	for i := range 5 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			var reply int
			if err := xc.Call(context.Background(), "Foo.Sum", &Args{Num1: i, Num2: i * i}, &reply); err != nil {
				t.Errorf("call fail:%s", err.Error())
			}
		}(i)
	}
	wg.Wait()
}

func broadcast(t *testing.T, addr1, addr2 string) {
	d := NewMultiServerDiscovery([]string{"tcp@" + addr1, "tcp@" + addr2})
	xc := NewXClient(d, RandomSelect, nil)
	defer func() { _ = xc.Close() }()
	var wg sync.WaitGroup
	for i := range 5 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			var reply int
			if err := xc.Broadcast(context.Background(), "Foo.Sum", &Args{Num1: i, Num2: i * i}, &reply); err != nil {
				t.Errorf("broadcast fail:%s", err.Error())
			}
		}(i)
	}
	wg.Wait()
}

func TestCall(t *testing.T) {
	ch1 := make(chan string)
	ch2 := make(chan string)
	// start two servers
	go startServer(ch1)
	go startServer(ch2)

	addr1 := <-ch1
	addr2 := <-ch2

	time.Sleep(time.Second)
	call(t, addr1, addr2)
}

func TestBroadcast(t *testing.T) {
	ch1 := make(chan string)
	ch2 := make(chan string)
	// start two servers
	go startServer(ch1)
	go startServer(ch2)

	addr1 := <-ch1
	addr2 := <-ch2

	time.Sleep(time.Second)
	broadcast(t, addr1, addr2)
}
