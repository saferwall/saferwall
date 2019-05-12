// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package clamav

import (
	"context"
	log "github.com/sirupsen/logrus"

	pb "github.com/saferwall/saferwall/core/multiav/clamav/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

const (
	address = "clamav-svc:50051"
	// address = "192.168.99.100:30051"
)

// MultiAVScanResult av result
type MultiAVScanResult struct {
	Output   string `json:"output"`
	Infected bool   `json:"infected"`
}

func checkgRPCErr(e error) {
	if e != nil {
		st, ok := status.FromError(e)
		if ok {
			log.Errorln(st.Message())
		}
	}
}

// GetVerion returns version
func GetVerion(client pb.ClamAVScannerClient) (*pb.VersionResponse, error) {
	versionRequest := &pb.VersionRequest{}
	return client.GetVersion(context.Background(), versionRequest)
}

// ScanFile scans file
func ScanFile(client pb.ClamAVScannerClient, path string) (MultiAVScanResult, error) {
	log.Println("Scanning:", path)
	scanFile := &pb.ScanFileRequest{Filepath: "eicar"}
	res, err := client.ScanFile(context.Background(), scanFile)
	checkgRPCErr(err)
	if err != nil {
		grpclog.Fatalf("fail to scan file: %v", err)
		return MultiAVScanResult{}, err
	}

	return MultiAVScanResult{
		Output:   res.Output,
		Infected: res.Infected,
	}, nil
}

// Init connection
func Init() pb.ClamAVScannerClient {

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	conn, err := grpc.Dial(address, []grpc.DialOption{grpc.WithInsecure()}...)
	checkgRPCErr(err)

	return pb.NewClamAVScannerClient(conn)
}
