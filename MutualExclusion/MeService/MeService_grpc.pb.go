// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.4
// source: MeService/MeService.proto

package MeService

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	MeService_ConnectionStatus_FullMethodName = "/MeService.MeService/ConnectionStatus"
	MeService_RequestEntry_FullMethodName     = "/MeService.MeService/RequestEntry"
)

// MeServiceClient is the client API for MeService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MeServiceClient interface {
	ConnectionStatus(ctx context.Context, in *ConnectionMsg, opts ...grpc.CallOption) (*Empty, error)
	RequestEntry(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Message, error)
}

type meServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMeServiceClient(cc grpc.ClientConnInterface) MeServiceClient {
	return &meServiceClient{cc}
}

func (c *meServiceClient) ConnectionStatus(ctx context.Context, in *ConnectionMsg, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, MeService_ConnectionStatus_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *meServiceClient) RequestEntry(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Message, error) {
	out := new(Message)
	err := c.cc.Invoke(ctx, MeService_RequestEntry_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MeServiceServer is the server API for MeService service.
// All implementations must embed UnimplementedMeServiceServer
// for forward compatibility
type MeServiceServer interface {
	ConnectionStatus(context.Context, *ConnectionMsg) (*Empty, error)
	RequestEntry(context.Context, *Message) (*Message, error)
	mustEmbedUnimplementedMeServiceServer()
}

// UnimplementedMeServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMeServiceServer struct {
}

func (UnimplementedMeServiceServer) ConnectionStatus(context.Context, *ConnectionMsg) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConnectionStatus not implemented")
}
func (UnimplementedMeServiceServer) RequestEntry(context.Context, *Message) (*Message, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestEntry not implemented")
}
func (UnimplementedMeServiceServer) mustEmbedUnimplementedMeServiceServer() {}

// UnsafeMeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MeServiceServer will
// result in compilation errors.
type UnsafeMeServiceServer interface {
	mustEmbedUnimplementedMeServiceServer()
}

func RegisterMeServiceServer(s grpc.ServiceRegistrar, srv MeServiceServer) {
	s.RegisterService(&MeService_ServiceDesc, srv)
}

func _MeService_ConnectionStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConnectionMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MeServiceServer).ConnectionStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MeService_ConnectionStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MeServiceServer).ConnectionStatus(ctx, req.(*ConnectionMsg))
	}
	return interceptor(ctx, in, info, handler)
}

func _MeService_RequestEntry_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Message)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MeServiceServer).RequestEntry(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MeService_RequestEntry_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MeServiceServer).RequestEntry(ctx, req.(*Message))
	}
	return interceptor(ctx, in, info, handler)
}

// MeService_ServiceDesc is the grpc.ServiceDesc for MeService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MeService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "MeService.MeService",
	HandlerType: (*MeServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ConnectionStatus",
			Handler:    _MeService_ConnectionStatus_Handler,
		},
		{
			MethodName: "RequestEntry",
			Handler:    _MeService_RequestEntry_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "MeService/MeService.proto",
}