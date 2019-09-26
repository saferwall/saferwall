// Copyright 2019 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package avira

import (
	"context"

	"github.com/saferwall/saferwall/core/multiav"
	pb "github.com/saferwall/saferwall/core/multiav/avira/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	address = "avira-svc:50051"
)

// ScanFile scans file
func ScanFile(client pb.AviraScannerClient, path string) (MultiAVScanResult, error) {
	log.Info("Scanning:", path)
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

// Init connection
func Init() (pb.AviraScannerClient, error) {

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	conn, err := grpc.Dial(address, []grpc.DialOption{grpc.WithInsecure()}...)
	if err != nil {
		return nil, err
	}
	return pb.NewAviraScannerClient(conn), nil
}
