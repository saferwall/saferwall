package pe

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
	"sort"

	mmap "github.com/edsrzf/mmap-go"
)

// A File represents an open PE file.
type File struct {
	data mmap.MMap
}

// Open opens the named file using os.Open and prepares it for use as a PE binary.
func Open(name string) (File, error) {

	// Init an File instance
	file := File{}

	f, err := os.Open(name)
	if err != nil {
		return file, err
	}

	// Memory map the file insead of using read/write.
	data, err := mmap.Map(f, mmap.RDONLY, 0)
	if err != nil {
		f.Close()
		return file, err
	}

	file.data = data
	return file, nil
}
