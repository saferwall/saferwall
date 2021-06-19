// Copyright 2020 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"context"

	"github.com/saferwall/saferwall/pkg/grpc/multiav"
	pb "github.com/saferwall/saferwall/pkg/grpc/multiav/avira/proto"
	"github.com/saferwall/saferwall/pkg/multiav/avira"
	"go.uber.org/zap"
	"google.golang.org/grpc/grpclog"
)

type server struct {
	avDbUpdateDate int64
	log            *zap.Logger
}

// ScanFile implements avira.AviraScanner.
func (s *server) ScanFile(ctx context.Context, in *pb.ScanFileRequest) (
	*pb.ScanResponse, error) {
	s.log.Info("Scanning :", zap.String("filepath", in.Filepath))
	res, err := avira.ScanFile(in.Filepath)
	return &pb.ScanResponse{
		Infected: res.Infected,
		Output:   res.Output,
		Update:   s.avDbUpdateDate}, err
}

// ActivateLicense implements avira.AviraScanner.
func (s *server) ActivateLicense(ctx context.Context, in *pb.LicenseRequest) (
	*pb.LicenseResponse, error) {
	r := bytes.NewReader(in.License)
	_, err := avira.ActivateLicense(r)
	return &pb.LicenseResponse{}, err
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

	// attach the AviraScanner service to the server
	pb.RegisterAviraScannerServer(
		s, &server{avDbUpdateDate: avDbUpdateDate, log: log})

	// register reflection service on gRPC server and serve.
	log.Info("Starting Avira gRPC server ...")
	err = multiav.Serve(s, lis)
	if err != nil {
		grpclog.Fatalf("failed to serve: %v", err)
	}
}
