package ql

import (
	"errors"

	"github.com/saferwall/saferwall/internal/ql/token"
)

var (
	// errors
	errBadQuerySyntax error = errors.New("bad query syntax")
	// empty return types
	emptySingleQuery = SingleQuery{}
)

// Parser implements a Pratt recursive descent parser.
type Parser struct {
	l *Lexer

	currToken token.Token
	peekToken token.Token
}

// NewParser creates a new instance of the parser.
func NewParser(l *Lexer) *Parser {
	p := &Parser{
		l: l,
	}
	p.NextToken()
	p.NextToken()
	return p
}

// NextToken reads the next token in the lexemes.
func (p *Parser) NextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// Parse the lexems in the search query.
func (p *Parser) Parse() (SearchQuery, error) {
	query := SearchQuery{}
	query.Statements = make([]QueryStatement, 0)

	for p.currToken.Kind != token.EOF {
		q, err := p.ParseQuery()
		if err != nil {
			return query, err
		}
		if q != nil {
			query.Statements = append(query.Statements, q)
		}
		p.NextToken()
	}
	return query, nil
}

// ParseQuery will parse a given query type (filetype,filesize,datetimes..,)
// and return the query statement.
func (p *Parser) ParseQuery() (QueryStatement, error) {
	return nil, nil
}

