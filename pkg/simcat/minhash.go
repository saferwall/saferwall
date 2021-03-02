package simcat

import (
	"bytes"
	"hash/fnv"
)

// MinHashes represents the number of minhashes to compute
const MinHashes = 256

// SketchRatio represents the number of minhashes to sketches ratio
const SketchRatio = 8

func minHash(features []string) {
	hashFunc := fnv.New128()
	minHashes := make([][]byte, 0, MinHashes)
	for i := 0; i < MinHashes; i++ {
		hashes := make([][]byte, 0, len(features))
		for _, feature := range features {
			h := hashFunc.Sum([]byte(feature))
			hashes = append(hashes, h)
		}
		minHashes = append(minHashes, minByteSlice(hashes))
	}
}

// minBytes compares two byte slices and returns the min value.
func minBytes(a, b []byte) []byte {

	compResult := bytes.Compare(a, b)
	if compResult < 0 {
		return a
	} else if compResult > 0 {
		return b
	}
	return a
}

// minByteSlice compares a list of byte slices and returns the smaller one.
func minByteSlice(b [][]byte) []byte {

	tmp := b[0]

	for i := range b {
		tmp = minBytes(tmp, b[i])
	}

	return tmp
}
