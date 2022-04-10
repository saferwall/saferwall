// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package agent

import (
	"context"

	pb "github.com/saferwall/saferwall/internal/agent/proto"
	"google.golang.org/grpc"
)

type Sandbox interface {
	Deploy(ctx context.Context, dest string, pkg []byte) (string, error)
	Analyze(config string, binary []byte) (FileScanResult, error)
}

type AgentClient struct {
	client pb.AgentClient
}

type FileScanResult struct {
	TraceLog []byte
}

func New(addr string) (AgentClient, error) {
	conn, err := grpc.Dial(addr,
		[]grpc.DialOption{grpc.WithInsecure()}...)
	if err != nil {
		return AgentClient{}, err
	}
	client := pb.NewAgentClient(conn)
	return AgentClient{client}, nil
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

	r, err := ac.client.Analyze(ctx, req)
	if err != nil {
		return FileScanResult{}, err
	}

	scanRes := FileScanResult{
		TraceLog: r.GetApitrace(),
	}

	return scanRes, nil
}
