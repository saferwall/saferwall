// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package bytestats

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestByteStatistics(t *testing.T) {
	t.Run("TestBinCount", func(t *testing.T) {
		testCases := []struct {
			testBuf  []byte
			expected []int
		}{
			{
				testBuf:  []byte{1, 2, 3, 4, 5},
				expected: []int{0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
		}

		for _, tt := range testCases {
			count := binCount(tt.testBuf, 256)
			for i := range tt.expected {
				if tt.expected[i] != count[i] {
					t.Fatalf("index[%d]: expected %v got %v", i, tt.expected[i], count[i])
				}
			}
		}
	})
	t.Run("TestRollingWindow", func(t *testing.T) {
		testCases := []struct {
			input    []byte
			window   int
			expected [][]byte
		}{
			{
				input:    []byte{1, 2, 3, 4, 5},
				window:   1,
				expected: [][]byte{{1}, {2}, {3}, {4}, {5}},
			}, {
				input:    []byte{1, 2, 3, 4, 5},
				window:   2,
				expected: [][]byte{{1, 2}, {2, 3}, {3, 4}, {4, 5}},
			}, {
				input:    []byte{1, 2, 3, 4, 5},
				window:   3,
				expected: [][]byte{{1, 2, 3}, {2, 3, 4}, {3, 4, 5}},
			},
		}

		for _, tt := range testCases {
			actual := rollingWindow(tt.input, tt.window)
			for i := range tt.expected {
				if !bytes.Equal(tt.expected[i], actual[i]) {
					t.Fatalf("failed to assert rolling window tests : expected %v got %v", tt.expected[i], actual[i])
				}
			}
		}
	})
	t.Run("TestByteEntropyHistogram", func(t *testing.T) {
		testCase := []struct {
			testBin  string
			expected []int
		}{
			{
				testBin:  "../../testdata/putty.exe",
				expected: []int{8192, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3916, 8, 26, 2, 10, 4, 78, 36, 2, 4, 2, 2, 0, 2, 0, 4, 3806, 14, 22, 4, 8, 18, 100, 38, 2, 0, 0, 4, 0, 2, 2, 76, 3700, 28, 17, 40, 64, 15, 37, 23, 33, 13, 23, 26, 17, 9, 13, 38, 6905, 97, 88, 79, 164, 83, 94, 58, 105, 42, 70, 96, 86, 55, 63, 107, 3256, 149, 41, 63, 61, 53, 52, 48, 64, 37, 53, 47, 59, 46, 31, 36, 8487, 286, 461, 114, 334, 134, 1049, 689, 129, 32, 40, 56, 200, 58, 80, 139, 25904, 551, 6651, 1519, 2775, 1796, 25909, 12468, 803, 87, 135, 283, 381, 232, 206, 2220, 22090, 247, 8808, 4597, 4901, 3952, 25076, 13368, 734, 81, 65, 462, 367, 62, 198, 1008, 13618, 881, 6100, 3368, 5958, 3123, 15770, 7321, 719, 317, 275, 232, 462, 290, 403, 555, 15419, 1413, 4928, 1439, 5240, 2248, 11695, 5950, 1386, 693, 659, 541, 750, 499, 1123, 1313, 22475, 3444, 5195, 30094, 6228, 2899, 7045, 4973, 3955, 2466, 2615, 2589, 3132, 2707, 3077, 3602, 31080, 2722, 4120, 3142, 4650, 4942, 3718, 2425, 14543, 1930, 1255, 3042, 5241, 1763, 4096, 9635, 135574, 19154, 25137, 18798, 30898, 35304, 17797, 24985, 78416, 6189, 5037, 11041, 36878, 9101, 27966, 48157, 116075, 24344, 30669, 24828, 37453, 42581, 22775, 29501, 76308, 6507, 7258, 9978, 40570, 13619, 30638, 52144, 39512, 35261, 35337, 35861, 35700, 35207, 35829, 35368, 34958, 34723, 34863, 34419, 34880, 34011, 34432, 34887},
			},
		}

		step := 1024
		window := 2048

		for _, tt := range testCase {
			bytez, _ := ioutil.ReadFile(tt.testBin)
			vec := byteEntropyHist(bytez, step, window)
			if len(tt.expected) != len(vec) {
				t.Fatalf("failed to assert vector length want %d got %d", len(tt.expected), len(vec))

			}
			for i := range tt.expected {
				if tt.expected[i] != vec[i] {
					t.Fatalf("failed to assert equality want %v got %v", tt.expected, vec)
				}
			}
		}
	})
}
