package main

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	proto "github.com/golang/protobuf/proto"
	"github.com/yplusplus/frog"
)

var (
	methods = make([]*frog.RpcMethod, 0)
	channel frog.RpcChannel
)

func registerMethods(method []*frog.RpcMethod) error {
	methods = append(methods, method...)
	return nil
}

type EchoServiceImpl int

func (impl *EchoServiceImpl) Echo(ctx context.Context, in *ProtoEchoRequest, out *ProtoEchoResponse) error {
	out.Text = proto.String(*in.Text)
	return nil
}

func (impl *EchoServiceImpl) Echo2(ctx context.Context, in *ProtoEchoRequest, out *ProtoEchoResponse) error {
	out.Text = proto.String(*in.Text)
	return nil
}

type Call struct {
	mu       *sync.Mutex
	request  proto.Message
	response proto.Message
	ch       chan struct{}
	err      error
	hasDone  bool
}

func (c *Call) Request() proto.Message {
	return c.request
}

func (c *Call) Response() proto.Message {
	return c.response
}

func (c *Call) Error() error {
	return c.err
}

func (c *Call) Done() chan struct{} {
	return c.ch
}

func (c *Call) done(err error) {
	c.mu.Lock()
	if !c.hasDone {
		c.hasDone = true
		c.err = err
		close(c.ch)
	}
	c.mu.Unlock()
}

type Channel int

func (_ *Channel) Go(method *frog.MethodDesc, ctx context.Context, request proto.Message, response proto.Message) frog.RpcCall {
	call := &Call{
		new(sync.Mutex),
		request,
		response,
		make(chan struct{}),
		nil,
		false,
	}
	var rpcMeth *frog.RpcMethod
	for _, meth := range methods {
		if meth.Descriptor() == method {
			rpcMeth = meth
			break
		}
	}
	if rpcMeth == nil {
		call.done(errors.New("method not found"))
		return call
	}

	// async invoke
	go func() {
		time.Sleep(time.Second * 5)
		err := frog.CallMethod(rpcMeth, ctx, request, response)
		call.done(err)
	}()

	if _, ok := ctx.Deadline(); ok {
		go func() {
			select {
			case <-call.Done():
				// already done, do nothing
			case <-ctx.Done():
				call.done(errors.New("request timeout"))
			}
		}()
	}

	return call
}

func main() {
	// implement service
	impl := new(EchoServiceImpl)
	var err error
	err = RegisterEchoService(impl, registerMethods)
	if err != nil {
		log.Fatal(err)
	}

	// create stub
	channel = new(Channel)
	stub := NewEchoServiceStub(channel)
	var request ProtoEchoRequest
	var response ProtoEchoResponse
	request.Text = proto.String("Hello, world!")

	// make a sync rpc
	log.Println("begin a sync rpc")
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	err = stub.Echo(ctx, &request, &response)
	if err != nil {
		log.Println(err)
	} else if request.GetText() != response.GetText() {
		log.Println("Text not match:", request.GetText(), "vs", response.GetText())
	} else {
		log.Println(request.GetText(), response.GetText())
	}
	log.Println("end a sync rpc")

	// make a async rpc
	log.Println("begin a async rpc")
	call := stub.AsyncEcho(context.TODO(), &request, &response)

	// wating for response
loop:
	for {
		select {
		case <-call.Done():
			log.Println("got a async call's response")
			break loop
		case <-time.Tick(time.Second):
			log.Println("wating for a async call's reponse")
		}
	}
	if call.Error() != nil {
		log.Println(call.Error())
	} else if request.GetText() != response.GetText() {
		log.Println("Text not match:", request.GetText(), "vs", response.GetText())
	} else {
		log.Println(request.GetText(), response.GetText())
	}
	log.Println("end a async rpc")
}
