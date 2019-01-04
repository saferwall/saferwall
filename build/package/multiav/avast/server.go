// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"context"
	"net"

	pb "github.com/saferwall/saferwall/api/protobuf-spec"
	"github.com/saferwall/saferwall/pkg/multiav/avast"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

// server is used to implement avast.AvastScanner.
type server struct{}

// GetVPSVersion implements avast.AvastScanner.
func (s *server) GetVPSVersion(ctx context.Context, in *pb.VersionRequest) (*pb.VersionResponse, error) {
	version, err := avast.GetVPSVersion()
	return &pb.VersionResponse{Version: version}, err
}

// GetProgramVersion implements avast.AvastScanner.
func (s *server) GetProgramVersion(ctx context.Context, in *pb.VersionRequest) (*pb.VersionResponse, error) {
	version, err := avast.GetProgramVersion()
	return &pb.VersionResponse{Version: version}, err
}

// ScanFile implements avast.AvastScanner.
func (s *server) ScanFile(ctx context.Context, in *pb.ScanFileRequest) (*pb.ScanResponse, error) {
	res, err := avast.ScanFile(in.Filepath)
	return &pb.ScanResponse{Infected: res.Infected, Output: res.Output}, err
}

// ActivateLicense implements avast.AvastScanner.
func (s *server) ActivateLicense(ctx context.Context, in *pb.LicenseActivationRequest) (*pb.LicenseActivationResponse, error) {
	r := bytes.NewReader(in.License)
	err := avast.ActivateLicense(r)
	return &pb.LicenseActivationResponse{}, err
}

// main start a gRPC server and waits for connection.
func main() {

	// create a listener on TCP port 50051
	lis, err := net.Listen("tcp", port)
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}

	// create a gRPC server object
	s := grpc.NewServer()

	// attach the AvastScanner service to the server
	pb.RegisterAvastScannerServer(s, &server{})

	// register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		grpclog.Fatalf("failed to serve: %v", err)
	}
}
