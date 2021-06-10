package strings

import (
	"io/ioutil"
	"testing"
)

func TestStrings(t *testing.T) {
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
		ascii := GetASCIIStrings(b, 5)
		wide := GetUnicodeStrings(b, 5)

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
}
