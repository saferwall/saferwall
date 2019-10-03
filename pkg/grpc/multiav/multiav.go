// Copyright 2019 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

// Package multiav implements common routines between all gRPC client/server av engines.

package multiav

import (
	"flag"
	"log"
	"net"
	"os"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/testdata"
)

// Scanner defines the common interface for AV engines
// to scan a file.
type Scanner interface {
	Scan()
}

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr         = flag.String("server_addr", "172.17.0.2:50051", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.saferwall.com", "The server name use to verify the hostname returned by TLS handshake")
	filePath           = flag.String("file_path", "/tmp/eicar.com", "The file path to scan")
)

const (
	// grpc library default is 4MB
	maxMsgSize = 1024 * 1024 * 20

	// Path to the file which holds the last time we updated the AV engine database.
	dbUpdateDateFilePath = "/av_db_update_date.txt"

	// port is the gRPC port the server listens on.
	port = ":50051"
)

// ScanResult av result
type ScanResult struct {
	Output   string `json:"output"`
	Infected bool   `json:"infected"`
	Update   int64  `json:"update"`
}

// DefaultServerOpts returns the set of default grpc ServerOption's that Tiller requires.
func DefaultServerOpts() []grpc.ServerOption {
	return []grpc.ServerOption{grpc.MaxRecvMsgSize(maxMsgSize),
		grpc.MaxSendMsgSize(maxMsgSize),
	}
}

// NewServer creates a new grpc server.
func NewServer(opts ...grpc.ServerOption) *grpc.Server {
	return grpc.NewServer(append(DefaultServerOpts(), opts...)...)
}

// CreateListener creates a listener on TCP port 50051
func CreateListener() (net.Listener, error) {
	lis, err := net.Listen("tcp", port)
	return lis, err
}

// Serve registers reflection service on gRPC server and start serving.
func Serve(s *grpc.Server, lis net.Listener) error {
	reflection.Register(s)
	return s.Serve(lis)
}

// read reads the entire file into memory
func read(filePath string) ([]byte, error) {
	// Start by getting a file descriptor over the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Get the file size to know how much we need to allocate
	fileinfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)

	// Read the whole binary
	_, err = file.Read(buffer)
	if err != nil {
		return nil, err
	}
	return buffer, nil

}

// UpdateDate returns a unix timestamp of the date when the
// database engine was updated.
func UpdateDate() (int64, error) {
	data, err := read(dbUpdateDateFilePath)
	if err != nil {
		return 0, err
	}
	updateDate, _ := strconv.Atoi(string(data))
	return int64(updateDate), nil
}

// ParseFlags parses the cmd line flags to create grpc conn.
func ParseFlags() (string, []grpc.DialOption, string) {
	flag.Parse()
	var opts []grpc.DialOption
	if *tls {
		if *caFile == "" {
			*caFile = testdata.Path("ca.pem")
		}
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	return *serverAddr, opts, *filePath
}