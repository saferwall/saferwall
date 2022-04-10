// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package bytestats

const (
	// ByteCount represents the number of possible values a byte can take.
	ByteCount = 256
	// RollingWindow represents the length of splits for a byte slice.
	RollingWindow = 2048
	// SkipStep represents each the number of skipped steps when compuying the entropy
	// histogram of a byte slice.
	SkipStep = 1024
)

// ByteHistogram computes a histogram of byte values according to their
// indexes, each index i represents the occurrences of the byte value i.
func ByteHistogram(buf []byte) []int {

	return binCount(buf, ByteCount)
}

// ByteEntropyHistogram computes the byte-entropy histogram based on local features
// following the description in https://arxiv.org/pdf/1508.03096.pdf.
func ByteEntropyHistogram(buf []byte) []int {
	return byteEntropyHist(buf, SkipStep, RollingWindow)
}
