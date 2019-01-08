// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"log"

	pb "github.com/saferwall/saferwall/api/protobuf-spec"
	"github.com/saferwall/saferwall/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

const (
	address = "127.0.0.1:50051"
)

func main() {

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewAviraScannerClient(conn)

	b, err := utils.ReadAll("/home/noteworthy/go/src/github.com/saferwall/saferwall/build/package/multiav/avira/hbedv.key")
	if err != nil {
		log.Fatal(err)
	}
	licenseRequest := &pb.LicenseRequest{
		License: b,
	}
	_, errLicense := client.ActivateLicense(context.Background(), licenseRequest)
	if errLicense != nil {
		grpclog.Fatalf("fail to dial: %v", errLicense)
	}
	log.Printf("License was activated")

	request := &pb.ScanFileRequest{
		Filepath: "/samples/eicar.com.txt",
	}
	response, err := client.ScanFile(context.Background(), request)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			// Error was not a status error
		}
		// Use st.Message() and st.Code()
		log.Fatal(st.Message(), st.Code(), ok)
	}
	log.Printf("Avast Result: %s", response.Output)
}
