// Code generated by protoc-gen-go. DO NOT EDIT.
// source: echo.proto

/*
Package main is a generated protocol buffer package.

It is generated from these files:
	echo.proto
	echo_service.proto

It has these top-level messages:
	ProtoEchoRequest
	ProtoEchoResponse
*/
package main

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type ProtoEchoRequest struct {
	Text             *string `protobuf:"bytes,1,opt,name=text" json:"text,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *ProtoEchoRequest) Reset()                    { *m = ProtoEchoRequest{} }
func (m *ProtoEchoRequest) String() string            { return proto.CompactTextString(m) }
func (*ProtoEchoRequest) ProtoMessage()               {}
func (*ProtoEchoRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ProtoEchoRequest) GetText() string {
	if m != nil && m.Text != nil {
		return *m.Text
	}
	return ""
}

type ProtoEchoResponse struct {
	Text             *string `protobuf:"bytes,1,opt,name=text" json:"text,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *ProtoEchoResponse) Reset()                    { *m = ProtoEchoResponse{} }
func (m *ProtoEchoResponse) String() string            { return proto.CompactTextString(m) }
func (*ProtoEchoResponse) ProtoMessage()               {}
func (*ProtoEchoResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ProtoEchoResponse) GetText() string {
	if m != nil && m.Text != nil {
		return *m.Text
	}
	return ""
}

func init() {
	proto.RegisterType((*ProtoEchoRequest)(nil), "ProtoEchoRequest")
	proto.RegisterType((*ProtoEchoResponse)(nil), "ProtoEchoResponse")
}

func init() { proto.RegisterFile("echo.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 82 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4a, 0x4d, 0xce, 0xc8,
	0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x57, 0x52, 0xe0, 0x12, 0x08, 0x00, 0x31, 0x5c, 0x93, 0x33,
	0xf2, 0x83, 0x52, 0x0b, 0x4b, 0x53, 0x8b, 0x4b, 0x84, 0x78, 0xb8, 0x58, 0x4a, 0x52, 0x2b, 0x4a,
	0x24, 0x18, 0x15, 0x18, 0x35, 0x38, 0x95, 0x14, 0xb9, 0x04, 0x91, 0x54, 0x14, 0x17, 0xe4, 0xe7,
	0x15, 0xa7, 0xa2, 0x2a, 0x01, 0x04, 0x00, 0x00, 0xff, 0xff, 0x4f, 0x3d, 0x63, 0x94, 0x51, 0x00,
	0x00, 0x00,
}
