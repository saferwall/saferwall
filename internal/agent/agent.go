// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package agent

import (
	"context"

	pb "github.com/saferwall/saferwall/internal/agent/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	gRPCMaxMsgSize = 1 * 1024 * 1024 * 1024 // 1GB
)

// Sandbox abstract a sandbox analyzer.
type Sandbox interface {
	Deploy(ctx context.Context, dest string, pkg []byte) (string, error)
	Analyze(config string, binary []byte) (FileScanResult, error)
}

type AgentClient struct {
	client pb.AgentClient
}

// FileScanResult represents the results of a file scan within the sandbox.
type FileScanResult struct {
	TraceLog    []byte
	AgentLog    []byte
	SandboxLog  []byte
	ProcessTree []byte
	Screenshots []*pb.AnalyzeFileReply_Screenshot
	Artifacts   []*pb.AnalyzeFileReply_Artifact
	APIBuffers  []*pb.AnalyzeFileReply_APIBuffer
}

type PingResult struct {
	ServerVersion string
	SysInfo       []byte
}

func New(addr string) (AgentClient, error) {

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return AgentClient{}, err
	}
	client := pb.NewAgentClient(conn)
	return AgentClient{client}, nil
}

func (ac AgentClient) Ping(ctx context.Context) (PingResult, error) {
	r, err := ac.client.Ping(ctx, &emptypb.Empty{})
	if err != nil {
		return PingResult{}, err
	}

	res := PingResult{ServerVersion: r.GetVersion(), SysInfo: r.GetSysinfo()}
	return res, nil
}

func (ac AgentClient) Deploy(ctx context.Context, dest string, pkg []byte) (
	string, error) {
	req := &pb.DeployRequest{Path: dest, Package: pkg}
	r, err := ac.client.Deploy(ctx, req)
	if err != nil {
		return "", err
	}
	return r.GetVersion(), nil
}

func (ac AgentClient) Analyze(ctx context.Context, config, binary []byte) (
	FileScanResult, error) {
	req := &pb.AnalyzeFileRequest{
		Binary: binary,
		Config: config,
	}

	opts := []grpc.CallOption{
		grpc.MaxCallSendMsgSize(gRPCMaxMsgSize),
		grpc.MaxCallRecvMsgSize(gRPCMaxMsgSize),
	}

	r, err := ac.client.Analyze(ctx, req, opts...)
	if err != nil {
		return FileScanResult{}, err
	}

	scanRes := FileScanResult{
		TraceLog:    r.GetAPITrace(),
		AgentLog:    r.GetServerLog(),
		SandboxLog:  r.GetControllerLog(),
		Artifacts:   r.GetArtifacts(),
		Screenshots: r.GetScreenshots(),
		ProcessTree: r.GetProcessTree(),
		APIBuffers:  r.GetAPIBuffers(),
	}

	return scanRes, nil
}
