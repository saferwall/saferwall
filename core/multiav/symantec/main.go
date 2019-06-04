// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"net"

	log "github.com/sirupsen/logrus"

	pb "github.com/saferwall/saferwall/core/multiav/symantec/proto"
	"github.com/saferwall/saferwall/pkg/multiav/eset"
	"github.com/saferwall/saferwall/pkg/multiav/symantec"
	"github.com/saferwall/saferwall/pkg/utils"
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

// server is used to implement symantec.SymantecScanner.
type server struct{}

// GetVersion implements eset.EsetAVScanner.
func (s *server) GetVersion(ctx context.Context, in *pb.VersionRequest) (*pb.VersionResponse, error) {
	version, err := eset.GetProgramVersion()
	return &pb.VersionResponse{Version: version}, err
}

// ScanFile implements symantec.SymantecScanner.
func (s *server) ScanFile(ctx context.Context, in *pb.ScanFileRequest) (*pb.ScanResponse, error) {
	res, err := symantec.ScanFile(in.Filepath)
	return &pb.ScanResponse{Infected: res.Infected, Output: res.Output}, err
}

// NewServer creates a new grpc server.
func NewServer(opts ...grpc.ServerOption) *grpc.Server {
	return grpc.NewServer(append(DefaultServerOpts(), opts...)...)
}

// main start a gRPC server and waits for connection.
func main() {

	log.Infoln("Starting Symantec daemon `symcfgd`")
	err := utils.StartCommand("sudo", "/etc/init.d/symcfgd", "start")
	if err != nil {
		log.Fatalf("StartCommand /etc/init.d/symcfgd failed: %v", err)
	}

	log.Infoln("Starting Symantec daemon `rtvscand`")
	err = utils.StartCommand("sudo", "/etc/init.d/rtvscand", "start")
	if err != nil {
		log.Fatalf("StartCommand /etc/init.d/rtvscand failed: %v", err)
	}

	log.Infoln("Starting Symantec gRPC server")

	// create a listener on TCP port 50051
	lis, err := net.Listen("tcp", port)
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}

	// create a gRPC server object
	s := NewServer()

	// attach the SymantecScanner service to the server
	pb.RegisterSymantecScannerServer(s, &server{})

	// register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		grpclog.Fatalf("failed to serve: %v", err)
	}

}
