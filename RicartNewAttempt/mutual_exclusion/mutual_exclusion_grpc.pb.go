// mutual_exclusion.proto

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.4
// source: mutual_exclusion/mutual_exclusion.proto

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
	MutualExclusion_RequestEntry_FullMethodName = "/mutual_exclusion.MutualExclusion/RequestEntry"
	MutualExclusion_ReleaseEntry_FullMethodName = "/mutual_exclusion.MutualExclusion/ReleaseEntry"
)

// MutualExclusionClient is the client API for MutualExclusion service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MutualExclusionClient interface {
	RequestEntry(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
	ReleaseEntry(ctx context.Context, in *Release, opts ...grpc.CallOption) (*Response, error)
}

type mutualExclusionClient struct {
	cc grpc.ClientConnInterface
}

func NewMutualExclusionClient(cc grpc.ClientConnInterface) MutualExclusionClient {
	return &mutualExclusionClient{cc}
}

func (c *mutualExclusionClient) RequestEntry(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, MutualExclusion_RequestEntry_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mutualExclusionClient) ReleaseEntry(ctx context.Context, in *Release, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, MutualExclusion_ReleaseEntry_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MutualExclusionServer is the server API for MutualExclusion service.
// All implementations must embed UnimplementedMutualExclusionServer
// for forward compatibility
type MutualExclusionServer interface {
	RequestEntry(context.Context, *Request) (*Response, error)
	ReleaseEntry(context.Context, *Release) (*Response, error)
	mustEmbedUnimplementedMutualExclusionServer()
}

// UnimplementedMutualExclusionServer must be embedded to have forward compatible implementations.
type UnimplementedMutualExclusionServer struct {
}

func (UnimplementedMutualExclusionServer) RequestEntry(context.Context, *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestEntry not implemented")
}
func (UnimplementedMutualExclusionServer) ReleaseEntry(context.Context, *Release) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReleaseEntry not implemented")
}
func (UnimplementedMutualExclusionServer) mustEmbedUnimplementedMutualExclusionServer() {}

// UnsafeMutualExclusionServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MutualExclusionServer will
// result in compilation errors.
type UnsafeMutualExclusionServer interface {
	mustEmbedUnimplementedMutualExclusionServer()
}

func RegisterMutualExclusionServer(s grpc.ServiceRegistrar, srv MutualExclusionServer) {
	s.RegisterService(&MutualExclusion_ServiceDesc, srv)
}

func _MutualExclusion_RequestEntry_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MutualExclusionServer).RequestEntry(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MutualExclusion_RequestEntry_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MutualExclusionServer).RequestEntry(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _MutualExclusion_ReleaseEntry_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Release)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MutualExclusionServer).ReleaseEntry(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MutualExclusion_ReleaseEntry_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MutualExclusionServer).ReleaseEntry(ctx, req.(*Release))
	}
	return interceptor(ctx, in, info, handler)
}

// MutualExclusion_ServiceDesc is the grpc.ServiceDesc for MutualExclusion service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MutualExclusion_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "mutual_exclusion.MutualExclusion",
	HandlerType: (*MutualExclusionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RequestEntry",
			Handler:    _MutualExclusion_RequestEntry_Handler,
		},
		{
			MethodName: "ReleaseEntry",
			Handler:    _MutualExclusion_ReleaseEntry_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mutual_exclusion/mutual_exclusion.proto",
}
