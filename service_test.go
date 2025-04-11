package gorpc

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Foo int

type Args struct{ Num1, Num2 int }

func (f Foo) Sum(args Args, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}

// it's not a exported Method
func (f Foo) sum(args Args, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}

func TestNewService(t *testing.T) {
	var foo Foo
	s := newService(&foo)
	assert.Equal(t, 1, len(s.method), "wrong service Method, expect 1, but got %d", len(s.method))
	mType := s.method["Sum"]
	assert.NotNil(t, mType, "wrong Method, Sum shouldn't nil")
}

func TestMethodType_Call(t *testing.T) {
	var foo Foo
	s := newService(&foo)
	mType := s.method["Sum"]

	argv := mType.newArgv()
	replyv := mType.newReplyv()
	argv.Set(reflect.ValueOf(Args{Num1: 1, Num2: 3}))
	err := s.call(mType, argv, replyv)
	assert.True(t, err == nil && *replyv.Interface().(*int) == 4 && mType.NumCalls() == 1, "failed to call Foo.Sum")
}
