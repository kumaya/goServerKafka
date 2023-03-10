// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: proto/manager/server.proto

package pb

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

// ManagerClient is the client API for Manager service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ManagerClient interface {
	Connect(ctx context.Context, opts ...grpc.CallOption) (Manager_ConnectClient, error)
	Outcome(ctx context.Context, in *OutcomeRequest, opts ...grpc.CallOption) (*OutcomeResponse, error)
}

type managerClient struct {
	cc grpc.ClientConnInterface
}

func NewManagerClient(cc grpc.ClientConnInterface) ManagerClient {
	return &managerClient{cc}
}

func (c *managerClient) Connect(ctx context.Context, opts ...grpc.CallOption) (Manager_ConnectClient, error) {
	stream, err := c.cc.NewStream(ctx, &Manager_ServiceDesc.Streams[0], "/manager.Manager/Connect", opts...)
	if err != nil {
		return nil, err
	}
	x := &managerConnectClient{stream}
	return x, nil
}

type Manager_ConnectClient interface {
	Send(*ClientRequest) error
	Recv() (*ServerResponse, error)
	grpc.ClientStream
}

type managerConnectClient struct {
	grpc.ClientStream
}

func (x *managerConnectClient) Send(m *ClientRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *managerConnectClient) Recv() (*ServerResponse, error) {
	m := new(ServerResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *managerClient) Outcome(ctx context.Context, in *OutcomeRequest, opts ...grpc.CallOption) (*OutcomeResponse, error) {
	out := new(OutcomeResponse)
	err := c.cc.Invoke(ctx, "/manager.Manager/Outcome", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ManagerServer is the server API for Manager service.
// All implementations must embed UnimplementedManagerServer
// for forward compatibility
type ManagerServer interface {
	Connect(Manager_ConnectServer) error
	Outcome(context.Context, *OutcomeRequest) (*OutcomeResponse, error)
	mustEmbedUnimplementedManagerServer()
}

// UnimplementedManagerServer must be embedded to have forward compatible implementations.
type UnimplementedManagerServer struct {
}

func (UnimplementedManagerServer) Connect(Manager_ConnectServer) error {
	return status.Errorf(codes.Unimplemented, "method Connect not implemented")
}
func (UnimplementedManagerServer) Outcome(context.Context, *OutcomeRequest) (*OutcomeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Outcome not implemented")
}
func (UnimplementedManagerServer) mustEmbedUnimplementedManagerServer() {}

// UnsafeManagerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ManagerServer will
// result in compilation errors.
type UnsafeManagerServer interface {
	mustEmbedUnimplementedManagerServer()
}

func RegisterManagerServer(s grpc.ServiceRegistrar, srv ManagerServer) {
	s.RegisterService(&Manager_ServiceDesc, srv)
}

func _Manager_Connect_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ManagerServer).Connect(&managerConnectServer{stream})
}

type Manager_ConnectServer interface {
	Send(*ServerResponse) error
	Recv() (*ClientRequest, error)
	grpc.ServerStream
}

type managerConnectServer struct {
	grpc.ServerStream
}

func (x *managerConnectServer) Send(m *ServerResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *managerConnectServer) Recv() (*ClientRequest, error) {
	m := new(ClientRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Manager_Outcome_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OutcomeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManagerServer).Outcome(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/manager.Manager/Outcome",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManagerServer).Outcome(ctx, req.(*OutcomeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Manager_ServiceDesc is the grpc.ServiceDesc for Manager service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Manager_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "manager.Manager",
	HandlerType: (*ManagerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Outcome",
			Handler:    _Manager_Outcome_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Connect",
			Handler:       _Manager_Connect_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "proto/manager/server.proto",
}
