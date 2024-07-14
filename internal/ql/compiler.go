package ql

import (
	"strings"
)

type Query interface{}

// QueryBuilder compiles a search query from its AST representation into
// a target query language, currently N1QL.
type QueryBuilder struct {
	bucket string
	fields map[string]string
	sb     strings.Builder
}

// NewQueryBuilder creates a new instance of `QueryBuilder`.
func NewQueryBuilder(bucket string) *QueryBuilder {
	return &QueryBuilder{
		bucket: bucket,
		fields: make(map[string]string),
		sb:     strings.Builder{},
	}
}

// Compile a search query into the target query language.
func (b *QueryBuilder) Compile(query Query) (string, error) {
	return "", nil
}
