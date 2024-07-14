package ql

import "testing"

func TestDatetime(t *testing.T) {
	t.Run("TestParseCompleteDatetime", func(t *testing.T) {
		tests := []struct {
			input        string
			expectedUTC  string
			expectedUnix int64
		}{
			{
				input:        "2012-08-21T16:59:22",
				expectedUTC:  "2012-08-21 16:59:22 +0000 UTC",
				expectedUnix: 1345568362,
			},
		}
		for i, tt := range tests {
			date, err := strToDatetime(tt.input)
			if err != nil {
				t.Fatalf("test[%d] failed with error %v", i, err)
			}
			if date.Unix() != tt.expectedUnix {
				t.Fatalf("test[%d] expected date=%d got=%d", i, tt.expectedUnix, date.Unix())

			}
			if date.String() != tt.expectedUTC {
				t.Fatalf("test[%d] expected date=%s got=%s", i, tt.expectedUTC, date.UTC())

			}
		}

	})
}
