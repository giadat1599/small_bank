// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.24.4
// source: service_small_bank.proto

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

// SmallBankClient is the client API for SmallBank service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SmallBankClient interface {
	CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error)
	LoginUser(ctx context.Context, in *LoginUserRequest, opts ...grpc.CallOption) (*LoginUserResponse, error)
}

type smallBankClient struct {
	cc grpc.ClientConnInterface
}

func NewSmallBankClient(cc grpc.ClientConnInterface) SmallBankClient {
	return &smallBankClient{cc}
}

func (c *smallBankClient) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error) {
	out := new(CreateUserResponse)
	err := c.cc.Invoke(ctx, "/pb.SmallBank/CreateUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *smallBankClient) LoginUser(ctx context.Context, in *LoginUserRequest, opts ...grpc.CallOption) (*LoginUserResponse, error) {
	out := new(LoginUserResponse)
	err := c.cc.Invoke(ctx, "/pb.SmallBank/LoginUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SmallBankServer is the server API for SmallBank service.
// All implementations must embed UnimplementedSmallBankServer
// for forward compatibility
type SmallBankServer interface {
	CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error)
	LoginUser(context.Context, *LoginUserRequest) (*LoginUserResponse, error)
	mustEmbedUnimplementedSmallBankServer()
}

// UnimplementedSmallBankServer must be embedded to have forward compatible implementations.
type UnimplementedSmallBankServer struct {
}

func (UnimplementedSmallBankServer) CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}
func (UnimplementedSmallBankServer) LoginUser(context.Context, *LoginUserRequest) (*LoginUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginUser not implemented")
}
func (UnimplementedSmallBankServer) mustEmbedUnimplementedSmallBankServer() {}

// UnsafeSmallBankServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SmallBankServer will
// result in compilation errors.
type UnsafeSmallBankServer interface {
	mustEmbedUnimplementedSmallBankServer()
}

func RegisterSmallBankServer(s grpc.ServiceRegistrar, srv SmallBankServer) {
	s.RegisterService(&SmallBank_ServiceDesc, srv)
}

func _SmallBank_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SmallBankServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.SmallBank/CreateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SmallBankServer).CreateUser(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SmallBank_LoginUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SmallBankServer).LoginUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.SmallBank/LoginUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SmallBankServer).LoginUser(ctx, req.(*LoginUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SmallBank_ServiceDesc is the grpc.ServiceDesc for SmallBank service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SmallBank_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.SmallBank",
	HandlerType: (*SmallBankServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateUser",
			Handler:    _SmallBank_CreateUser_Handler,
		},
		{
			MethodName: "LoginUser",
			Handler:    _SmallBank_LoginUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service_small_bank.proto",
}
