// Copyright 2019 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"github.com/saferwall/saferwall/pkg/grpc/multiav"
	pb "github.com/saferwall/saferwall/pkg/grpc/multiav/bitdefender/proto"
	"google.golang.org/grpc"
	"log"
)

// GetVerion returns version
func GetVerion(client pb.BitdefenderScannerClient) (*pb.VersionResponse, error) {
	versionRequest := &pb.VersionRequest{}
	return client.GetProgramVersion(context.Background(), versionRequest)
}

// ScanFile scans file
func ScanFile(client pb.BitdefenderScannerClient, path string) (multiav.ScanResult, error) {
	scanFile := &pb.ScanFileRequest{Filepath: path}
	res, err := client.ScanFile(context.Background(), scanFile)
	if err != nil {
		return multiav.ScanResult{}, err
	}

	return multiav.ScanResult{
		Output:   res.Output,
		Infected: res.Infected,
		Update:   res.Update,
	}, nil
}

func main() {
	serverAddr, opts, filePath := multiav.ParseFlags()
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewBitdefenderScannerClient(conn)

	// ScanFile
	res, err := ScanFile(client, filePath)
	if err != nil {
		log.Fatalf("fail to scanfile: %v", err)
	}
	log.Println(res)
}
