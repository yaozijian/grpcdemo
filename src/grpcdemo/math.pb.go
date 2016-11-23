// Code generated by protoc-gen-go.
// source: math.proto
// DO NOT EDIT!

/*
Package main is a generated protocol buffer package.

It is generated from these files:
	math.proto

It has these top-level messages:
	OpRequest
	OpReply
	Operator
*/
package main

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type OpRequest struct {
	A int32 `protobuf:"varint,1,opt,name=a" json:"a,omitempty"`
	B int32 `protobuf:"varint,2,opt,name=b" json:"b,omitempty"`
}

func (m *OpRequest) Reset()                    { *m = OpRequest{} }
func (m *OpRequest) String() string            { return proto.CompactTextString(m) }
func (*OpRequest) ProtoMessage()               {}
func (*OpRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type OpReply struct {
	A int32 `protobuf:"varint,1,opt,name=a" json:"a,omitempty"`
	B int32 `protobuf:"varint,2,opt,name=b" json:"b,omitempty"`
	R int32 `protobuf:"varint,3,opt,name=r" json:"r,omitempty"`
}

func (m *OpReply) Reset()                    { *m = OpReply{} }
func (m *OpReply) String() string            { return proto.CompactTextString(m) }
func (*OpReply) ProtoMessage()               {}
func (*OpReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type Operator struct {
	A int32 `protobuf:"varint,1,opt,name=a" json:"a,omitempty"`
}

func (m *Operator) Reset()                    { *m = Operator{} }
func (m *Operator) String() string            { return proto.CompactTextString(m) }
func (*Operator) ProtoMessage()               {}
func (*Operator) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func init() {
	proto.RegisterType((*OpRequest)(nil), "main.op_request")
	proto.RegisterType((*OpReply)(nil), "main.op_reply")
	proto.RegisterType((*Operator)(nil), "main.operator")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for Math service

type MathClient interface {
	Add(ctx context.Context, in *OpRequest, opts ...grpc.CallOption) (*OpReply, error)
	Addlist(ctx context.Context, opts ...grpc.CallOption) (Math_AddlistClient, error)
	Random(ctx context.Context, in *Operator, opts ...grpc.CallOption) (Math_RandomClient, error)
	Addstream(ctx context.Context, opts ...grpc.CallOption) (Math_AddstreamClient, error)
}

type mathClient struct {
	cc *grpc.ClientConn
}

func NewMathClient(cc *grpc.ClientConn) MathClient {
	return &mathClient{cc}
}

func (c *mathClient) Add(ctx context.Context, in *OpRequest, opts ...grpc.CallOption) (*OpReply, error) {
	out := new(OpReply)
	err := grpc.Invoke(ctx, "/main.math/add", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mathClient) Addlist(ctx context.Context, opts ...grpc.CallOption) (Math_AddlistClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Math_serviceDesc.Streams[0], c.cc, "/main.math/addlist", opts...)
	if err != nil {
		return nil, err
	}
	x := &mathAddlistClient{stream}
	return x, nil
}

type Math_AddlistClient interface {
	Send(*Operator) error
	CloseAndRecv() (*Operator, error)
	grpc.ClientStream
}

type mathAddlistClient struct {
	grpc.ClientStream
}

func (x *mathAddlistClient) Send(m *Operator) error {
	return x.ClientStream.SendMsg(m)
}

func (x *mathAddlistClient) CloseAndRecv() (*Operator, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Operator)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *mathClient) Random(ctx context.Context, in *Operator, opts ...grpc.CallOption) (Math_RandomClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Math_serviceDesc.Streams[1], c.cc, "/main.math/random", opts...)
	if err != nil {
		return nil, err
	}
	x := &mathRandomClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Math_RandomClient interface {
	Recv() (*Operator, error)
	grpc.ClientStream
}

type mathRandomClient struct {
	grpc.ClientStream
}

func (x *mathRandomClient) Recv() (*Operator, error) {
	m := new(Operator)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *mathClient) Addstream(ctx context.Context, opts ...grpc.CallOption) (Math_AddstreamClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Math_serviceDesc.Streams[2], c.cc, "/main.math/addstream", opts...)
	if err != nil {
		return nil, err
	}
	x := &mathAddstreamClient{stream}
	return x, nil
}

type Math_AddstreamClient interface {
	Send(*OpRequest) error
	Recv() (*OpReply, error)
	grpc.ClientStream
}

type mathAddstreamClient struct {
	grpc.ClientStream
}

func (x *mathAddstreamClient) Send(m *OpRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *mathAddstreamClient) Recv() (*OpReply, error) {
	m := new(OpReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Math service

type MathServer interface {
	Add(context.Context, *OpRequest) (*OpReply, error)
	Addlist(Math_AddlistServer) error
	Random(*Operator, Math_RandomServer) error
	Addstream(Math_AddstreamServer) error
}

func RegisterMathServer(s *grpc.Server, srv MathServer) {
	s.RegisterService(&_Math_serviceDesc, srv)
}

func _Math_Add_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OpRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MathServer).Add(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.math/Add",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MathServer).Add(ctx, req.(*OpRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Math_Addlist_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MathServer).Addlist(&mathAddlistServer{stream})
}

type Math_AddlistServer interface {
	SendAndClose(*Operator) error
	Recv() (*Operator, error)
	grpc.ServerStream
}

type mathAddlistServer struct {
	grpc.ServerStream
}

func (x *mathAddlistServer) SendAndClose(m *Operator) error {
	return x.ServerStream.SendMsg(m)
}

func (x *mathAddlistServer) Recv() (*Operator, error) {
	m := new(Operator)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Math_Random_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Operator)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MathServer).Random(m, &mathRandomServer{stream})
}

