// Copyright 2021 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"context"

	"github.com/saferwall/saferwall/pkg/grpc/multiav"
	pb "github.com/saferwall/saferwall/pkg/grpc/multiav/clamav/proto"
	"github.com/saferwall/saferwall/pkg/multiav/clamav"
	"github.com/saferwall/saferwall/pkg/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc/grpclog"
)

type server struct {
	avDbUpdateDate int64
	log            *zap.Logger
}

// GetVersion implements clamav.ClamAVScanner.
func (s *server) GetVersion(ctx context.Context, in *pb.VersionRequest) (
	*pb.VersionResponse, error) {
	version, err := clamav.GetVersion()
	return &pb.VersionResponse{Version: version}, err
}

// ScanFile implements clamav.ClamAVScanner.
func (s *server) ScanFile(ctx context.Context, in *pb.ScanFileRequest) (
	*pb.ScanResponse, error) {
	s.log.Info("Scanning :", zap.String("filepath", in.Filepath))
	res, err := clamav.ScanFile(in.Filepath)
	return &pb.ScanResponse{
		Infected: res.Infected,
		Output:   res.Output,
		Update:   s.avDbUpdateDate}, err
}

// main start a gRPC server and waits for connection.
func main() {

	log := multiav.SetupLogging()

	// start clamav daemon
	log.Info("Starting clamav daemon ...")
	_, err := utils.ExecCommand("clamd")
	if err != nil {
		grpclog.Fatalf("failed to start clamav daemon: %v", err)
	}

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

	// attach the AviraScanner service to the server
	pb.RegisterClamAVScannerServer(
		s, &server{avDbUpdateDate: avDbUpdateDate, log: log})

	// register reflection service on gRPC server and serve.
	log.Info("Starting ClamAV gRPC server ...")
	err = multiav.Serve(s, lis)
	if err != nil {
		grpclog.Fatalf("failed to serve: %v", err)
	}
}