/*
// ParseQuery will parse a given query type (filetype,filesize,datetimes..,)
// and return the query statement.
func (p *Parser) ParseQuery() (QueryStatement, error) {
	switch p.currToken.Kind {
	case token.FileSize:
		return p.parseFileSizeQuery()
	case token.FileType:
		return p.parseFileTypeQuery()
	case token.DateLiteral:
		return p.parseDatetimeQuery()
	case token.FirstSeen:
		return p.parseDatetimeQuery()
	case token.Positives:
		return p.parsePositivesQuery()
	case token.Sections:
		return p.parseSectionQuery()
	case token.Imports:
		return p.parseImportsQuery()
	case token.Exports:
		return p.parseExportsQuery()
	default:
		return nil, errBadQuerySyntax
	}
}

// parseFileSizeQuery parses size queries by including the range
// and transforming filesizes from KB or MB to Bytes.
func (p *Parser) parseFileSizeQuery() (SingleQuery, error) {
	query := SingleQuery{
		Token: p.currToken,
	}
	// we expect an assignment (:) next
	if !p.expectPeek(token.Colon) {
		return emptySingleQuery, newSyntaxError(, ASSIGN, string(p.peekToken.Type))
	}
	// we expect a digit next
	if !p.expectPeek(INT) {
		return emptySingleQuery, newSyntaxError(SIZE, INT, string(p.peekToken.Type))
	}
	// query value is number
	query.Value = string(p.currToken.Literal)
	// if peek token is EOF for e.g `size:200`
	// then return early
	if p.expectPeek(EOF) {
		return query, nil
	}
	// if not EOF i.e `size:200+` or `size:1mb-`
	// then we need to process more
	// we expect +/- or mb/kb or EOF
	if !p.expectPeek(GREATERTHAN) && !p.expectPeek(LESSERTHAN) && !p.expectPeek(KB) && !p.expectPeek(MB) {
		return emptySingleQuery, newSyntaxError(SIZE, "one of [+,-,kb,mb]", string(p.peekToken.Type))
	}
	// At this point current token is either KB/MB or +/-
	// e.g `size:200+` => not KB/MB conversion only operation changes
	// also implies after the Op (+/-) we expect EOF
	switch p.currToken.Type {
	case KB:
		query.Value = kbToBytesStr(query.Value)
	case MB:
		query.Value = mbToBytesStr(query.Value)
	case GREATERTHAN:
		query.Op = GREATERTHAN
	case LESSERTHAN:
		query.Op = LESSERTHAN
	}
	// if previously we had KB/MB then next we expect either +/- or EOF
	switch p.peekToken.Type {
	case GREATERTHAN:
		query.Op = GREATERTHAN
	case LESSERTHAN:
		query.Op = LESSERTHAN
	}
	for !p.currTokenIs(EOF) && !p.currTokenIs(COMMA) {
		p.NextToken()
	}
	return query, nil
}

// parseFileTypeQuery parses filetype queries.
func (p *Parser) parseFileTypeQuery() (SingleQuery, error) {
	query := SingleQuery{
		Token: p.currToken,
		Op:    EQUAL,
	}
	if !p.expectPeek(ASSIGN) {
		return emptySingleQuery, newSyntaxError(TYPE, ASSIGN, string(p.peekToken.Type))
	}
	if !p.expectPeek(IDENT) {
		return emptySingleQuery, newSyntaxError(TYPE, IDENT, string(p.peekToken.Type))
	}
	query.Value = string(p.currToken.Literal)
	for !p.currTokenIs(EOF) && !p.currTokenIs(COMMA) {
		p.NextToken()
	}
	return query, nil
}

// parsePositivesQuery parses av positives queries.
func (p *Parser) parsePositivesQuery() (SingleQuery, error) {
	query := SingleQuery{
		Token: p.currToken,
		Op:    EQUAL,
	}
	if !p.expectPeek(ASSIGN) {
		return emptySingleQuery, newSyntaxError(POSITIVES, ASSIGN, string(p.peekToken.Type))
	}
	if !p.expectPeek(INT) {
		return emptySingleQuery, newSyntaxError(POSITIVES, INT, string(p.peekToken.Type))
	}
	query.Value = string(p.currToken.Literal)

	switch p.peekToken.Type {
	case GREATERTHAN:
		query.Op = GREATERTHAN
	case LESSERTHAN:
		query.Op = LESSERTHAN
	}

	for !p.currTokenIs(EOF) && !p.currTokenIs(COMMA) {
		p.NextToken()
	}
	return query, nil
}

// parseDatetimeQuery parses datetime queries
// such as last_seen, first_seen.
func (p *Parser) parseDatetimeQuery() (SingleQuery, error) {
	query := SingleQuery{
		Token: p.currToken,
		Op:    EQUAL,
	}
	if !p.expectPeek(ASSIGN) {
		return emptySingleQuery, newSyntaxError("DATE", ASSIGN, string(p.peekToken.Type))
	}
	if !p.expectPeek(STRING) {
		return emptySingleQuery, newSyntaxError("DATE", STRING, string(p.peekToken.Type))
	}
	query.Value = string(p.currToken.Literal)
	for !p.currTokenIs(EOF) && !p.currTokenIs(COMMA) {
		p.NextToken()
	}
	return query, nil
}

// TODO: parseStringsQuery parses strings queries.
func (p *Parser) parseStringsQuery() (SingleQuery, error) {
	return emptySingleQuery, nil
}

// parseSectionQuery parses SECTION queries.
func (p *Parser) parseSectionQuery() (SingleQuery, error) {
	query := SingleQuery{
		Token: p.currToken,
		Op:    EQUAL,
	}
	if !p.expectPeek(ASSIGN) {
		return emptySingleQuery, newSyntaxError(SECTION, ASSIGN, string(p.peekToken.Type))
	}
	if !p.expectPeek(STRING) {
		return emptySingleQuery, newSyntaxError(SECTION, STRING, string(p.peekToken.Type))
	}
	query.Value = string(p.currToken.Literal)
	for !p.currTokenIs(EOF) && !p.currTokenIs(COMMA) {
		p.NextToken()
	}
	return query, nil
}

// parseImportsQuery parses IMPORTS queries.
func (p *Parser) parseImportsQuery() (SingleQuery, error) {
	query := SingleQuery{
		Token: p.currToken,
		Op:    EQUAL,
	}
	if !p.expectPeek(ASSIGN) {
		return emptySingleQuery, newSyntaxError(IMPORTS, ASSIGN, string(p.peekToken.Type))
	}
	if !p.expectPeek(STRING) {
		return emptySingleQuery, newSyntaxError(IMPORTS, STRING, string(p.peekToken.Type))
	}
	query.Value = string(p.currToken.Literal)
	for !p.currTokenIs(EOF) && !p.currTokenIs(COMMA) {
		p.NextToken()
	}
	return query, nil
}

// parseExportsQuery parses EXPORTS queries.
func (p *Parser) parseExportsQuery() (SingleQuery, error) {
	query := SingleQuery{
		Token: p.currToken,
		Op:    EQUAL,
	}
	if !p.expectPeek(ASSIGN) {
		return emptySingleQuery, newSyntaxError(EXPORTS, ASSIGN, string(p.peekToken.Type))
	}
	if !p.expectPeek(STRING) {
		return emptySingleQuery, newSyntaxError(EXPORTS, STRING, string(p.peekToken.Type))
	}
	query.Value = string(p.currToken.Literal)
	for !p.currTokenIs(EOF) && !p.currTokenIs(COMMA) {
		p.NextToken()
	}
	return query, nil
}

// parseFileTypeQuery parses
// currTokenIs asserts the current token type.
func (p *Parser) currTokenIs(t Type) bool {
	return p.currToken.Type == t
}

// peekTokenIs asserts the next token type.
func (p *Parser) peekTokenIs(t Type) bool {
	return p.peekToken.Type == t
}

// expectPeek returns true or false depending on whether the next
// token is the expected one.
func (p *Parser) expectPeek(t Type) bool {
	if p.peekTokenIs(t) {
		p.NextToken()
		return true
	}
	return false
}
*/
