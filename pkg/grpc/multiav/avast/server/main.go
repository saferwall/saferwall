// Copyright 2019 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"context"

	"github.com/saferwall/saferwall/core/multiav"
	pb "github.com/saferwall/saferwall/core/multiav/avast/proto"
	"github.com/saferwall/saferwall/pkg/multiav/avast"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/grpclog"
)

// GetVPSVersion implements avast.AvastScanner.
type server struct {
	avDbUpdateDate int64
}

func (s *server) GetVPSVersion(ctx context.Context, in *pb.VersionRequest) (*pb.VersionResponse, error) {
	version, err := avast.GetVPSVersion()
	return &pb.VersionResponse{Version: version}, err
}

// GetProgramVersion implements avast.AvastScanner.
func (s *server) GetProgramVersion(ctx context.Context, in *pb.VersionRequest) (*pb.VersionResponse, error) {
	version, err := avast.GetProgramVersion()
	return &pb.VersionResponse{Version: version}, err
}

// ScanFilePath implements avast.AvastScanner.
func (s *server) ScanFilePath(ctx context.Context, in *pb.ScanFileRequest) (*pb.ScanResponse, error) {
	res, err := avast.ScanFilePath(in.Filepath)
	return &pb.ScanResponse{
		Infected: res.Infected,
		Output:   res.Output,
		Update:   s.avDbUpdateDate}, err
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

func main() {

	// Start by running avast daemon
	log.Infoln("Starting avast daemon")
	err := avast.StartDaemon()
	if err != nil {
		log.Error(err)
	}

	log.Infoln("Starting avast gRPC server")

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

	// attach the AvastScanner service to the server
	pb.RegisterAvastScannerServer(s, &server{})

	// register reflection service on gRPC server and serve.
	err = multiav.Serve(s, lis)
	if err != nil {
		grpclog.Fatalf("failed to serve: %v", err)
	}
}
