// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"log"

	"github.com/fatih/color"
	flags "github.com/jessevdk/go-flags"
	pb "github.com/saferwall/saferwall/api/protobuf-spec"
	"github.com/saferwall/saferwall/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

var opts struct {
	Scan            string `short:"s" long:"scan" description:"Scan a file|directory|url"`
	Update          bool   `short:"u" long:"update" description:"Update VPS database"`
	VPSVersion      bool   `short:"v" long:"vps-version" description:"Get VPS version"`
	ProgramVersion  bool   `short:"p" long:"program-version" description:"Get program version"`
	ActivateLicense string `short:"l" long:"license" description:"Activate a license from a license key [license.avastlic]"`
	StatusLicense   bool   `short:"e" long:"status-license" description:"Check expirate of license"`
}

var client pb.AvastScannerClient

const (
	address = "127.0.0.1:50051"
)

func check(e error) {
	if e != nil {
		log.Fatalln(e.Error())
		return
	}
}

func checkgRPCErr(e error) {
	if e != nil {
		st, ok := status.FromError(e)
		if ok {
			log.Println(st.Message())
		}
	}
}

func activateLicense(licensePath string) {
	b, err := utils.ReadAll(licensePath)
	if err != nil {
		log.Fatal(err)
	}
	licenseRequest := &pb.LicenseActivationRequest{
		License: b,
	}
	_, errLicense := client.ActivateLicense(context.Background(), licenseRequest)
	if errLicense != nil {
		grpclog.Fatalf("fail to dial: %v", errLicense)
	}
	log.Printf("License was activated")
}

func getVPSVersion() {
	versionRequest := &pb.VersionRequest{}
	vpsVersion, err := client.GetVPSVersion(context.Background(), versionRequest)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	log.Printf("VPS version: %s", vpsVersion.Version)
}

func getProgramVersion() {
	versionRequest := &pb.VersionRequest{}
	vpsVersion, err := client.GetProgramVersion(context.Background(), versionRequest)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	log.Printf("Program version: %s", vpsVersion.Version)
}

func scanFileBinary(data []byte) {
	request := &pb.ScanFileBinaryRequest{
		File: data,
	}
	response, err := client.ScanFileBinary(context.Background(), request)
	checkgRPCErr(err)

	if err == nil {

		if response.Infected {
			color.Set(color.FgRed)
			log.Printf("Scan Result: %s", response.Output)
			color.Unset() // Don't forget to unset
		} else {
			color.Set(color.FgGreen)
			log.Printf("Scan Result: [CLEAN]")
			color.Unset() // Don't forget to unset
		}
	}
}

func scanFile(path string) {

	isDirectory, err := utils.IsDirectory(path)
	check(err)

	// Loop through all files in directory
	if isDirectory {
		filePaths, err := utils.WalkAllFilesInDir(path)
		check(err)

		for _, filepath := range filePaths {
			log.Println("Scanning:", filepath)
			b, err := utils.ReadAll(filepath)
			check(err)
			scanFileBinary(b)

		}
	}

}

func updateVPS() {
	_, err := client.UpdateVPS(context.Background(), &pb.UpdateVPSRequest{})
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	log.Printf("Updated with success")
}

func getLicenseStatus() {
	_, err := client.GetLicenseStatus(context.Background(), &pb.LicenseStatusRequest{})
	checkgRPCErr(err)
	if err != nil {
		log.Fatal("Please Activate the license ...")
	} else {
		log.Println("License is valid")
	}
}

func main() {

	// parse command line arguments
	flags.Parse(&opts)

	conn, err := grpc.Dial(address, []grpc.DialOption{grpc.WithInsecure()}...)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client = pb.NewAvastScannerClient(conn)

	if opts.Scan != "" {
		scanFile(opts.Scan)
	}

	if opts.VPSVersion {
		getVPSVersion()
	}

	if opts.ProgramVersion {
		getProgramVersion()
	}

	if opts.ActivateLicense != "" {
		activateLicense(opts.ActivateLicense)
	}

	if opts.Update {
		updateVPS()
	}

	if opts.StatusLicense {
		getLicenseStatus()
	}

}
