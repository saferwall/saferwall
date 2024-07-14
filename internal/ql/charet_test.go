package ql

import "testing"

func TestCharet(t *testing.T) {
	t.Run("TestIsWhitespace", func(t *testing.T) {
		testCases := []struct {
			v        byte
			expected bool
		}{
			{
				' ',
				true,
			}, {
				' ',
				true,
			}, {
				'a',
				false,
			},
		}

		for _, tt := range testCases {
			if isWhitespace(tt.v) != tt.expected {
				t.Error("failed to identify whitespace")
			}
		}
	})
	t.Run("TestIsLetter", func(t *testing.T) {
		testCases := []struct {
			v        byte
			expected bool
		}{
			{
				'a',
				true,
			}, {
				'b',
				true,
			}, {
				',',
				false,
			}, {
				':',
				false,
			}, {
				'3',
				false,
			},
		}
		for _, tt := range testCases {
			if isLetter(tt.v) != tt.expected {
				t.Error("failed to identify letters")
			}
		}
	})
	t.Run("TestIsDigit", func(t *testing.T) {
		testCases := []struct {
			v        byte
			expected bool
		}{
			{
				'a',
				false,
			}, {
				'8',
				true,
			}, {
				',',
				false,
			}, {
				':',
				false,
			}, {
				'3',
				true,
			},
		}
		for _, tt := range testCases {
			if isDigit(tt.v) != tt.expected {
				t.Error("failed to identify letters")
			}
		}
	})
	t.Run("TestIsOperator", func(t *testing.T) {
		testCases := []struct {
			v        byte
			expected bool
		}{
			{
				'+',
				true,
			}, {
				'-',
				true,
			}, {
				',',
				false,
			}, {
				':',
				false,
			}, {
				'3',
				false,
			}, {
				'a',
				false,
			},
		}
		for _, tt := range testCases {
			if isOperator(tt.v) != tt.expected {
				t.Error("failed to identify operators")
			}
		}
	})
}
