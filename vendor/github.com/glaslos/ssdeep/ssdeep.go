package ssdeep

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	// ErrFileTooSmall is returned when a file contains too few bytes.
	ErrFileTooSmall = errors.New("did not process files large enough to produce meaningful results")

	// ErrBlockSizeTooSmall is returned when a file can't produce a large enough block size.
	ErrBlockSizeTooSmall = errors.New("unable to establish a sufficient block size")
)

const (
	rollingWindow uint32 = 7
	blockMin             = 3
	spamSumLength        = 64
	minFileSize          = 4096
	hashPrime     uint32 = 0x01000193
	hashInit      uint32 = 0x28021967
	b64String            = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
)

var (
	b64 = []byte(b64String)
	// Force calculates the hash on invalid input
	Force = false
)

type rollingState struct {
	window []byte
	h1     uint32
	h2     uint32
	h3     uint32
	n      uint32
}

func (rs *rollingState) rollSum() uint32 {
	return rs.h1 + rs.h2 + rs.h3
}

type ssdeepState struct {
	rollingState rollingState
	blockSize    int
	hashString1  string
	hashString2  string
	blockHash1   uint32
	blockHash2   uint32
}

func newSsdeepState() ssdeepState {
	return ssdeepState{
		blockHash1: hashInit,
		blockHash2: hashInit,
		rollingState: rollingState{
			window: make([]byte, rollingWindow),
		},
	}
}

func (state *ssdeepState) newRollingState() {
	state.rollingState = rollingState{}
	state.rollingState.window = make([]byte, rollingWindow)
}

// sumHash based on FNV hash
func sumHash(c byte, h uint32) uint32 {
	return (h * hashPrime) ^ uint32(c)
}

// rollHash based on Adler checksum
func (state *ssdeepState) rollHash(c byte) {
	rs := &state.rollingState
	rs.h2 -= rs.h1
	rs.h2 += rollingWindow * uint32(c)
	rs.h1 += uint32(c)
	rs.h1 -= uint32(rs.window[rs.n])
	rs.window[rs.n] = c
	rs.n++
	if rs.n == rollingWindow {
		rs.n = 0
	}
	rs.h3 = rs.h3 << 5
	rs.h3 ^= uint32(c)
}

// getBlockSize calculates the block size based on file size
func (state *ssdeepState) setBlockSize(n int) {
	blockSize := blockMin
	for blockSize*spamSumLength < n {
		blockSize = blockSize * 2
	}
	state.blockSize = blockSize
}

func (state *ssdeepState) processByte(b byte) {
	state.blockHash1 = sumHash(b, state.blockHash1)
	state.blockHash2 = sumHash(b, state.blockHash2)
	state.rollHash(b)
	rh := int(state.rollingState.rollSum())
	if rh%state.blockSize == (state.blockSize - 1) {
		if len(state.hashString1) < spamSumLength-1 {
			state.hashString1 += string(b64[state.blockHash1%64])
			state.blockHash1 = hashInit
		}
		if rh%(state.blockSize*2) == ((state.blockSize * 2) - 1) {
			if len(state.hashString2) < spamSumLength/2-1 {
				state.hashString2 += string(b64[state.blockHash2%64])
				state.blockHash2 = hashInit
			}
		}
	}
}

// Reader is the minimum interface that ssdeep needs in order to calculate the fuzzy hash.
// Reader groups io.Seeker and io.Reader.
type Reader interface {
	io.Seeker
	io.Reader
}

func (state *ssdeepState) process(r *bufio.Reader) {
	state.newRollingState()
	b, err := r.ReadByte()
	for err == nil {
		state.processByte(b)
		b, err = r.ReadByte()
	}
}

// FuzzyReader computes the fuzzy hash of a Reader interface with a given input size.
// It is the caller's responsibility to append the filename, if any, to result after computation.
// Returns an error when ssdeep could not be computed on the Reader.
func FuzzyReader(f Reader, fileSize int) (out string, err error) {
	if fileSize < minFileSize {
		err = ErrFileTooSmall
		if !Force {
			return
		}
	}
	state := newSsdeepState()
	state.setBlockSize(fileSize)
	for {
		if _, seekErr := f.Seek(0, 0); seekErr != nil {
			return "", seekErr
		}
		r := bufio.NewReader(f)
		state.process(r)
		if state.blockSize < blockMin {
			err = ErrBlockSizeTooSmall
			if !Force {
				return
			}
		}
		if len(state.hashString1) < spamSumLength/2 {
			state.blockSize = state.blockSize / 2
			state.blockHash1 = hashInit
			state.blockHash2 = hashInit
			state.hashString1 = ""
			state.hashString2 = ""
		} else {
			rh := state.rollingState.rollSum()
			if rh != 0 {
				// Finalize the hash string with the remaining data
				state.hashString1 += string(b64[state.blockHash1%64])
				state.hashString2 += string(b64[state.blockHash2%64])
			}
			break
		}
	}
	return fmt.Sprintf("%d:%s:%s", state.blockSize, state.hashString1, state.hashString2), err
}

// FuzzyFilename computes the fuzzy hash of a file.
// FuzzyFilename will opens, reads, and hashes the contents of the file 'filename'.
// It is the caller's responsibility to append the filename to the result after computation.
// Returns an error when the file doesn't exist or ssdeep could not be computed on the file.
func FuzzyFilename(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	return FuzzyFile(f)
}

// FuzzyFile computes the fuzzy hash of a file using os.File pointer.
// FuzzyFile will computes the fuzzy hash of the contents of the open file, starting at the beginning of the file.
// When finished, the file pointer is returned to its original position.
// If an error occurs, the file pointer's value is undefined.
// It is the callers's responsibility to append the filename to the result after computation.
// Returns an error when ssdeep could not be computed on the file.
func FuzzyFile(f *os.File) (out string, err error) {
	currentPosition, err := f.Seek(0, io.SeekCurrent)
	if err != nil {
		return
	}
	if _, err = f.Seek(0, io.SeekStart); err != nil {
		return
	}
	stat, err := f.Stat()
	if err != nil {
		return
	}
	out, err = FuzzyReader(f, int(stat.Size()))
	if err != nil {
		return
	}
	_, err = f.Seek(currentPosition, io.SeekStart)
	return
}

// FuzzyBytes computes the fuzzy hash of a slice of byte.
// It is the caller's responsibility to append the filename, if any, to result after computation.
// Returns an error when ssdeep could not be computed on the buffer.
func FuzzyBytes(buffer []byte) (string, error) {
	n := len(buffer)
	br := bytes.NewReader(buffer)

	result, err := FuzzyReader(br, n)
	if err != nil {
		return "", err
	}

	return result, nil
}
