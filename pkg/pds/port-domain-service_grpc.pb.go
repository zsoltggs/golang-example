// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.4
// source: services/port-domain-service/resources/port-domain-service.proto

package pds

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ServiceClient is the client API for Service service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ServiceClient interface {
	UpsertPort(ctx context.Context, in *UpsertPortRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetPortsPaginated(ctx context.Context, in *GetPortsPaginatedRequest, opts ...grpc.CallOption) (*GetPortsPaginatedResponse, error)
	GetPortByID(ctx context.Context, in *GetPortByIDRequest, opts ...grpc.CallOption) (*GetPortByIDResponse, error)
}

type serviceClient struct {
	cc grpc.ClientConnInterface
}

func NewServiceClient(cc grpc.ClientConnInterface) ServiceClient {
	return &serviceClient{cc}
}

func (c *serviceClient) UpsertPort(ctx context.Context, in *UpsertPortRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/pds.Service/UpsertPort", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) GetPortsPaginated(ctx context.Context, in *GetPortsPaginatedRequest, opts ...grpc.CallOption) (*GetPortsPaginatedResponse, error) {
	out := new(GetPortsPaginatedResponse)
	err := c.cc.Invoke(ctx, "/pds.Service/GetPortsPaginated", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) GetPortByID(ctx context.Context, in *GetPortByIDRequest, opts ...grpc.CallOption) (*GetPortByIDResponse, error) {
	out := new(GetPortByIDResponse)
	err := c.cc.Invoke(ctx, "/pds.Service/GetPortByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServiceServer is the server API for Service service.
// All implementations must embed UnimplementedServiceServer
// for forward compatibility
type ServiceServer interface {
	UpsertPort(context.Context, *UpsertPortRequest) (*emptypb.Empty, error)
	GetPortsPaginated(context.Context, *GetPortsPaginatedRequest) (*GetPortsPaginatedResponse, error)
	GetPortByID(context.Context, *GetPortByIDRequest) (*GetPortByIDResponse, error)
	mustEmbedUnimplementedServiceServer()
}

// UnimplementedServiceServer must be embedded to have forward compatible implementations.
type UnimplementedServiceServer struct {
}

func (UnimplementedServiceServer) UpsertPort(context.Context, *UpsertPortRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpsertPort not implemented")
}
func (UnimplementedServiceServer) GetPortsPaginated(context.Context, *GetPortsPaginatedRequest) (*GetPortsPaginatedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPortsPaginated not implemented")
}
func (UnimplementedServiceServer) GetPortByID(context.Context, *GetPortByIDRequest) (*GetPortByIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPortByID not implemented")
}
func (UnimplementedServiceServer) mustEmbedUnimplementedServiceServer() {}

// UnsafeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ServiceServer will
// result in compilation errors.
type UnsafeServiceServer interface {
	mustEmbedUnimplementedServiceServer()
}

func RegisterServiceServer(s grpc.ServiceRegistrar, srv ServiceServer) {
	s.RegisterService(&Service_ServiceDesc, srv)
}

func _Service_UpsertPort_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpsertPortRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).UpsertPort(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pds.Service/UpsertPort",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).UpsertPort(ctx, req.(*UpsertPortRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_GetPortsPaginated_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPortsPaginatedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).GetPortsPaginated(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pds.Service/GetPortsPaginated",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).GetPortsPaginated(ctx, req.(*GetPortsPaginatedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_GetPortByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPortByIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).GetPortByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pds.Service/GetPortByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).GetPortByID(ctx, req.(*GetPortByIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Service_ServiceDesc is the grpc.ServiceDesc for Service service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Service_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pds.Service",
	HandlerType: (*ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpsertPort",
			Handler:    _Service_UpsertPort_Handler,
		},
		{
			MethodName: "GetPortsPaginated",
			Handler:    _Service_GetPortsPaginated_Handler,
		},
		{
			MethodName: "GetPortByID",
			Handler:    _Service_GetPortByID_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "services/port-domain-service/resources/port-domain-service.proto",
}
