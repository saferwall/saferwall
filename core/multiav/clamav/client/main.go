// Copyright 2019 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package clamav

import (
	"context"

	log "github.com/sirupsen/logrus"

	pb "github.com/saferwall/saferwall/core/multiav/clamav/proto"
	"google.golang.org/grpc"
)

const (
	address = "clamav-svc:50051"
)

// MultiAVScanResult av result
type MultiAVScanResult struct {
	Output   string `json:"output"`
	Infected bool   `json:"infected"`
}

// GetVerion returns version
func GetVerion(client pb.ClamAVScannerClient) (*pb.VersionResponse, error) {
	versionRequest := &pb.VersionRequest{}
	return client.GetVersion(context.Background(), versionRequest)
}

// ScanFile scans file
func ScanFile(client pb.ClamAVScannerClient, path string) (MultiAVScanResult, error) {
	log.Println("ClamAV Scanning:", path)
	scanFile := &pb.ScanFileRequest{Filepath: path}
	res, err := client.ScanFile(context.Background(), scanFile)
	if err != nil {
		return MultiAVScanResult{}, err
	}

	return MultiAVScanResult{
		Output:   res.Output,
		Infected: res.Infected,
	}, nil
}

// Init connection
func Init() (pb.ClamAVScannerClient, error) {

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	conn, err := grpc.Dial(address, []grpc.DialOption{grpc.WithInsecure()}...)
	if err != nil {
		return nil , err
	}

	return pb.NewClamAVScannerClient(conn), nil
}
