package multiav

import (
	"net"
	"os"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	// grpc port to listen.
	port = ":50051"

	// grpc library default is 4MB
	maxMsgSize = 1024 * 1024 * 20

	// Path to the file which holds the last time we updated the AV engine database.
	dbUpdateDateFilePath = "av_db_update_date.txt"
)

// ScanResult av result
type ScanResult struct {
	Output   string `json:"output"`
	Infected bool   `json:"infected"`
	Update   string `json:"update"`
}

// DefaultServerOpts returns the set of default grpc ServerOption's that Tiller requires.
func DefaultServerOpts() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.MaxRecvMsgSize(maxMsgSize),
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
