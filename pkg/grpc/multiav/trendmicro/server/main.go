// Copyright 2020 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"context"

	"github.com/saferwall/saferwall/pkg/grpc/multiav"
	"github.com/saferwall/saferwall/pkg/utils"
	pb "github.com/saferwall/saferwall/pkg/grpc/multiav/trendmicro/proto"
	"github.com/saferwall/saferwall/pkg/multiav/trendmicro"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/grpclog"
)

// server is used to implement trendmicro.TrendMicroAVScanner.
type server struct {
	avDbUpdateDate int64
}

// GetVersion implements trendmicro.TrendMicroAVScanner.
func (s *server) GetVersion(ctx context.Context, in *pb.VersionRequest) (*pb.VersionResponse, error) {
	version, err := trendmicro.GetVersion()
	return &pb.VersionResponse{Version: version.EngineVersion}, err
}

// ScanFile implements trendmicro.TrendMicroAVScanner.
func (s *server) ScanFile(ctx context.Context, in *pb.ScanFileRequest) (*pb.ScanResponse, error) {
	log.Printf("Scanning %s", in.Filepath)
	res, err := trendmicro.ScanFile(in.Filepath)
	return &pb.ScanResponse{
		Infected: res.Infected,
		Output:   res.Output,
		Update:   s.avDbUpdateDate}, err
}

// main start a gRPC server and waits for connection.
func main() {

	log.SetFormatter(&log.JSONFormatter{})

	// start TrendMicro daemon
	log.Infoln("Starting TrendMicro service ...")
	_, err := utils.ExecCommand("sudo", "/etc/init.d/splx", "restart")
	if err != nil {
		grpclog.Fatalf("failed to start TrendMicro daemon: %v", err)
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

	// attach the TrendMicroScanner service to the server
	pb.RegisterTrendMicroScannerServer(
		s, &server{avDbUpdateDate: avDbUpdateDate})

	// register reflection service on gRPC server and serve.
	log.Infoln("Starting TrendMciro gRPC server ...")
	err = multiav.Serve(s, lis)
	if err != nil {
		grpclog.Fatalf("failed to serve: %v", err)
	}
}
