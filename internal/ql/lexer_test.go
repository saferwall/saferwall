package ql

import (
	"testing"

	"github.com/saferwall/saferwall/internal/ql/token"
)

func TestNextToken(t *testing.T) {
	t.Run("TestNextTokenSingle", func(t *testing.T) {
		input := `:,+-=`

		tests := []struct {
			expectedType token.Kind
			expectedLit  token.Literal
		}{
			{token.Colon, ":"},
			{token.Comma, ","},
			{token.Plus, "+"},
			{token.Minus, "-"},
		}

		l := NewLexer(input)

		for i, tt := range tests {
			tok := l.NextToken()
			if tok.Kind != tt.expectedType {
				t.Fatalf("tests[%d] - tokenType wrong, expected=%q, got=%q", i, tt.expectedType, tok.Kind)
			}
			if tok.Literal != tt.expectedLit {
				t.Fatalf("tests[%d] - tokenLiteral wrong, expected=%q, got=%q", i, tt.expectedLit, tok.Literal)
			}
		}
	})
	t.Run("TestNextTokenStrings", func(t *testing.T) {
		input := `name:"winshell.ocx"`

		tests := []struct {
			expectedType token.Kind
			expectedLit  token.Literal
		}{
			{token.FileName, "name"},
			{token.Colon, ":"},
			{token.StringLiteral, "winshell.ocx"},
		}

		l := NewLexer(input)

		for i, tt := range tests {
			tok := l.NextToken()
			if tok.Kind != tt.expectedType {
				t.Fatalf("tests[%d] - tokenType wrong, expected=%q, got=%q", i, tt.expectedType, tok.Kind)
			}
			if tok.Literal != tt.expectedLit {
				t.Fatalf("tests[%d] - tokenLiteral wrong, expected=%q, got=%q", i, tt.expectedLit, tok.Literal)
			}
		}
	})
	t.Run("TestNextTokenNumber", func(t *testing.T) {
		input := `-123456789kb,fs:2019`

		tests := []struct {
			expectedType token.Kind
			expectedLit  token.Literal
		}{
			{token.Minus, "-"},
			{token.IntegerLiteral, "123456789"},
			{token.KB, "kb"},
			{token.Comma, ","},
			{token.FirstSeen, "fs"},
			{token.Colon, ":"},
			{token.IntegerLiteral, "2019"},
		}

		l := NewLexer(input)

		for i, tt := range tests {
			tok := l.NextToken()
			if tok.Kind != tt.expectedType {
				t.Fatalf("tests[%d] - tokenType wrong, expected=%q, got=%q", i, tt.expectedType, tok.Kind)
			}
			if tok.Literal != tt.expectedLit {
				t.Fatalf("tests[%d] - tokenLiteral wrong, expected=%q, got=%q", i, tt.expectedLit, tok.Literal)
			}
		}
	})
	t.Run("TestNextTokenDate", func(t *testing.T) {
		input := `fs:+"2009-01-01T19:59:22"`

		tests := []struct {
			expectedType token.Kind
			expectedLit  token.Literal
		}{
			{token.FirstSeen, "fs"},
			{token.Colon, ":"},
			{token.Plus, "+"},
			{token.DateLiteral, "2009-01-01T19:59:22"},
		}

		l := NewLexer(input)

		for i, tt := range tests {
			tok := l.NextToken()
			if tok.Kind != tt.expectedType {
				t.Fatalf("tests[%d] - tokenType wrong, expected=%q, got=%q", i, tt.expectedType, tok.Kind)
			}
			if tok.Literal != tt.expectedLit {
				t.Fatalf("tests[%d] - tokenLiteral wrong, expected=%q, got=%q", i, tt.expectedLit, tok.Literal)
			}
		}
	})
	t.Run("TestNextTokenFull", func(t *testing.T) {
		input := `type:pe,size:102mb+,content:"k40s",section:".k40s", imports:"crypt32.dll"`

		tests := []struct {
			expectedType token.Kind
			expectedLit  token.Literal
		}{
			{token.FileType, "type"},
			{token.Colon, ":"},
			{token.Pe, "pe"},
			{token.Comma, ","},
			{token.FileSize, "size"},
			{token.Colon, ":"},
			{token.IntegerLiteral, "102"},
			{token.MB, "mb"},
			{token.Plus, "+"},
			{token.Comma, ","},
			{token.FileContent, "content"},
			{token.Colon, ":"},
			{token.LBracket, "["},
			{token.StringLiteral, "k40s"},
			{token.RBracket, "]"},
			{token.Comma, ","},
			{token.Sections, "section"},
			{token.Colon, ":"},
			{token.StringLiteral, ".k40s"},
			{token.Comma, ","},
			{token.Imports, "imports"},
			{token.Comma, ":"},
			{token.StringLiteral, "crypt32.dll"},
		}

		l := NewLexer(input)

		for i, tt := range tests {
			tok := l.NextToken()
			if tok.Kind != tt.expectedType {
				t.Fatalf("tests[%d] - tokenType wrong, expected=%q, got=%q", i, tt.expectedType, tok.Kind)
			}
			if tok.Literal != tt.expectedLit {
				t.Fatalf("tests[%d] - tokenLiteral wrong, expected=%q, got=%q", i, tt.expectedLit, tok.Literal)
			}
		}
	})
	t.Run("TestNextTokenPartial", func(t *testing.T) {
		input := `type:pe,size:102mb+,content:["k40s",".elf",".asm",".code"],positives:5+`

		tests := []struct {
			expectedType token.Kind
			expectedLit  token.Literal
		}{
			{token.FileType, "type"},
			{token.Comma, ":"},
			{token.StringLiteral, "pe"},
			{token.Comma, ","},
			{token.FileSize, "size"},
			{token.Colon, ":"},
			{token.IntegerLiteral, "102"},
			{token.MB, "mb"},
			{token.Plus, "+"},
			{token.Comma, ","},
			{token.FileContent, "content"},
			{token.Colon, ":"},
			{token.LBracket, "["},
			{token.StringLiteral, "k40s"},
			{token.Comma, ","},
			{token.StringLiteral, ".elf"},
			{token.Comma, ","},
			{token.StringLiteral, ".asm"},
			{token.Comma, ","},
			{token.StringLiteral, ".code"},
			{token.RBracket, "]"},
			{token.Comma, ","},
			{token.Positives, "positives"},
			{token.Comma, ":"},
			{token.IntegerLiteral, "5"},
			{token.Plus, "+"},
		}

		l := NewLexer(input)

		for i, tt := range tests {
			tok := l.NextToken()
			if tok.Kind != tt.expectedType {
				t.Fatalf("tests[%d] - tokenType wrong, expected=%q, got=%q", i, tt.expectedType, tok.Kind)
			}
			if tok.Literal != tt.expectedLit {
				t.Fatalf("tests[%d] - tokenLiteral wrong, expected=%q, got=%q", i, tt.expectedLit, tok.Literal)
			}
		}
	})
	t.Run("TestNextTokenIllegal", func(t *testing.T) {
		input := `!??_`

		tests := []struct {
			expectedType token.Kind
			expectedLit  token.Literal
		}{
			{token.Unknown, "!"},
			{token.Unknown, "?"},
			{token.Unknown, "?"},
			{token.Unknown, "_"},
		}

		l := NewLexer(input)

		for i, tt := range tests {
			tok := l.NextToken()
			if tok.Kind != tt.expectedType {
				t.Fatalf("tests[%d] - tokenType wrong, expected=%q, got=%q", i, tt.expectedType, tok.Kind)
			}
			if tok.Literal != tt.expectedLit {
				t.Fatalf("tests[%d] - tokenLiteral wrong, expected=%q, got=%q", i, tt.expectedLit, tok.Literal)
			}
		}
	})
	t.Run("TestLexSingle", func(t *testing.T) {
		input := `:,+-=`

		expectedTokens := []token.Token{
			{Literal: ":", Kind: token.Colon},
			{Literal: ",", Kind: token.Comma},
			{Literal: "+", Kind: token.Plus},
			{Literal: "-", Kind: token.Minus},
		}

		l := NewLexer(input)
		tokens := l.Lex()

		for i, tok := range expectedTokens {
			if tokens[i].Kind != tok.Kind {
				t.Fatalf("tests[%d] - tokenType wrong, expected=%q, got=%q", i,
					tok.Kind, tokens[i].Kind)
			}
			if tokens[i].Literal != tok.Literal {
				t.Fatalf("tests[%d] - tokenLiteral wrong, expected=%q, got=%q",
					i, tok.Literal, tokens[i].Literal)
			}
		}
	})
	t.Run("TestLexStrings", func(t *testing.T) {
		input := `name:"winshell.ocx"`

		expectedTokens := []token.Token{
			{Literal: "name", Kind: token.FileName},
			{Literal: ":", Kind: token.Colon},
			{Literal: "winshell.ocx", Kind: token.StringLiteral},
		}

		l := NewLexer(input)
		tokens := l.Lex()

		for i, tok := range expectedTokens {
			if tokens[i].Kind != tok.Kind {
				t.Fatalf("tests[%d] - tokenType wrong, expected=%q, got=%q",
					i, tok.Kind, tokens[i].Kind)
			}
			if tokens[i].Literal != tok.Literal {
				t.Fatalf("tests[%d] - tokenLiteral wrong, expected=%q, got=%q",
					i, tok.Literal, tokens[i].Literal)
			}
		}
	})
	t.Run("TestLexNumbers", func(t *testing.T) {
		input := `-123456789kb,fs:2019`

		expectedTokens := []token.Token{
			{Literal: "-", Kind: token.Minus},
			{Literal: "123456789", Kind: token.IntegerLiteral},
			{Literal: "kb", Kind: token.KB},
			{Literal: ",", Kind: token.Comma},
			{Literal: "fs", Kind: token.FirstSeen},
			{Literal: ":", Kind: token.Colon},
			{Literal: "2019", Kind: token.StringLiteral},
		}

		l := NewLexer(input)
		tokens := l.Lex()

		for i, tok := range expectedTokens {
			if tokens[i].Kind != tok.Kind {
				t.Fatalf("tests[%d] - tokenType wrong, expected=%q, got=%q",
					i, tok.Kind, tokens[i].Kind)
			}
			if tokens[i].Literal != tok.Literal {
				t.Fatalf("tests[%d] - tokenLiteral wrong, expected=%q, got=%q",
					i, tok.Literal, tokens[i].Literal)
			}
		}
	})
}
