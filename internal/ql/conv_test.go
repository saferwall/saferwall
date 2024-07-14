package ql

import "testing"

func TestConv(t *testing.T) {
	t.Run("TestMBSizeConversion", func(t *testing.T) {
		testCases := []struct {
			v        int
			expected int
		}{
			{
				5,
				5000000,
			},
		}

		for _, tt := range testCases {
			if mbToBytes(tt.v) != tt.expected {
				t.Errorf("conversion failed expected %d got %d : ", tt.expected, mbToBytes(tt.v))
			}
		}
	})
	t.Run("TestKBSizeConversion", func(t *testing.T) {
		testCases := []struct {
			v        int
			expected int
		}{
			{
				200,
				200000,
			},
		}

		for _, tt := range testCases {
			if kbToBytes(tt.v) != tt.expected {
				t.Errorf("conversion failed expected %d got %d : ", tt.expected, kbToBytes(tt.v))
			}
		}
	})
	t.Run("TestMBSizeConversionString", func(t *testing.T) {
		testCases := []struct {
			v        string
			expected string
		}{
			{
				"5",
				"5000000",
			},
		}

		for _, tt := range testCases {
			if mbToBytesStr(tt.v) != tt.expected {
				t.Errorf("conversion failed expected %s got %s : ", tt.expected, mbToBytesStr(tt.v))
			}
		}
	})
	t.Run("TestKBSizeConversionStrings", func(t *testing.T) {
		testCases := []struct {
			v        string
			expected string
		}{
			{
				"200",
				"200000",
			},
		}

		for _, tt := range testCases {
			if kbToBytesStr(tt.v) != tt.expected {
				t.Errorf("conversion failed expected %s got %s : ", tt.expected, kbToBytesStr(tt.v))
			}
		}
	})
	t.Run("TestOpConversionStrings", func(t *testing.T) {
		testCases := []struct {
			v        string
			expected string
		}{
			{
				"+",
				">=",
			}, {
				"-",
				"<=",
			}, {
				"",
				"=",
			},
		}

		for _, tt := range testCases {
			if convOp(tt.v) != tt.expected {
				t.Errorf("conversion failed expected %s got %s : ", tt.expected, kbToBytesStr(tt.v))
			}
		}
	})
}
