package main

//go:generate ./gen_proto.sh

import (
	"context"
	"errors"
	"log"
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
	request  proto.Message
	response proto.Message
	ch       chan frog.RpcCall
	err      error
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

func (c *Call) Done() chan frog.RpcCall {
	return c.ch
}

func (c *Call) done() {
	select {
	case c.ch <- c:
	default:
		// dont block
	}
}

type Channel int

func (_ *Channel) Go(method *frog.MethodDesc, ctx context.Context, request proto.Message, response proto.Message) frog.RpcCall {
	call := &Call{
		request,
		response,
		make(chan frog.RpcCall, 5),
		nil,
	}
	var rpcMeth *frog.RpcMethod
	for _, meth := range methods {
		if meth.Descriptor() == method {
			rpcMeth = meth
			break
		}
	}
	if rpcMeth == nil {
		call.err = errors.New("method not found")
		call.done()
		return call
	}

	// async invoke
	go func() {
		time.Sleep(time.Second * 5)
		err := frog.CallMethod(rpcMeth, ctx, request, response)
		call.err = err
		call.done()
	}()

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
	err = stub.Echo(context.TODO(), &request, &response)
	if err != nil {
		log.Fatal(err)
	}
	if request.GetText() != response.GetText() {
		log.Fatal("Text not match:", request.GetText(), "vs", response.GetText())
	}
	log.Println(request.GetText(), response.GetText())
	log.Println("end a sync rpc")

	// make a async rpc
	log.Println("begin a async rpc")
	call := stub.AsyncEcho(context.TODO(), &request, &response)
	log.Println("get a call from async call:", call)

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
		log.Fatal(call.Error())
	}
	if request.GetText() != response.GetText() {
		log.Fatal("Text not match:", request.GetText(), "vs", response.GetText())
	}
	log.Println(request.GetText(), response.GetText())
	log.Println("end a async rpc")
}
