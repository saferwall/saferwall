// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package bytestats

import (
	"math"
)

// binCount counts the number of occurrences of each byte value in a buffer.
func binCount(buf []byte, minlength int) []int {

	count := make([]int, minlength)

	for _, b := range buf {
		count[int(b)]++
	}

	return count
}

// rollingWindow returns a rolling window of values accross 1D slice.
func rollingWindow(buf []byte, window int) [][]byte {
	r := make([][]byte, 0, len(buf))
	for i := window - 1; i < len(buf); i++ {
		r = append(r, buf[(i-window+1):i+1])
	}
	return r
}

// entropyBinCount calculates the coarse entropy histogram of byte values.
func entropyBinCount(block []byte, window int) (int, []int) {

	var H float32
	shiftedBlock := shiftBytes(block, 4)
	c := binCount(shiftedBlock, 16)
	p := apply(func(b float32) float32 {
		return b / float32(window)
	}, asFloat32(c))
	wh := nonZeroEntries(c)
	var tmp float64
	for _, entry := range wh {
		a := -float64(p[entry]) * math.Log2(float64(p[entry]))
		tmp += a
	}
	H = float32(tmp * 2.)
	Hbin := int(H * 2)
	if Hbin == 16. {
		Hbin = 15.
	}
	return Hbin, c
}

// byteEntropyHist encodes the histogram entropy values to a singular vector
func byteEntropyHist(buf []byte, step, window int) []int {

	var output = make([][]int, 16)

	for idx := range output {
		output[idx] = make([]int, 16)
	}

	if len(buf) < window {
		Hbin, c := entropyBinCount(buf, window)
		output[Hbin], _ = vectorizeSum(output[Hbin], c)
	} else {
		blocks := sliceWithStep(rollingWindow(buf, window), step)
		for _, block := range blocks {
			Hbin, c := entropyBinCount(block, window)
			output[Hbin], _ = vectorizeSum(output[Hbin], c)
		}
	}
	return flatten2D(output)
}
