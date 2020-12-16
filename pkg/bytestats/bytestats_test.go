package bytestats

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBinCount(t *testing.T) {
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
		assert.EqualValues(t, tt.expected, count)
	}
}

func TestRollingWindow(t *testing.T) {
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
		assert.EqualValues(t, tt.expected, actual)
	}
}
