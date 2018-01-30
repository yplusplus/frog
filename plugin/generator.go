// Package frog outputs rpc service descriptions in Go code.
// It runs as a plugin for the Go protocol buffer compiler plugin.
// It is linked in to protoc-gen-go.
package plugin

import (
	"fmt"
	"path"
	"strconv"
	"strings"

	desc "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/generator"
	"github.com/yplusplus/frog"
)

// Paths for packages used by code generated in this file,
// relative to the import_prefix of the generator.Generator.
const (
	frogPkgPath = "github.com/yplusplus/frog"
)

func init() {
	generator.RegisterPlugin(new(frogGen))
}

// frogGen is an implementation of the Go protocol buffer compiler's
// plugin architecture.  It generates bindings for frog support.
type frogGen struct {
	gen *generator.Generator
}

// Name returns the name of this plugin, "frog".
func (g *frogGen) Name() string {
	return "frog"
}

// The names for packages imported in the generated code.
// They may vary from the final path component of the import path
// if the name is used by other packages.
var (
	frogPkg string
)

// Init initializes the plugin.
func (g *frogGen) Init(gen *generator.Generator) {
	g.gen = gen
	frogPkg = generator.RegisterUniquePackageName("frog", nil)
}

// Given a type name defined in a .proto, return its object.
// Also record that we're using it, to guarantee the associated import.
func (g *frogGen) objectNamed(name string) generator.Object {
	g.gen.RecordTypeUse(name)
	return g.gen.ObjectNamed(name)
}

// Given a type name defined in a .proto, return its name as we will print it.
func (g *frogGen) typeName(str string) string {
	return g.gen.TypeName(g.objectNamed(str))
}

// P forwards to g.gen.P.
func (g *frogGen) P(args ...interface{}) { g.gen.P(args...) }

// Generate generates code for the services in the given file.
func (g *frogGen) Generate(file *generator.FileDescriptor) {

	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}

	g.P("// Reference imports to suppress errors if they are not otherwise used.")
	g.P("var _ context.Context")
	g.P()

	for i, service := range file.FileDescriptorProto.Service {
		g.generateService(file, service, i)
	}

	// generate ServiceDescriptor
	g.P("var (")
	for _, service := range file.FileDescriptorProto.Service {
		servDescName := frog.GenerateServiceDescName(service.GetName())
		g.P(servDescName, " *"+frogPkg+".ServiceDesc")
	}
	g.P(")")
	g.P()

	g.gen.PrintComments("Register and Bind descriptor")
	g.P("func init() { ")
	g.P("frog.GenerateServiceDesc(", file.VarName(), ")")

	for _, service := range file.FileDescriptorProto.Service {
		servDescName := frog.GenerateServiceDescName(service.GetName())
		g.P(servDescName, " = frog.ServiceDescriptor(\""+servDescName+"\")")
	}

	g.P("}")
	g.P()
}

// GenerateImports generates the import declaration for this file.
func (g *frogGen) GenerateImports(file *generator.FileDescriptor) {
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}
	g.P("import (")
	g.P("\"context\"")
	g.P(frogPkg, " ", strconv.Quote(path.Join(g.gen.ImportPrefix, frogPkgPath)))
	g.P(")")
	g.P()
}

func unexport(s string) string { return strings.ToLower(s[:1]) + s[1:] }

