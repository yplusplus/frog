package frog

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	proto "github.com/golang/protobuf/proto"
)

var (
	TypeOfContext = reflect.TypeOf((*context.Context)(nil)).Elem()
)

// RpcCall represents an RPC call
type RpcCall interface {
	Request() proto.Message  // The request to the rpc
	Response() proto.Message // The response from rpc
	Error() error            // After completion, the error status.
	Done() chan struct{}     // Strobes when call is complete.
}

// DefaultCall implements RpcCall interface and is enough to deal with most situations.
type DefaultCall struct {
	request  proto.Message
	response proto.Message
	ch       chan struct{}

	mu      *sync.Mutex // protect following fields
	err     error
	hasDone bool // protect ch from being closed twice
}

func NewDefaultCall(request, response proto.Message) *DefaultCall {
	return &DefaultCall{
		request,
		response,
		make(chan struct{}),
		new(sync.Mutex),
		nil,
		false,
	}
}

func (c *DefaultCall) Request() proto.Message {
	return c.request
}

func (c *DefaultCall) Response() proto.Message {
	return c.response
}

func (c *DefaultCall) Error() error {
	return c.err
}

func (c *DefaultCall) Done() chan struct{} {
	return c.ch
}

func (c *DefaultCall) Close(err error) {
	c.mu.Lock()
	if !c.hasDone {
		c.hasDone = true
		c.err = err
		close(c.ch)
	}
	c.mu.Unlock()
}

// RpcChannel represents a communication line to a Service which can
// be used to call that Service's methods.  The Service may be running
// on another machine.  Normally, you should not call an RpcChannel
// directly, but instead construct a stub Service wrapping it.
// Example:
//   channel := NewMyRpcChannel("remotehost.example.com:1234")
//   stub := NewMyServiceStub(channel)
//	 err := stub.MyMethod(ctx, &request, &response)
//	 call := stub.AsyncMyMethod(ctx, &request, &response)
//	 <-call.Done()
type RpcChannel interface {

	// Go invokes the function asynchronously. It returns the RpcCall structure
	// representing the invocation and will be called by many goroutines simultaneously.
	Go(method *MethodDesc, ctx context.Context, request proto.Message, response proto.Message) RpcCall
}

// RpcMethod represents a rpc method metadata
type RpcMethod struct {
	desc         *MethodDesc
	receiver     reflect.Value
	requestType  reflect.Type
	responseType reflect.Type
	method       reflect.Method
}

// Name returns method name which format is "service.method"
func (meta *RpcMethod) Name() string {
	return meta.desc.service.GetName() + "." + meta.desc.GetName()
}

// Descriptor returns method descriptor
func (meta *RpcMethod) Descriptor() *MethodDesc {
	return meta.desc
}

// NewInput news a input variance
func (meta *RpcMethod) NewRequest() proto.Message {
	return reflect.New(meta.requestType.Elem()).Interface().(proto.Message)
}

// NewInput news a output variance
func (meta *RpcMethod) NewResponse() proto.Message {
	return reflect.New(meta.responseType.Elem()).Interface().(proto.Message)
}

// CallMethod invokes meta.method with given arguments
func CallMethod(meta *RpcMethod, ctx context.Context, request proto.Message, response proto.Message) (err error) {
	args := []reflect.Value{meta.receiver, reflect.ValueOf(ctx), reflect.ValueOf(request), reflect.ValueOf(response)}
	retValues := meta.method.Func.Call(args)
	if retValues[0].Interface() != nil {
		err = retValues[0].Interface().(error)
	}
	return
}

type MethodsRegister func([]*RpcMethod) error

// RegisterService registers all rpc methods in service with given register
// It is called from generated code
func RegisterService(sd *ServiceDesc, service interface{}, register MethodsRegister) (err error) {
	methodsMap := make(map[string]*MethodDesc, len(sd.methods))
	for _, methDesc := range sd.methods {
		methodsMap[methDesc.GetName()] = methDesc
	}

	st := reflect.TypeOf(service)
	rpcMeths := make([]*RpcMethod, 0, len(methodsMap))
	for i := 0; i < st.NumMethod(); i++ {
		method := st.Method(i)
		mname := method.Name
		mtype := method.Type
		methDesc, ok := methodsMap[mname]
		if !ok {
			continue
		}

		// method needs four ins: receiver, context, request, response
		if mtype.NumIn() != 4 {
			panic(fmt.Sprintln("method has wrong number of ins:", mtype.NumIn()))
		}

		recvType := mtype.In(0)
		if recvType.Kind() != reflect.Ptr {
			panic(fmt.Sprintln("method", mname, "receiver type not a pointer:", recvType))
		}

		ctxType := mtype.In(1)
		if ctxType != TypeOfContext {
			panic(fmt.Sprintln("method", mname, "context type is not context.Context:", ctxType))
		}

		// TODO: check if request and response are proto.Message
		requestType := mtype.In(2)
		if requestType.Kind() != reflect.Ptr {
			panic(fmt.Sprintln("method", mname, "request type not a pointer:", requestType))
		}

		responseType := mtype.In(3)
		if responseType.Kind() != reflect.Ptr {
			panic(fmt.Sprintln("method", mname, "response type not a pointer:", responseType))
		}

		rpcMeth := &RpcMethod{
			methDesc,                 // desc
			reflect.ValueOf(service), // receiver
			requestType,              // request type
			responseType,             // response type
			method,                   // method
		}

		rpcMeths = append(rpcMeths, rpcMeth)
	}

	// method's number must be equal len(methodsMap)
	if len(rpcMeths) != len(methodsMap) {
		panic("methods number not match")
	}

	err = register(rpcMeths)
	return
}
