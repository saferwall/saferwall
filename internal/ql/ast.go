package ql

import (
	"strings"
)

// Node represents an AST node
type Node interface {
	Literal() string
}

// QueryStatement represents a query statement in the form (IDENTIFER:VALUE).
type QueryStatement interface {
	Node
	statementNode()
}

// SearchQuery represents a list of query statements.
type SearchQuery struct {
	Statements []QueryStatement
}

// Literal implements the Node interface for search query.
func (s SearchQuery) Literal() string {
	if len(s.Statements) > 0 {
		return s.Statements[0].Literal()
	}
	return ""
}

// SingleQuery represent queries of the form identifier:value
// such as type:elf or size:+100kb these operates with a single rhs value.
// Such queries are parsed as {SIZE,=/>/<,100} or {TYPE,=,ELF}
type SingleQuery struct {
	Token
	Op    string
	Value string
}

// statmentNode to implement the QueryStatement interface.
func (sq SingleQuery) statementNode() {}

// Literal to implement the QueryStatement interface.
func (sq SingleQuery) Literal() string {
	var sb strings.Builder
	sb.WriteString(string(sq.Token.Literal))
	sb.WriteString(sq.Op)
	sb.WriteString(sq.Value)
	return sb.String()
}
