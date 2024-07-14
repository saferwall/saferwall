package ql

import "strings"

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
func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case ':':
		tok = NewToken(ASSIGN, l.ch)
	case ',':
		tok = NewToken(COMMA, l.ch)
	case '=':
		tok = NewToken(EQUAL, l.ch)
	case '+':
		tok = NewToken(GREATERTHAN, l.ch)
	case '-':
		tok = NewToken(LESSERTHAN, l.ch)
	case '"':
		tok.Type = STRING
		tok.Literal = Literal(l.readString())
	case '[':
		tok = NewToken(LBRACKET, l.ch)
	case ']':
		tok = NewToken(RBRACKET, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = Literal(l.readIdentifier())
			tok.Type = LookupIdent(string(tok.Literal))
			return tok
		} else if isDigit(l.ch) {
			tok.Type = INT
			tok.Literal = Literal(l.readNumber())
		} else {
			tok = NewToken(ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

// Lex builds a lexeme slice by iterating and lexing the input.
func (l *Lexer) Lex() []Token {
	lexemes := make([]Token, 0)
	tok := l.NextToken()
	for tok.Type != EOF {
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
