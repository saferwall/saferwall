// Copyright 2020 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"context"

	"github.com/saferwall/saferwall/pkg/grpc/multiav"
	pb "github.com/saferwall/saferwall/pkg/grpc/multiav/fsecure/proto"
	"github.com/saferwall/saferwall/pkg/multiav/fsecure"
	"go.uber.org/zap"
	"google.golang.org/grpc/grpclog"
)

type server struct {
	avDbUpdateDate int64
	log            *zap.Logger
}

// ScanFile implements bitdefender.BitdefenderScanner.
func (s *server) ScanFile(ctx context.Context, in *pb.ScanFileRequest) (*pb.ScanResponse, error) {
	s.log.Info("Scanning :", zap.String("filepath", in.Filepath))
	res, err := fsecure.ScanFile(in.Filepath)
	output := ""
	if res.Infected {
		if res.FSE != "" {
			output = res.FSE
		} else {
			output = res.Aquarius
		}
	}
	return &pb.ScanResponse{
		Infected: res.Infected,
		Output:   output,
		Update:   s.avDbUpdateDate}, err
}

// main start a gRPC server and waits for connection.
func main() {

	log := multiav.SetupLogging()
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

	// attach the FSecureScanner service to the server
	pb.RegisterFSecureScannerServer(
		s, &server{avDbUpdateDate: avDbUpdateDate, log: log})

	// register reflection service on gRPC server and serve.
	log.Info("Starting FSecure gRPC server ...")
	err = multiav.Serve(s, lis)
	if err != nil {
		grpclog.Fatalf("failed to serve: %v", err)
	}
}
