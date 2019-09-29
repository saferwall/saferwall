// Copyright 2019 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"github.com/saferwall/saferwall/pkg/grpc/multiav"
	pb "github.com/saferwall/saferwall/pkg/grpc/multiav/bitdefender/proto"
	"github.com/saferwall/saferwall/pkg/multiav/bitdefender"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/grpclog"
)

type server struct {
	avDbUpdateDate int64
}

// GetProgramVersion implements bitdefender.BitdefenderScanner.
func (s *server) GetProgramVersion(ctx context.Context, in *pb.VersionRequest) (*pb.VersionResponse, error) {
	version, err := bitdefender.GetProgramVersion()
	return &pb.VersionResponse{Version: version}, err
}

// ScanFile implements bitdefender.BitdefenderScanner.
func (s *server) ScanFile(ctx context.Context, in *pb.ScanFileRequest) (*pb.ScanResponse, error) {
	res, err := bitdefender.ScanFile(in.Filepath)
	return &pb.ScanResponse{
		Infected: res.Infected,
		Output:   res.Output,
		Update:   s.avDbUpdateDate}, err
}

// main start a gRPC server and waits for connection.
func main() {
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
	pb.RegisterBitdefenderScannerServer(
		s, &server{avDbUpdateDate: avDbUpdateDate})

	// register reflection service on gRPC server and serve.
	log.Infoln("Starting Bitdefender gRPC server ...")
	err = multiav.Serve(s, lis)
	if err != nil {
		grpclog.Fatalf("failed to serve: %v", err)
	}
}
