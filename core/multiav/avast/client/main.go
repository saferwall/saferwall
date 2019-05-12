// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package avast

import (
	"context"
	log "github.com/sirupsen/logrus"

	pb "github.com/saferwall/saferwall/core/multiav/avast/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

const (
	// address = "avast-svc:50051"
	address = "172.17.0.2:50051"
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
func GetVerion(client pb.AvastScannerClient) (*pb.VersionResponse, error) {
	versionRequest := &pb.VersionRequest{}
	return client.GetProgramVersion(context.Background(), versionRequest)
}

// ScanFile scans file
func ScanFile(client pb.AvastScannerClient, path string) (MultiAVScanResult, error) {
	log.Println("Scanning:", path)
	scanFile := &pb.ScanFileRequest{Filepath: "/samples/eicar"}
	res, err := client.ScanFilePath(context.Background(), scanFile)
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
func Init() pb.AvastScannerClient {

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	conn, err := grpc.Dial(address, []grpc.DialOption{grpc.WithInsecure()}...)
	checkgRPCErr(err)

	return pb.NewAvastScannerClient(conn)
}
