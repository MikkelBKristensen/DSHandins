// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
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

const (
	Consensus_Sync_FullMethodName          = "/Consensus.Consensus/Sync"
	Consensus_Ping_FullMethodName          = "/Consensus.Consensus/Ping"
	Consensus_ConnectStatus_FullMethodName = "/Consensus.Consensus/ConnectStatus"
	Consensus_Election_FullMethodName      = "/Consensus.Consensus/Election"
	Consensus_Leader_FullMethodName        = "/Consensus.Consensus/Leader"
)

// ConsensusClient is the client API for Consensus service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ConsensusClient interface {
	Sync(ctx context.Context, in *ClientBid, opts ...grpc.CallOption) (*Ack, error)
	Ping(ctx context.Context, in *Ack, opts ...grpc.CallOption) (*Ack, error)
	ConnectStatus(ctx context.Context, in *Ack, opts ...grpc.CallOption) (*Ack, error)
	Election(ctx context.Context, in *Command, opts ...grpc.CallOption) (*Empty, error)
	Leader(ctx context.Context, in *Coordinator, opts ...grpc.CallOption) (*Empty, error)
}

type consensusClient struct {
	cc grpc.ClientConnInterface
}

func NewConsensusClient(cc grpc.ClientConnInterface) ConsensusClient {
	return &consensusClient{cc}
}

func (c *consensusClient) Sync(ctx context.Context, in *ClientBid, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := c.cc.Invoke(ctx, Consensus_Sync_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *consensusClient) Ping(ctx context.Context, in *Ack, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := c.cc.Invoke(ctx, Consensus_Ping_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *consensusClient) ConnectStatus(ctx context.Context, in *Ack, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := c.cc.Invoke(ctx, Consensus_ConnectStatus_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *consensusClient) Election(ctx context.Context, in *Command, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, Consensus_Election_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *consensusClient) Leader(ctx context.Context, in *Coordinator, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, Consensus_Leader_FullMethodName, in, out, opts...)
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
	ConnectStatus(context.Context, *Ack) (*Ack, error)
	Election(context.Context, *Command) (*Empty, error)
	Leader(context.Context, *Coordinator) (*Empty, error)
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
func (UnimplementedConsensusServer) ConnectStatus(context.Context, *Ack) (*Ack, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConnectStatus not implemented")
}
func (UnimplementedConsensusServer) Election(context.Context, *Command) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Election not implemented")
}
func (UnimplementedConsensusServer) Leader(context.Context, *Coordinator) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Leader not implemented")
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
		FullMethod: Consensus_Sync_FullMethodName,
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
		FullMethod: Consensus_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConsensusServer).Ping(ctx, req.(*Ack))
	}
	return interceptor(ctx, in, info, handler)
}

func _Consensus_ConnectStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Ack)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConsensusServer).ConnectStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Consensus_ConnectStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConsensusServer).ConnectStatus(ctx, req.(*Ack))
	}
	return interceptor(ctx, in, info, handler)
}

func _Consensus_Election_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Command)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConsensusServer).Election(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Consensus_Election_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConsensusServer).Election(ctx, req.(*Command))
	}
	return interceptor(ctx, in, info, handler)
}

func _Consensus_Leader_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Coordinator)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConsensusServer).Leader(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Consensus_Leader_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConsensusServer).Leader(ctx, req.(*Coordinator))
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
		{
			MethodName: "ConnectStatus",
			Handler:    _Consensus_ConnectStatus_Handler,
		},
		{
			MethodName: "Election",
			Handler:    _Consensus_Election_Handler,
		},
		{
			MethodName: "Leader",
			Handler:    _Consensus_Leader_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "Consensus.proto",
}
