package main

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

type Channel int

func (_ *Channel) Go(method *frog.MethodDesc, ctx context.Context, request proto.Message, response proto.Message) frog.RpcCall {
	call := frog.NewDefaultCall(request, response)

	var rpcMeth *frog.RpcMethod
	for _, meth := range methods {
		if meth.Descriptor() == method {
			rpcMeth = meth
			break
		}
	}
	if rpcMeth == nil {
		call.Close(errors.New("method not found"))
		return call
	}

	// async invoke
	go func() {
		time.Sleep(time.Second * 5)
		err := frog.CallMethod(rpcMeth, ctx, request, response)
		call.Close(err)
	}()

	if _, ok := ctx.Deadline(); ok {
		go func() {
			select {
			case <-call.Done():
				// already close, do nothing
			case <-ctx.Done():
				call.Close(errors.New("request timeout"))
			}
		}()
	} else {
		go func() {
			t := time.NewTimer(time.Second)
			select {
			case <-call.Done():
				// already closed, stop the timer
				t.Stop()
			case <-t.C:
				call.Close(errors.New("request timeout"))
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
