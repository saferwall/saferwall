// Copyright 2020 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"sync"

	"github.com/saferwall/saferwall/pkg/grpc/multiav"
	pb "github.com/saferwall/saferwall/pkg/grpc/multiav/symantec/proto"
	"github.com/saferwall/saferwall/pkg/multiav/symantec"
	"github.com/saferwall/saferwall/pkg/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc/grpclog"
)

const (
	symcfgd  = "/opt/Symantec/symantec_antivirus/symcfgd"
	rtvscand = "/opt/Symantec/symantec_antivirus/rtvscand"
)

// server is used to implement symantec.SymantecScanner.
type server struct {
	avDbUpdateDate int64
	log            *zap.Logger
	mu             sync.Mutex
}

// GetVersion implements eset.EsetAVScanner.
func (s *server) GetVersion(ctx context.Context, in *pb.VersionRequest) (*pb.VersionResponse, error) {
	version, err := symantec.GetProgramVersion()
	return &pb.VersionResponse{Version: version}, err
}

// ScanFile implements symantec.SymantecScanner.
func (s *server) ScanFile(ctx context.Context, in *pb.ScanFileRequest) (*pb.ScanResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.log.Info("Scanning :", zap.String("filepath", in.Filepath))
	res, err := symantec.ScanFile(in.Filepath)
	return &pb.ScanResponse{
		Infected: res.Infected,
		Output:   res.Output,
		Update:   s.avDbUpdateDate}, err
}

// main start a gRPC server and waits for connection.
func main() {
	log := multiav.SetupLogging()
	log.Info("Starting Symantec daemon `symcfgd`")
	out, err := utils.ExecCommand("sudo", symcfgd, "-x")
	if err != nil {
		grpclog.Fatalf("failed to start symcfgd: %v", err)
	}
	log.Info(out)

	log.Info("Starting Symantec daemon `rtvscand`")
	out, err = utils.ExecCommand("sudo", rtvscand, "-x")
	if err != nil {
		grpclog.Fatalf("failed to start rtvscand: %v", err)
	}
	log.Info(out)

	// create a listener on TCP port 50051
	lis, err := multiav.CreateListener()
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}

	// create a gRPC server object
	s := multiav.NewServer()

	// get the av db update date
	avDbUpdateDate, err := multiav.UpdateDate()
	if err != nil {
		grpclog.Fatalf("failed to read av db update date %v", err)
	}

	// attach the Symanteccanner service to the server
	pb.RegisterSymantecScannerServer(
		s, &server{avDbUpdateDate: avDbUpdateDate, log: log})

	// register reflection service on gRPC server and serve.
	log.Info("Starting Symantec gRPC server ...")
	err = multiav.Serve(s, lis)
	if err != nil {
		grpclog.Fatalf("failed to serve: %v", err)
	}
}