type Math_RandomServer interface {
	Send(*Operator) error
	grpc.ServerStream
}

type mathRandomServer struct {
	grpc.ServerStream
}

func (x *mathRandomServer) Send(m *Operator) error {
	return x.ServerStream.SendMsg(m)
}

func _Math_Addstream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MathServer).Addstream(&mathAddstreamServer{stream})
}

type Math_AddstreamServer interface {
	Send(*OpReply) error
	Recv() (*OpRequest, error)
	grpc.ServerStream
}

type mathAddstreamServer struct {
	grpc.ServerStream
}

func (x *mathAddstreamServer) Send(m *OpReply) error {
	return x.ServerStream.SendMsg(m)
}

func (x *mathAddstreamServer) Recv() (*OpRequest, error) {
	m := new(OpRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Math_serviceDesc = grpc.ServiceDesc{
	ServiceName: "main.math",
	HandlerType: (*MathServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "add",
			Handler:    _Math_Add_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "addlist",
			Handler:       _Math_Addlist_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "random",
			Handler:       _Math_Random_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "addstream",
			Handler:       _Math_Addstream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: fileDescriptor0,
}

func init() { proto.RegisterFile("math.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 199 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0xca, 0x4d, 0x2c, 0xc9,
	0xd0, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0xc9, 0x4d, 0xcc, 0xcc, 0x53, 0xd2, 0xe0, 0xe2,
	0xca, 0x2f, 0x88, 0x2f, 0x4a, 0x2d, 0x2c, 0x4d, 0x2d, 0x2e, 0x11, 0xe2, 0xe1, 0x62, 0x4c, 0x94,
	0x60, 0x54, 0x60, 0xd4, 0x60, 0x0d, 0x62, 0x4c, 0x04, 0xf1, 0x92, 0x24, 0x98, 0x20, 0xbc, 0x24,
	0x25, 0x13, 0x2e, 0x0e, 0xb0, 0xca, 0x82, 0x9c, 0x4a, 0x7c, 0xea, 0x40, 0xbc, 0x22, 0x09, 0x66,
	0x08, 0xaf, 0x48, 0x49, 0x02, 0xa4, 0x2b, 0xb5, 0x28, 0xb1, 0x24, 0xbf, 0x08, 0x55, 0x97, 0xd1,
	0x61, 0x46, 0x2e, 0x16, 0x90, 0x73, 0x84, 0x34, 0xb9, 0x98, 0x13, 0x53, 0x52, 0x84, 0x04, 0xf4,
	0x40, 0x0e, 0xd2, 0x43, 0xb8, 0x46, 0x8a, 0x0f, 0x49, 0x04, 0x68, 0xab, 0x12, 0x83, 0x90, 0x2e,
	0x17, 0x3b, 0x50, 0x69, 0x4e, 0x26, 0xd0, 0xa9, 0x70, 0x49, 0x88, 0xe1, 0x52, 0x68, 0x7c, 0x25,
	0x06, 0x0d, 0x46, 0x21, 0x1d, 0x2e, 0xb6, 0xa2, 0xc4, 0xbc, 0x94, 0xfc, 0x5c, 0xc2, 0xaa, 0x0d,
	0x18, 0x85, 0x8c, 0xb9, 0x38, 0x81, 0x86, 0x17, 0x97, 0x14, 0xa5, 0x26, 0xe6, 0x12, 0xe3, 0x1a,
	0x0d, 0x46, 0x03, 0xc6, 0x24, 0x36, 0x70, 0x60, 0x1a, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x9f,
	0xf2, 0xec, 0xbc, 0x5a, 0x01, 0x00, 0x00,
}