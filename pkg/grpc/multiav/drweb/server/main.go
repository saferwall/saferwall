// Copyright 2020 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"context"

	"github.com/saferwall/saferwall/pkg/grpc/multiav"
	"github.com/saferwall/saferwall/pkg/utils"
	pb "github.com/saferwall/saferwall/pkg/grpc/multiav/drweb/proto"
	"github.com/saferwall/saferwall/pkg/multiav/drweb"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/grpclog"
)

// server is used to implement drweb.DrWebAVScanner.
type server struct {
	avDbUpdateDate int64
}

// GetVersion implements drweb.DrWebAVScanner.
func (s *server) GetVersion(ctx context.Context, in *pb.VersionRequest) (*pb.VersionResponse, error) {
	version, err := drweb.GetVersion()
	return &pb.VersionResponse{Version: version.CoreEngineVersion}, err
}

// ScanFile implements drweb.DrWebAVScanner.
func (s *server) ScanFile(ctx context.Context, in *pb.ScanFileRequest) (*pb.ScanResponse, error) {
	res, err := drweb.ScanFile(in.Filepath)
	return &pb.ScanResponse{
		Infected: res.Infected,
		Output:   res.Output,
		Update:   s.avDbUpdateDate}, err
}

// main start a gRPC server and waits for connection.
func main() {

	// start DrWeb daemon
	log.Infoln("Starting drweb daemon ...")
	_, err := utils.ExecCommand("sudo", drweb.ConfigD, "-d")
	if err != nil {
		grpclog.Fatalf("failed to start drweb daemon: %v", err)
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

	// attach the DrWebScanner service to the server
	pb.RegisterDrWebScannerServer(
		s, &server{avDbUpdateDate: avDbUpdateDate})

	// register reflection service on gRPC server and serve.
	log.Infoln("Starting DrWeb gRPC server ...")
	err = multiav.Serve(s, lis)
	if err != nil {
		grpclog.Fatalf("failed to serve: %v", err)
	}
}
