// Copyright 2019 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package avast

import (
	"context"

	"github.com/saferwall/saferwall/core/multiav"
	pb "github.com/saferwall/saferwall/core/multiav/avast/proto"
	"google.golang.org/grpc"
)

const (
	address = "avast-svc:50051"
)

// GetVerion returns version
func GetVerion(client pb.AvastScannerClient) (*pb.VersionResponse, error) {
	versionRequest := &pb.VersionRequest{}
	return client.GetProgramVersion(context.Background(), versionRequest)
}

// ScanFile scans file
func ScanFile(client pb.AvastScannerClient, path string) (multiav.ScanResult, error) {
	scanFile := &pb.ScanFileRequest{Filepath: path}
	res, err := client.ScanFilePath(context.Background(), scanFile)
	if err != nil {
		return multiav.ScanResult{}, err
	}

	return multiav.ScanResult{
		Output:   res.Output,
		Infected: res.Infected,
		Update:   res.Update,
	}, nil
}

// Init connection
func Init() (pb.AvastScannerClient, error) {
	conn, err := grpc.Dial(address, []grpc.DialOption{grpc.WithInsecure()}...)
	if err != nil {
		return nil, err
	}
	return pb.NewAvastScannerClient(conn), nil
}
