package strings

import (
	"io/ioutil"
	"testing"
)

func TestStrings(t *testing.T) {
	t.Run("TestExtractStrings", func(t *testing.T) {
		testCases := []struct {
			filename string
			expected []string
		}{
			{
				filename: "../../testdata/eicar.com",
				expected: []string{`X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*`},
			},
		}

		for _, tt := range testCases {
			b, _ := ioutil.ReadFile(tt.filename)
			ascii := GetASCIIStrings(&b, 5)
			wide := GetUnicodeStrings(&b, 5)

			for i, s := range ascii {
				if s != tt.expected[i] {
					t.Fatalf("failed to get ASCII strings expected %s got %s", tt.expected[i], s)
				}
			}
			for i, s := range wide {
				if s != tt.expected[i] {
					t.Fatalf("failed to get Unicode strings expected %s got %s", tt.expected[i], s)
				}
			}
		}
	})
	t.Run("TestStringDecoder", func(t *testing.T) {
		testCases := []struct {
			in  []byte
			out string
		}{
			{
				[]byte{0x61, 0x6f, 0x73, 0x6a, 0x64, 0x66, 0x6b, 0x7a, 0x6c, 0x7a, 0x6b, 0x64, 0x6f, 0x61, 0x73, 0x6c, 0x63, 0x6b, 0x6a, 0x7a, 0x6e, 0x78},
				"潡橳晤穫穬摫慯汳正穪确",
			},
		}

		for _, tt := range testCases {
			s, err := decodeUTF16([]byte(tt.in))
			if err != nil || s != tt.out {
				t.Fatalf("failed to decode UTF-16 string expected %v got %v with error : %v", s, tt.out, err)
			}
		}
	})
	t.Run("TestASMStrings", func(t *testing.T) {
		testCases := []struct {
			filename string
			expected []string
		}{
			{
				filename: "../../testdata/ls",
				expected: []string{""},
			},
		}
		for _, tt := range testCases {
			b, _ := ioutil.ReadFile(tt.filename)
			asm := GetAsmStrings(&b)
			if len(asm) > 0 {
				for i, s := range tt.expected {
					if asm[i] != s {
						t.Fatalf("failed to get ASM string expected {%v} got {%v}", s, asm[i])
					}
				}
			}
		}

	})
}
