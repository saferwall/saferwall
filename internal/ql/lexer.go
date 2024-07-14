package ql

import (
	"strings"

	"github.com/saferwall/saferwall/internal/ql/token"
)

const (
	// MaxStringSize sets an upper bound on lengths of strings to be parsed
	MaxStringSize = 256
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

// NewLexer creates a new lexer instance.
func NewLexer(input string) *Lexer {
	l := &Lexer{
		input: input,
	}
	l.readChar()
	return l
}

// NextToken parses and returns the next token.
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case ':':
		tok = token.New(token.Colon, string(l.ch))
	case ',':
		tok = token.New(token.Comma, string(l.ch))
	case '+':
		tok = token.New(token.Plus, string(l.ch))
	case '-':
		tok = token.New(token.Minus, string(l.ch))
	case '"':
		tok.Kind = token.StringLiteral
		tok.Literal = token.Literal(l.readString())
	case '[':
		tok = token.New(token.LBracket, string(l.ch))
	case ']':
		tok = token.New(token.RBracket, string(l.ch))
	case 0:
		tok.Literal = ""
		tok.Kind = token.EOF
	default:
		if isLetter(l.ch) {
			ident := l.readIdentifier()
			if modifier, ok := token.GetModifier(ident); ok {
				tok.Literal = token.Literal(ident)
				tok.Kind = modifier
			} else {

				tok.Literal = token.Literal(l.readIdentifier())
				tok.Kind = token.StringLiteral
			}
			return tok
		} else if isDigit(l.ch) {
			tok.Kind = token.IntegerLiteral
			tok.Literal = token.Literal(l.readNumber())
		} else {
			tok = token.New(token.Unknown, string(l.ch))
		}
	}
	l.readChar()
	return tok
}

// Lex builds a lexeme slice by iterating and lexing the input.
func (l *Lexer) Lex() []token.Token {
	lexemes := make([]token.Token, 0)
	tok := l.NextToken()
	for tok.Kind != token.EOF {
		lexemes = append(lexemes, tok)
		tok = l.NextToken()
	}
	return lexemes
}

// Analyze runs a heuristic search analysis on the lexer output
// to check for common typos and invalid inputs.
func (l *Lexer) Analyze() error {
	return nil
}

// readChar reads the next character.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

// readIdentifier reads an identifier (rhs value).
func (l *Lexer) readIdentifier() string {
	pos := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.position]
}

// readNumber reads a number (rhs value).
func (l *Lexer) readNumber() string {
	pos := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	// when integers are followed by chars e.g 100kb
	// the last digit will make read position increment
	// skipping the 'k' when reading it next
	// to avoid this we decrement readPos before returning
	l.readPosition--
	return l.input[pos:l.position]
}

// readString reads a quoted string (rhs value).
func (l *Lexer) readString() string {
	var sb strings.Builder
	len := 0

	for {
		l.readChar()
		len++
		if l.ch == '"' || l.ch == 0 || len > MaxStringSize {
			break
		}
		sb.WriteByte(l.ch)
	}
	return sb.String()
}

// skipWhitespace consumes whitespace.
func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.ch) {
		l.readChar()
	}
}
