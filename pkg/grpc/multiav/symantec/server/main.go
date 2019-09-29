// Copyright 2019 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"github.com/saferwall/saferwall/pkg/grpc/multiav"
	log "github.com/sirupsen/logrus"
	pb "github.com/saferwall/saferwall/pkg/grpc/multiav/symantec/proto"
	"github.com/saferwall/saferwall/pkg/multiav/symantec"
	"github.com/saferwall/saferwall/pkg/utils"
	"google.golang.org/grpc/grpclog"

)

const (
	symcfgd = "/opt/Symantec/symantec_antivirus/symcfgd"
	rtvscand = "/opt/Symantec/symantec_antivirus/rtvscand"
)


// server is used to implement symantec.SymantecScanner.
type server struct {
	avDbUpdateDate int64
}

// GetVersion implements eset.EsetAVScanner.
func (s *server) GetVersion(ctx context.Context, in *pb.VersionRequest) (*pb.VersionResponse, error) {
	version, err := symantec.GetProgramVersion()
	return &pb.VersionResponse{Version: version}, err
}

// ScanFile implements symantec.SymantecScanner.
func (s *server) ScanFile(ctx context.Context, in *pb.ScanFileRequest) (*pb.ScanResponse, error) {
	res, err := symantec.ScanFile(in.Filepath)
	return &pb.ScanResponse{
		Infected: res.Infected,
		Output:   res.Output,
		Update:   s.avDbUpdateDate}, err
}


// main start a gRPC server and waits for connection.
func main() {

	log.Infoln("Starting Symantec daemon `symcfgd`")
	out, err := utils.ExecCommand("sudo", symcfgd, "-x")
	if err != nil {
		grpclog.Fatalf("failed to start symcfgd: %v", err)
	}
	log.Infoln(out)

	log.Infoln("Starting Symantec daemon `rtvscand`")
	out, err = utils.ExecCommand("sudo", rtvscand, "-x")
	if err != nil {
		grpclog.Fatalf("failed to start rtvscand: %v", err)
	}
	log.Infoln(out)

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
	pb.RegisterSymantecScannerServer(
		s, &server{avDbUpdateDate: avDbUpdateDate})

	// register reflection service on gRPC server and serve.
	log.Infoln("Starting Symantec gRPC server ...")
	err = multiav.Serve(s, lis)
	if err != nil {
		grpclog.Fatalf("failed to serve: %v", err)
	}
}
