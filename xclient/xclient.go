package xclient

import (
	"context"
	"gorpc"
	"io"
	"reflect"
	"sync"

	"golang.org/x/sync/errgroup"
)

type XClient struct {
	d       Discovery
	mode    SelectMode
	opt     *gorpc.Option
	mu      sync.Mutex // protect following
	clients map[string]*gorpc.Client
}

var _ io.Closer = (*XClient)(nil)

func NewXClient(d Discovery, mode SelectMode, opt *gorpc.Option) *XClient {
	return &XClient{d: d, mode: mode, opt: opt, clients: make(map[string]*gorpc.Client)}
}

func (xc *XClient) Close() error {
	xc.mu.Lock()
	defer xc.mu.Unlock()
	for key, client := range xc.clients {
		// I have no idea how to deal with error, just ignore it.
		_ = client.Close()
		delete(xc.clients, key)
	}
	return nil
}

func (xc *XClient) dial(rpcAddr string) (*gorpc.Client, error) {
	xc.mu.Lock()
	defer xc.mu.Unlock()
	client, ok := xc.clients[rpcAddr]
	if ok && !client.IsAvailable() {
		_ = client.Close()
		delete(xc.clients, rpcAddr)
		client = nil
	}
	if client == nil {
		var err error
		client, err = gorpc.Dial(rpcAddr, xc.opt)
		if err != nil {
			return nil, err
		}
		xc.clients[rpcAddr] = client
	}
	return client, nil
}

func (xc *XClient) call(rpcAddr string, ctx context.Context, serviceMethod string, args, reply any) error {
	client, err := xc.dial(rpcAddr)
	if err != nil {
		return err
	}
	return client.Call(ctx, serviceMethod, args, reply)
}

// Call invokes the named function, waits for it to complete,
// and returns its error status.
// xc will choose a proper server.
func (xc *XClient) Call(ctx context.Context, serviceMethod string, args, reply any) error {
	rpcAddr, err := xc.d.Get(xc.mode)
	if err != nil {
		return err
	}
	return xc.call(rpcAddr, ctx, serviceMethod, args, reply)
}

// Broadcast invokes the named function for every server registered in discovery
// 调用所有服务的 serviceMethod，如果有服务失败,返回其中一个 error
// reply 为任意一个成功的返回值
func (xc *XClient) Broadcast(ctx context.Context, serviceMethod string, args, reply any) error {
	servers, err := xc.d.GetAll()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)

	var mu sync.Mutex         // protect replyDone
	replyDone := reply == nil // if reply is nil, don't need to set value
	for _, rpcAddr := range servers {
		g.Go(func() error {
			var cloneReply any
			if reply != nil {
				cloneReply = reflect.New(reflect.ValueOf(reply).Elem().Type()).Interface()
			}
			if err := xc.call(rpcAddr, ctx, serviceMethod, args, cloneReply); err != nil {
				return err
			}
			mu.Lock()
			if !replyDone {
				reflect.ValueOf(reply).Elem().Set(reflect.ValueOf(cloneReply).Elem())
				replyDone = true
			}
			mu.Unlock()
			return nil
		})
	}

	return g.Wait()
}
