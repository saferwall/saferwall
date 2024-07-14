package ql

import (
	"testing"
)

func TestTranspiler(t *testing.T) {
	t.Run("TestTranspileSearchQuery", func(t *testing.T) {
		tests := []struct {
			in  string
			out string
		}{
			{
				in:  `size:100kb+`,
				out: `SELECT * FROM bucket WHERE size >= 100000`,
			},
			{
				in:  `type:pe`,
				out: `SELECT * FROM bucket WHERE fileformat = "pe"`,
			},
			{
				in:  `fs:"2012-08-2116:59:22"`,
				out: `SELECT * FROM bucket WHERE last_seen = 1345568362`,
			},
			{
				in:  `positives:10+`,
				out: `SELECT * FROM bucket WHERE ARRAY_COUNT ( ARRAY_FLATTEN ( ARRAY i.infected for i in OBJECT_VALUES(multiav.last_scan) WHEN i.infected=true end, 1)) >= 10`,
			},
			{
				in:  `positives:5`,
				out: `SELECT * FROM bucket WHERE ARRAY_COUNT ( ARRAY_FLATTEN ( ARRAY i.infected for i in OBJECT_VALUES(multiav.last_scan) WHEN i.infected=true end, 1)) = 5`,
			},
		}
		n1bl := NewQueryBuilder("bucket")
		for i, tt := range tests {
			sq, err := NewParser(NewLexer(tt.in)).ParseQuery()
			if err != nil {
				t.Fatalf("test[%d] failed with error:%q", i, err)
			}
			query, err := n1bl.Compile(sq)
			if err != nil {
				t.Fatalf("test[%d] encountered an error %s", err)
			}
			if query != tt.out {
				t.Fatalf("test[%d] failed expected %s got %s", i, tt.out, query)
			}
		}
	})
}
