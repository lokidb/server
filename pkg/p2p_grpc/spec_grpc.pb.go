// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: spec.proto

package p2p_grcp

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ConfigurationNodeServiceClient is the client API for ConfigurationNodeService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ConfigurationNodeServiceClient interface {
	// Get peer state
	GetState(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*StateResponse, error)
	// Add new peer to cluster
	NewPeer(ctx context.Context, in *NewPeerRequest, opts ...grpc.CallOption) (*NewPeerResponse, error)
}

type configurationNodeServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewConfigurationNodeServiceClient(cc grpc.ClientConnInterface) ConfigurationNodeServiceClient {
	return &configurationNodeServiceClient{cc}
}

func (c *configurationNodeServiceClient) GetState(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*StateResponse, error) {
	out := new(StateResponse)
	err := c.cc.Invoke(ctx, "/lokidb.ConfigurationNodeService/GetState", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configurationNodeServiceClient) NewPeer(ctx context.Context, in *NewPeerRequest, opts ...grpc.CallOption) (*NewPeerResponse, error) {
	out := new(NewPeerResponse)
	err := c.cc.Invoke(ctx, "/lokidb.ConfigurationNodeService/NewPeer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ConfigurationNodeServiceServer is the server API for ConfigurationNodeService service.
// All implementations must embed UnimplementedConfigurationNodeServiceServer
// for forward compatibility
type ConfigurationNodeServiceServer interface {
	// Get peer state
	GetState(context.Context, *empty.Empty) (*StateResponse, error)
	// Add new peer to cluster
	NewPeer(context.Context, *NewPeerRequest) (*NewPeerResponse, error)
	mustEmbedUnimplementedConfigurationNodeServiceServer()
}

// UnimplementedConfigurationNodeServiceServer must be embedded to have forward compatible implementations.
type UnimplementedConfigurationNodeServiceServer struct {
}

func (UnimplementedConfigurationNodeServiceServer) GetState(context.Context, *empty.Empty) (*StateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetState not implemented")
}
func (UnimplementedConfigurationNodeServiceServer) NewPeer(context.Context, *NewPeerRequest) (*NewPeerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewPeer not implemented")
}
func (UnimplementedConfigurationNodeServiceServer) mustEmbedUnimplementedConfigurationNodeServiceServer() {
}

// UnsafeConfigurationNodeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ConfigurationNodeServiceServer will
// result in compilation errors.
type UnsafeConfigurationNodeServiceServer interface {
	mustEmbedUnimplementedConfigurationNodeServiceServer()
}

func RegisterConfigurationNodeServiceServer(s grpc.ServiceRegistrar, srv ConfigurationNodeServiceServer) {
	s.RegisterService(&ConfigurationNodeService_ServiceDesc, srv)
}

func _ConfigurationNodeService_GetState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigurationNodeServiceServer).GetState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lokidb.ConfigurationNodeService/GetState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigurationNodeServiceServer).GetState(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConfigurationNodeService_NewPeer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewPeerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigurationNodeServiceServer).NewPeer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lokidb.ConfigurationNodeService/NewPeer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigurationNodeServiceServer).NewPeer(ctx, req.(*NewPeerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ConfigurationNodeService_ServiceDesc is the grpc.ServiceDesc for ConfigurationNodeService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ConfigurationNodeService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "lokidb.ConfigurationNodeService",
	HandlerType: (*ConfigurationNodeServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetState",
			Handler:    _ConfigurationNodeService_GetState_Handler,
		},
		{
			MethodName: "NewPeer",
			Handler:    _ConfigurationNodeService_NewPeer_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "spec.proto",
}