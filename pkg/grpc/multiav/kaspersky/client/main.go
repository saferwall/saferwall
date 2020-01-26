// Copyright 2019 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package client

import (
	"context"
	"github.com/saferwall/saferwall/pkg/grpc/multiav"
	pb "github.com/saferwall/saferwall/pkg/grpc/multiav/kaspersky/proto"
)

// ScanFile scans file
func ScanFile(client pb.KasperskyScannerClient, path string) (multiav.ScanResult, error) {
	scanFileRequest := &pb.ScanFileRequest{Filepath: path}
	res, err := client.ScanFile(context.Background(), scanFileRequest)
	if err != nil {
		return multiav.ScanResult{}, err
	}

	return multiav.ScanResult{
		Output:   res.Output,
		Infected: res.Infected,
		Update:   res.Update,
	}, nil
}