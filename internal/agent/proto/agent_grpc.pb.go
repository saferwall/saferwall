// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.3
// source: agent.proto

package agent

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

// AgentClient is the client API for Agent service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AgentClient interface {
	// Deploy installs all the malware sandbox component files.
	// This include the dll to be injected, the driver, the loader, etc ...
	Deploy(ctx context.Context, in *DeployRequest, opts ...grpc.CallOption) (*DeployReply, error)
	// Analyze executes the sample inside the virtual machine and monitor its
	// behavior.
	Analyze(ctx context.Context, in *AnalyzeFileRequest, opts ...grpc.CallOption) (*AnalyzeFileReply, error)
}

type agentClient struct {
	cc grpc.ClientConnInterface
}

func NewAgentClient(cc grpc.ClientConnInterface) AgentClient {
	return &agentClient{cc}
}

func (c *agentClient) Deploy(ctx context.Context, in *DeployRequest, opts ...grpc.CallOption) (*DeployReply, error) {
	out := new(DeployReply)
	err := c.cc.Invoke(ctx, "/sandbox.Agent/Deploy", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentClient) Analyze(ctx context.Context, in *AnalyzeFileRequest, opts ...grpc.CallOption) (*AnalyzeFileReply, error) {
	out := new(AnalyzeFileReply)
	err := c.cc.Invoke(ctx, "/sandbox.Agent/Analyze", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AgentServer is the server API for Agent service.
// All implementations must embed UnimplementedAgentServer
// for forward compatibility
type AgentServer interface {
	// Deploy installs all the malware sandbox component files.
	// This include the dll to be injected, the driver, the loader, etc ...
	Deploy(context.Context, *DeployRequest) (*DeployReply, error)
	// Analyze executes the sample inside the virtual machine and monitor its
	// behavior.
	Analyze(context.Context, *AnalyzeFileRequest) (*AnalyzeFileReply, error)
	mustEmbedUnimplementedAgentServer()
}

// UnimplementedAgentServer must be embedded to have forward compatible implementations.
type UnimplementedAgentServer struct {
}

func (UnimplementedAgentServer) Deploy(context.Context, *DeployRequest) (*DeployReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Deploy not implemented")
}
func (UnimplementedAgentServer) Analyze(context.Context, *AnalyzeFileRequest) (*AnalyzeFileReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Analyze not implemented")
}
func (UnimplementedAgentServer) mustEmbedUnimplementedAgentServer() {}

// UnsafeAgentServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AgentServer will
// result in compilation errors.
type UnsafeAgentServer interface {
	mustEmbedUnimplementedAgentServer()
}

func RegisterAgentServer(s grpc.ServiceRegistrar, srv AgentServer) {
	s.RegisterService(&Agent_ServiceDesc, srv)
}

func _Agent_Deploy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeployRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServer).Deploy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sandbox.Agent/Deploy",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServer).Deploy(ctx, req.(*DeployRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Agent_Analyze_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AnalyzeFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServer).Analyze(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sandbox.Agent/Analyze",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServer).Analyze(ctx, req.(*AnalyzeFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Agent_ServiceDesc is the grpc.ServiceDesc for Agent service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Agent_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sandbox.Agent",
	HandlerType: (*AgentServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Deploy",
			Handler:    _Agent_Deploy_Handler,
		},
		{
			MethodName: "Analyze",
			Handler:    _Agent_Analyze_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "agent.proto",
}
