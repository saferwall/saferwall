// Copyright 2019 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package client

import (
	"context"
	"github.com/saferwall/saferwall/pkg/grpc/multiav"
	pb "github.com/saferwall/saferwall/pkg/grpc/multiav/sophos/proto"
	"google.golang.org/grpc"
	"log"
)

// ScanFile scans file
func ScanFile(client pb.SophosScannerClient, path string) (multiav.ScanResult, error) {
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
	client := pb.NewSophosScannerClient(conn)

	// ScanFile
	res, err := ScanFile(client, filePath)
	if err != nil {
		log.Fatalf("fail to scanfile: %v", err)
	}
	log.Println(res)
}
