package ql

import "strings"

const (
	// SELECTFromBucketQuery
	SELECTFromBucketQuery = "SELECT * FROM bucket WHERE "
)

// N1QLBuilder implements a query builder for N1QL.
type N1QLBuilder struct {
	bucket string
	fields map[string]string
	sb     strings.Builder
}

// NewN1QLBuilder creates a new N1QL builder.
func NewN1QLBuilder(bucket string) *N1QLBuilder {
	return &N1QLBuilder{
		bucket: bucket,
		fields: make(map[string]string),
		sb:     strings.Builder{},
	}
}

// TranspileSingleQuery transpiles an N1QL query for single queries.
func (n *N1QLBuilder) TranspileSingleQuery(sq SingleQuery) string {
	n.sb.Reset()
	n.sb.WriteString(SELECTFromBucketQuery)
	switch sq.Type {
	case SIZE:
		n.sb.WriteString("size")
		n.sb.WriteByte(' ')
		n.sb.WriteString(convOp(sq.Op))
		n.sb.WriteByte(' ')
		n.sb.WriteString(sq.Value)
		return n.sb.String()
	}
	return ""
}
