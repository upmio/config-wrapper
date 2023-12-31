// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.20.1
// source: app/service/pb/service.proto

package service

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
	ServiceLifecycle_StartService_FullMethodName = "/service.ServiceLifecycle/StartService"
	ServiceLifecycle_StopService_FullMethodName  = "/service.ServiceLifecycle/StopService"
)

// ServiceLifecycleClient is the client API for ServiceLifecycle service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ServiceLifecycleClient interface {
	StartService(ctx context.Context, in *ServiceRequest, opts ...grpc.CallOption) (*ServiceResponse, error)
	StopService(ctx context.Context, in *ServiceRequest, opts ...grpc.CallOption) (*ServiceResponse, error)
}

type serviceLifecycleClient struct {
	cc grpc.ClientConnInterface
}

func NewServiceLifecycleClient(cc grpc.ClientConnInterface) ServiceLifecycleClient {
	return &serviceLifecycleClient{cc}
}

func (c *serviceLifecycleClient) StartService(ctx context.Context, in *ServiceRequest, opts ...grpc.CallOption) (*ServiceResponse, error) {
	out := new(ServiceResponse)
	err := c.cc.Invoke(ctx, ServiceLifecycle_StartService_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceLifecycleClient) StopService(ctx context.Context, in *ServiceRequest, opts ...grpc.CallOption) (*ServiceResponse, error) {
	out := new(ServiceResponse)
	err := c.cc.Invoke(ctx, ServiceLifecycle_StopService_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServiceLifecycleServer is the server API for ServiceLifecycle service.
// All implementations must embed UnimplementedServiceLifecycleServer
// for forward compatibility
type ServiceLifecycleServer interface {
	StartService(context.Context, *ServiceRequest) (*ServiceResponse, error)
	StopService(context.Context, *ServiceRequest) (*ServiceResponse, error)
	mustEmbedUnimplementedServiceLifecycleServer()
}

// UnimplementedServiceLifecycleServer must be embedded to have forward compatible implementations.
type UnimplementedServiceLifecycleServer struct {
}

func (UnimplementedServiceLifecycleServer) StartService(context.Context, *ServiceRequest) (*ServiceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartService not implemented")
}
func (UnimplementedServiceLifecycleServer) StopService(context.Context, *ServiceRequest) (*ServiceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StopService not implemented")
}
func (UnimplementedServiceLifecycleServer) mustEmbedUnimplementedServiceLifecycleServer() {}

// UnsafeServiceLifecycleServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ServiceLifecycleServer will
// result in compilation errors.
type UnsafeServiceLifecycleServer interface {
	mustEmbedUnimplementedServiceLifecycleServer()
}

func RegisterServiceLifecycleServer(s grpc.ServiceRegistrar, srv ServiceLifecycleServer) {
	s.RegisterService(&ServiceLifecycle_ServiceDesc, srv)
}

func _ServiceLifecycle_StartService_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ServiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceLifecycleServer).StartService(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ServiceLifecycle_StartService_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceLifecycleServer).StartService(ctx, req.(*ServiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServiceLifecycle_StopService_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ServiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceLifecycleServer).StopService(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ServiceLifecycle_StopService_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceLifecycleServer).StopService(ctx, req.(*ServiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ServiceLifecycle_ServiceDesc is the grpc.ServiceDesc for ServiceLifecycle service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ServiceLifecycle_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "service.ServiceLifecycle",
	HandlerType: (*ServiceLifecycleServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "StartService",
			Handler:    _ServiceLifecycle_StartService_Handler,
		},
		{
			MethodName: "StopService",
			Handler:    _ServiceLifecycle_StopService_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "app/service/pb/service.proto",
}
