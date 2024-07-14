package ql

import "testing"

func TestNextToken(t *testing.T) {
	t.Run("TestNextTokenSingle", func(t *testing.T) {
		input := `:,+-=`

		tests := []struct {
			expectedType Type
			expectedLit  Literal
		}{
			{ASSIGN, ":"},
			{COMMA, ","},
			{GREATERTHAN, "+"},
			{LESSERTHAN, "-"},
			{EQUAL, "="},
		}

		l := NewLexer(input)

		for i, tt := range tests {
			tok := l.NextToken()
			if tok.Type != tt.expectedType {
				t.Fatalf("tests[%d] - tokenType wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
			}
			if tok.Literal != tt.expectedLit {
				t.Fatalf("tests[%d] - tokenLiteral wrong, expected=%q, got=%q", i, tt.expectedLit, tok.Literal)
			}
		}
	})
	t.Run("TestNextTokenStrings", func(t *testing.T) {
		input := `name:"winshell.ocx"`

		tests := []struct {
			expectedType Type
			expectedLit  Literal
		}{
			{NAME, "name"},
			{ASSIGN, ":"},
			{STRING, "winshell.ocx"},
		}

		l := NewLexer(input)

		for i, tt := range tests {
			tok := l.NextToken()
			if tok.Type != tt.expectedType {
				t.Fatalf("tests[%d] - tokenType wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
			}
			if tok.Literal != tt.expectedLit {
				t.Fatalf("tests[%d] - tokenLiteral wrong, expected=%q, got=%q", i, tt.expectedLit, tok.Literal)
			}
		}
	})
	t.Run("TestNextTokenNumber", func(t *testing.T) {
		input := `-123456789kb,fs:2019`

		tests := []struct {
			expectedType Type
			expectedLit  Literal
		}{
			{LESSERTHAN, "-"},
			{INT, "123456789"},
			{KB, "kb"},
			{COMMA, ","},
			{FS, "fs"},
			{ASSIGN, ":"},
			{INT, "2019"},
		}

		l := NewLexer(input)

		for i, tt := range tests {
			tok := l.NextToken()
			if tok.Type != tt.expectedType {
				t.Fatalf("tests[%d] - tokenType wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
			}
			if tok.Literal != tt.expectedLit {
				t.Fatalf("tests[%d] - tokenLiteral wrong, expected=%q, got=%q", i, tt.expectedLit, tok.Literal)
			}
		}
	})
	t.Run("TestNextTokenDate", func(t *testing.T) {
		input := `fs:+"2009-01-01T19:59:22"`

		tests := []struct {
			expectedType Type
			expectedLit  Literal
		}{
			{FS, "fs"},
			{ASSIGN, ":"},
			{GREATERTHAN, "+"},
			{STRING, "2009-01-01T19:59:22"},
		}

		l := NewLexer(input)

		for i, tt := range tests {
			tok := l.NextToken()
			if tok.Type != tt.expectedType {
				t.Fatalf("tests[%d] - tokenType wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
			}
			if tok.Literal != tt.expectedLit {
				t.Fatalf("tests[%d] - tokenLiteral wrong, expected=%q, got=%q", i, tt.expectedLit, tok.Literal)
			}
		}
	})
	t.Run("TestNextTokenFull", func(t *testing.T) {
		input := `type:pe,size:102mb+,strings:["k40s"],section:".k40s", imports:"crypt32.dll"`

		tests := []struct {
			expectedType Type
			expectedLit  Literal
		}{
			{TYPE, "type"},
			{ASSIGN, ":"},
			{IDENT, "pe"},
			{COMMA, ","},
			{SIZE, "size"},
			{ASSIGN, ":"},
			{INT, "102"},
			{MB, "mb"},
			{GREATERTHAN, "+"},
			{COMMA, ","},
			{STRINGS, "strings"},
			{ASSIGN, ":"},
			{LBRACKET, "["},
			{STRING, "k40s"},
			{RBRACKET, "]"},
			{COMMA, ","},
			{SECTION, "section"},
			{ASSIGN, ":"},
			{STRING, ".k40s"},
			{COMMA, ","},
			{IMPORTS, "imports"},
			{ASSIGN, ":"},
			{STRING, "crypt32.dll"},
		}

		l := NewLexer(input)

		for i, tt := range tests {
			tok := l.NextToken()
			if tok.Type != tt.expectedType {
				t.Fatalf("tests[%d] - tokenType wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
			}
			if tok.Literal != tt.expectedLit {
				t.Fatalf("tests[%d] - tokenLiteral wrong, expected=%q, got=%q", i, tt.expectedLit, tok.Literal)
			}
		}
	})
	t.Run("TestNextTokenPartial", func(t *testing.T) {
		input := `type:pe,size:102mb+,strings:["k40s",".elf",".asm",".code"],positives:+5`

		tests := []struct {
			expectedType Type
			expectedLit  Literal
		}{
			{TYPE, "type"},
			{ASSIGN, ":"},
			{IDENT, "pe"},
			{COMMA, ","},
			{SIZE, "size"},
			{ASSIGN, ":"},
			{INT, "102"},
			{MB, "mb"},
			{GREATERTHAN, "+"},
			{COMMA, ","},
			{STRINGS, "strings"},
			{ASSIGN, ":"},
			{LBRACKET, "["},
			{STRING, "k40s"},
			{COMMA, ","},
			{STRING, ".elf"},
			{COMMA, ","},
			{STRING, ".asm"},
			{COMMA, ","},
			{STRING, ".code"},
			{RBRACKET, "]"},
			{COMMA, ","},
			{POSITIVES, "positives"},
			{ASSIGN, ":"},
			{GREATERTHAN, "+"},
			{INT, "5"},
		}

		l := NewLexer(input)

		for i, tt := range tests {
			tok := l.NextToken()
			if tok.Type != tt.expectedType {
				t.Fatalf("tests[%d] - tokenType wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
			}
			if tok.Literal != tt.expectedLit {
				t.Fatalf("tests[%d] - tokenLiteral wrong, expected=%q, got=%q", i, tt.expectedLit, tok.Literal)
			}
		}
	})
	t.Run("TestNextTokenIllegal", func(t *testing.T) {
		input := `!??_`

		tests := []struct {
			expectedType Type
			expectedLit  Literal
		}{
			{ILLEGAL, "!"},
			{ILLEGAL, "?"},
			{ILLEGAL, "?"},
			{ILLEGAL, "_"},
		}

		l := NewLexer(input)

		for i, tt := range tests {
			tok := l.NextToken()
			if tok.Type != tt.expectedType {
				t.Fatalf("tests[%d] - tokenType wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
			}
			if tok.Literal != tt.expectedLit {
				t.Fatalf("tests[%d] - tokenLiteral wrong, expected=%q, got=%q", i, tt.expectedLit, tok.Literal)
			}
		}
	})
	t.Run("TestLexSingle", func(t *testing.T) {
		input := `:,+-=`

		expectedTokens := []Token{
			{ASSIGN, ":"},
			{COMMA, ","},
			{GREATERTHAN, "+"},
			{LESSERTHAN, "-"},
			{EQUAL, "="},
		}

		l := NewLexer(input)
		tokens := l.Lex()

		for i, tok := range expectedTokens {
			if tokens[i].Type != tok.Type {
				t.Fatalf("tests[%d] - tokenType wrong, expected=%q, got=%q", i, tok.Type, tokens[i].Type)
			}
			if tokens[i].Literal != tok.Literal {
				t.Fatalf("tests[%d] - tokenLiteral wrong, expected=%q, got=%q", i, tok.Literal, tokens[i].Literal)
			}
		}
	})
	t.Run("TestLexStrings", func(t *testing.T) {
		input := `name:"winshell.ocx"`

		expectedTokens := []Token{
			{NAME, "name"},
			{ASSIGN, ":"},
			{STRING, "winshell.ocx"},
		}

		l := NewLexer(input)
		tokens := l.Lex()

		for i, tok := range expectedTokens {
			if tokens[i].Type != tok.Type {
				t.Fatalf("tests[%d] - tokenType wrong, expected=%q, got=%q", i, tok.Type, tokens[i].Type)
			}
			if tokens[i].Literal != tok.Literal {
				t.Fatalf("tests[%d] - tokenLiteral wrong, expected=%q, got=%q", i, tok.Literal, tokens[i].Literal)
			}
		}
	})
	t.Run("TestLexNumbers", func(t *testing.T) {
		input := `-123456789kb,fs:2019`

		expectedTokens := []Token{
			{LESSERTHAN, "-"},
			{INT, "123456789"},
			{KB, "kb"},
			{COMMA, ","},
			{FS, "fs"},
			{ASSIGN, ":"},
			{INT, "2019"},
		}

		l := NewLexer(input)
		tokens := l.Lex()

		for i, tok := range expectedTokens {
			if tokens[i].Type != tok.Type {
				t.Fatalf("tests[%d] - tokenType wrong, expected=%q, got=%q", i, tok.Type, tokens[i].Type)
			}
			if tokens[i].Literal != tok.Literal {
				t.Fatalf("tests[%d] - tokenLiteral wrong, expected=%q, got=%q", i, tok.Literal, tokens[i].Literal)
			}
		}
	})
	t.Run("TestLexFull", func(t *testing.T) {
		input := `type:pe,size:500kb-,strings:["k40s",".elf",".asm",".code"],positives:+5`

		expectedTokens := []Token{
			{TYPE, "type"},
			{ASSIGN, ":"},
			{IDENT, "pe"},
			{COMMA, ","},
			{SIZE, "size"},
			{ASSIGN, ":"},
			{INT, "500"},
			{KB, "kb"},
			{LESSERTHAN, "-"},
			{COMMA, ","},
			{STRINGS, "strings"},
			{ASSIGN, ":"},
			{LBRACKET, "["},
			{STRING, "k40s"},
			{COMMA, ","},
			{STRING, ".elf"},
			{COMMA, ","},
			{STRING, ".asm"},
			{COMMA, ","},
			{STRING, ".code"},
			{RBRACKET, "]"},
			{COMMA, ","},
			{POSITIVES, "positives"},
			{ASSIGN, ":"},
			{GREATERTHAN, "+"},
			{INT, "5"},
		}

		l := NewLexer(input)
		tokens := l.Lex()

		for i, tok := range expectedTokens {
			if tokens[i].Type != tok.Type {
				t.Fatalf("tests[%d] - tokenType wrong, expected=%q, got=%q", i, tok.Type, tokens[i].Type)
			}
			if tokens[i].Literal != tok.Literal {
				t.Fatalf("tests[%d] - tokenLiteral wrong, expected=%q, got=%q", i, tok.Literal, tokens[i].Literal)
			}
		}
	})
	t.Run("TestLexPartial", func(t *testing.T) {
		input := `type:pe,  size:102mb+,strings:["k40s"]`

		expectedTokens := []Token{
			{TYPE, "type"},
			{ASSIGN, ":"},
			{IDENT, "pe"},
			{COMMA, ","},
			{SIZE, "size"},
			{ASSIGN, ":"},
			{INT, "102"},
			{MB, "mb"},
			{GREATERTHAN, "+"},
			{COMMA, ","},
			{STRINGS, "strings"},
			{ASSIGN, ":"},
			{LBRACKET, "["},
			{STRING, "k40s"},
			{RBRACKET, "]"},
		}

		l := NewLexer(input)
		tokens := l.Lex()

		for i, tok := range expectedTokens {
			if tokens[i].Type != tok.Type {
				t.Fatalf("tests[%d] - tokenType wrong, expected=%q, got=%q", i, tok.Type, tokens[i].Type)
			}
			if tokens[i].Literal != tok.Literal {
				t.Fatalf("tests[%d] - tokenLiteral wrong, expected=%q, got=%q", i, tok.Literal, tokens[i].Literal)
			}
		}
	})
}
