// Copyright 2019 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"context"
	"net"

	log "github.com/sirupsen/logrus"

	pb "github.com/saferwall/saferwall/core/multiav/avast/proto"
	"github.com/saferwall/saferwall/pkg/multiav/avast"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"

	// grpc library default is 4MB
	maxMsgSize = 1024 * 1024 * 20
)

// DefaultServerOpts returns the set of default grpc ServerOption's that Tiller requires.
func DefaultServerOpts() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.MaxRecvMsgSize(maxMsgSize),
		grpc.MaxSendMsgSize(maxMsgSize),
	}
}

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
func (s *server) ScanFilePath(ctx context.Context, in *pb.ScanFileRequest) (*pb.ScanResponse, error) {
	res, err := avast.ScanFilePath(in.Filepath)
	return &pb.ScanResponse{Infected: res.Infected, Output: res.Output}, err
}

// ScanFileBinary implements avast.AvastScanner.
func (s *server) ScanFileBinary(ctx context.Context, in *pb.ScanFileBinaryRequest) (*pb.ScanResponse, error) {
	r := bytes.NewReader(in.File)
	res, err := avast.ScanFileBinary(r)
	return &pb.ScanResponse{Infected: res.Infected, Output: res.Output}, err
}

// ActivateLicense implements avast.AvastScanner.
func (s *server) ActivateLicense(ctx context.Context, in *pb.LicenseActivationRequest) (*pb.LicenseActivationResponse, error) {
	r := bytes.NewReader(in.License)
	err := avast.ActivateLicense(r)
	return &pb.LicenseActivationResponse{}, err
}

// GetLicenseStatus implements avast.AvastScanner.
func (s *server) GetLicenseStatus(ctx context.Context, in *pb.LicenseStatusRequest) (*pb.LicenseStatusResponse, error) {
	isExpired, err := avast.IsLicenseExpired()
	return &pb.LicenseStatusResponse{Expired: isExpired}, err
}

// UpdateVPS implements avast.AvastScanner.
func (s *server) UpdateVPS(ctx context.Context, in *pb.UpdateVPSRequest) (*pb.UpdateVPSResponse, error) {
	err := avast.UpdateVPS()
	return &pb.UpdateVPSResponse{}, err
}

// NewServer creates a new grpc server.
func NewServer(opts ...grpc.ServerOption) *grpc.Server {
	return grpc.NewServer(append(DefaultServerOpts(), opts...)...)
}

// main start a gRPC server and waits for connection.
func main() {

	// Start by running avast daemon
	log.Infoln("Starting avast daemon")
	err := avast.StartDaemon()
	if err != nil {
		log.Error(err)
	}

	log.Infoln("Starting avast gRPC server")

	// create a listener on TCP port 50051
	lis, err := net.Listen("tcp", port)
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}

	// create a gRPC server object
	s := NewServer()

	// attach the AvastScanner service to the server
	pb.RegisterAvastScannerServer(s, &server{})

	// register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		grpclog.Fatalf("failed to serve: %v", err)
	}

}
