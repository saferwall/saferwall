// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

// Package bytestats : vec.go implements slice processing utilities.
package bytestats

import (
	"errors"
)

var (
	errIncompatibleShape = errors.New("operands could not be broadcasted together with incompatible dimensions")
)

// shiftBytes applies right shift operator to all buffer values.
func shiftBytes(buf []byte, s int) []byte {
	b := make([]byte, len(buf))
	for idx, v := range buf {
		b[idx] = v >> s
	}
	return b
}

// apply applies fun to all entries in the list
func apply(f func(b float32) float32, list []float32) []float32 {

	r := make([]float32, len(list))

	for idx, v := range list {
		r[idx] = f(v)
	}

	return r
}

// asFloat32 casts a slice of integers to float32.
func asFloat32(ar []int) []float32 {
	newar := make([]float32, len(ar))
	for idx, v := range ar {
		newar[idx] = float32(v)
	}
	return newar
}

// nonZeroEntries returns a slice of non-zero indexes in a given slice.
func nonZeroEntries(a []int) []int {
	b := make([]int, 0, len(a))

	for idx, v := range a {
		if v != 0 {
			b = append(b, idx)
		}
	}
	return b
}

// flatten2D flattens a 2D slice to a 1D slice
func flatten2D(a [][]int) []int {
	r := make([]int, 0)

	for _, row := range a {
		r = append(r, row...)
	}
	return r
}

// vectorizeSum computes the sum of two float32 vectors
func vectorizeSum(a []int, b []int) ([]int, error) {
	if len(a) != len(b) {
		return nil, errIncompatibleShape
	}
	c := make([]int, len(b))

	for idx := range b {
		c[idx] = a[idx] + b[idx]
	}
	return c, nil
}

func sliceWithStep(blocks [][]byte, step int) [][]byte {
	r := make([][]byte, len(blocks))

	for idx := 0; idx < len(blocks); idx += step {
		r[idx] = append(r[idx], blocks[idx]...)
	}
	return r

}
