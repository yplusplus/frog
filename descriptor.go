package frog

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"

	proto "github.com/golang/protobuf/proto"
	desc "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

// ServiceDesc wraps ServiceDescriptorProto and provide more functions
type ServiceDesc struct {
	*desc.ServiceDescriptorProto
	methods []*MethodDesc
}

func (sd *ServiceDesc) NumMethod() int {
	return len(sd.methods)
}

func (sd *ServiceDesc) Method(index int) *MethodDesc {
	if index < 0 || index >= sd.NumMethod() {
		panic("")
	}
	return sd.methods[index]
}

// MethodDesc wraps MethodDescriptorProto and provide more functions
type MethodDesc struct {
	*desc.MethodDescriptorProto
	service *ServiceDesc
}

var (
	serviceDescriptors = make(map[string]*ServiceDesc)
)

// GenerateServiceAndMethodDesc called by generated code to generate ServiceDesc and MethodDesc
func GenerateServiceDesc(fileDescriptor []byte) {
	r, err := gzip.NewReader(bytes.NewReader(fileDescriptor))
	if err != nil {
		panic(fmt.Sprintf("failed to open gzip reader: %v", err))
	}
	defer r.Close()

	b, err := ioutil.ReadAll(r)
	if err != nil {
		panic(fmt.Errorf("failed to uncompress descriptor: %v", err))
	}

	fd := new(desc.FileDescriptorProto)
	if err := proto.Unmarshal(b, fd); err != nil {
		panic(fmt.Sprintf("unmarshal FileDescriptorProto failed, err=%s", err))
	}

	for _, service := range fd.Service {
		servDesc := &ServiceDesc{service, make([]*MethodDesc, 0, len(service.Method))}
		servDescName := GenerateServiceDescName(service.GetName())
		serviceDescriptors[servDescName] = servDesc

		for _, method := range service.Method {
			methDesc := &MethodDesc{method, servDesc}
			servDesc.methods = append(servDesc.methods, methDesc)
		}
	}
}

// ServiceDescriptorProto called by generated code
func ServiceDescriptor(servicename string) *ServiceDesc {
	return serviceDescriptors[servicename]
}