// generateService generates all the code for the named service.
func (g *frogGen) generateService(file *generator.FileDescriptor, service *desc.ServiceDescriptorProto, index int) {
	path := fmt.Sprintf("6,%d", index) // 6 means service.
	_ = path

	origServName := service.GetName()
	fullServName := origServName
	if pkg := file.GetPackage(); pkg != "" {
		fullServName = pkg + "." + fullServName
	}
	servName := generator.CamelCase(origServName)

	g.P()
	g.P("// Stub for ", servName, " service")
	g.P()

	// Stub structure.
	g.P("type ", unexport(servName), "Stub struct {")
	g.P("channel ", frogPkg, ".RpcChannel")
	g.P("}")
	g.P()

	// NewStub factory.
	g.P("func New", servName, "Stub (channel ", frogPkg, ".RpcChannel) *", unexport(servName), "Stub {")
	g.P("return &", unexport(servName), "Stub{channel}")
	g.P("}")
	g.P()

	// Stub method implementations.
	for index, method := range service.Method {
		if !method.GetServerStreaming() && !method.GetClientStreaming() {
		} else {
			g.gen.Fail("not support streaming method")
		}
		g.generateStubMethod(servName, fullServName, method, index)
	}

	// Stub call method and go method
	g.P("func (stub *", unexport(servName), "Stub) Call(method *"+frogPkg+".MethodDesc, ctx context.Context, in proto.Message, out proto.Message) error {")
	g.P("call := stub.channel.Go(method, ctx, in, out)")
	g.P("<-call.Done()")
	g.P("return call.Error()")
	g.P("}")
	g.P()

	g.P("func (stub *", unexport(servName), "Stub) Go(method *"+frogPkg+".MethodDesc, ctx context.Context, in proto.Message, out proto.Message)  ", frogPkg, ".RpcCall {")
	g.P("return stub.channel.Go(method, ctx, in, out)")
	g.P("}")
	g.P()

	// Server interface.
	g.P("type ", servName, " interface {")
	for i, method := range service.Method {
		g.gen.PrintComments(fmt.Sprintf("%s,2,%d", path, i)) // 2 means method in a service.
		g.P(g.generateServerSignature(servName, method))
	}
	g.P("}")
	g.P()

	// generate registerXxxService
	// func RegisterEchoService(service EchoService, register frog.MethodsRegister) error {
	//     return frog.RegisterService(EchoService_ServiceDesc, service)
	// }
	g.P(fmt.Sprintf("func Register%s(service %s, register %s.MethodsRegister) error {", servName, servName, frogPkg))
	g.P(fmt.Sprintf("return %s.RegisterService(%s, service, register)", frogPkg, frog.GenerateServiceDescName(servName)))
	g.P("}")
}

// generateStubSignature returns the stub signature for a method.
func (g *frogGen) generateStubSignature(servName string, method *desc.MethodDescriptorProto) string {
	origMethName := method.GetName()
	methName := generator.CamelCase(origMethName)
	inType := g.typeName(method.GetInputType())
	outType := g.typeName(method.GetOutputType())
	return fmt.Sprintf("%s(ctx context.Context, in *%s, out *%s) error", methName, inType, outType)
}

func (g *frogGen) generateStubAsyncSignature(servName string, method *desc.MethodDescriptorProto) string {
	origMethName := method.GetName()
	methName := generator.CamelCase(origMethName)
	inType := g.typeName(method.GetInputType())
	outType := g.typeName(method.GetOutputType())
	return fmt.Sprintf("Async%s(ctx context.Context, in *%s, out *%s) %s.RpcCall", methName, inType, outType, frogPkg)
}

func (g *frogGen) generateStubMethod(servName, fullServName string, method *desc.MethodDescriptorProto, index int) {
	if method.GetServerStreaming() || method.GetClientStreaming() {
		g.gen.Fail("not support streaming method")
	}

	servDescName := frog.GenerateServiceDescName(servName)
	// sync method
	g.P("func (stub *", unexport(servName), "Stub) ", g.generateStubSignature(servName, method), "{")
	g.P("err := stub.Call("+servDescName+".Method(", index, "), ctx, in, out)")
	g.P("return err")
	g.P("}")
	g.P()

	// async method
	g.P("func (stub *", unexport(servName), "Stub) ", g.generateStubAsyncSignature(servName, method), "{")
	g.P("return stub.Go("+servDescName+".Method(", index, "), ctx, in, out)")
	g.P("}")
	g.P()
}

// generateServerSignature returns the server-side signature for a method.
func (g *frogGen) generateServerSignature(servName string, method *desc.MethodDescriptorProto) string {
	origMethName := method.GetName()
	methName := generator.CamelCase(origMethName)
	inType := g.typeName(method.GetInputType())
	outType := g.typeName(method.GetOutputType())
	return methName + "(context.Context, *" + inType + ", *" + outType + ") error"
}
