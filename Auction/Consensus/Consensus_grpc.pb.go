// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.24.3
// source: Consensus.proto

package Consensus

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

// ConsensusClient is the client API for Consensus service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ConsensusClient interface {
	Sync(ctx context.Context, in *ClientBid, opts ...grpc.CallOption) (*Ack, error)
	Ping(ctx context.Context, in *Ack, opts ...grpc.CallOption) (*Ack, error)
}

type consensusClient struct {
	cc grpc.ClientConnInterface
}

func NewConsensusClient(cc grpc.ClientConnInterface) ConsensusClient {
	return &consensusClient{cc}
}

func (c *consensusClient) Sync(ctx context.Context, in *ClientBid, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := c.cc.Invoke(ctx, "/Consensus.Consensus/Sync", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *consensusClient) Ping(ctx context.Context, in *Ack, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := c.cc.Invoke(ctx, "/Consensus.Consensus/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ConsensusServer is the server API for Consensus service.
// All implementations must embed UnimplementedConsensusServer
// for forward compatibility
type ConsensusServer interface {
	Sync(context.Context, *ClientBid) (*Ack, error)
	Ping(context.Context, *Ack) (*Ack, error)
	mustEmbedUnimplementedConsensusServer()
}

// UnimplementedConsensusServer must be embedded to have forward compatible implementations.
type UnimplementedConsensusServer struct {
}

func (UnimplementedConsensusServer) Sync(context.Context, *ClientBid) (*Ack, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Sync not implemented")
}
func (UnimplementedConsensusServer) Ping(context.Context, *Ack) (*Ack, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedConsensusServer) mustEmbedUnimplementedConsensusServer() {}

// UnsafeConsensusServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ConsensusServer will
// result in compilation errors.
type UnsafeConsensusServer interface {
	mustEmbedUnimplementedConsensusServer()
}

func RegisterConsensusServer(s grpc.ServiceRegistrar, srv ConsensusServer) {
	s.RegisterService(&Consensus_ServiceDesc, srv)
}

func _Consensus_Sync_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClientBid)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConsensusServer).Sync(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Consensus.Consensus/Sync",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConsensusServer).Sync(ctx, req.(*ClientBid))
	}
	return interceptor(ctx, in, info, handler)
}

func _Consensus_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Ack)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConsensusServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Consensus.Consensus/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConsensusServer).Ping(ctx, req.(*Ack))
	}
	return interceptor(ctx, in, info, handler)
}

// Consensus_ServiceDesc is the grpc.ServiceDesc for Consensus service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Consensus_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Consensus.Consensus",
	HandlerType: (*ConsensusServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Sync",
			Handler:    _Consensus_Sync_Handler,
		},
		{
			MethodName: "Ping",
			Handler:    _Consensus_Ping_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "Consensus.proto",
}
