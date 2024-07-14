package ql

// Type defines the support modifiers we want to support in a search query
type Type string

// Literal encodes a literal value
type Literal string

// Token represents the actual token holds the type and it's literal representation.
type Token struct {
	Type
	Literal
}

// New creates a new token instance
func NewToken(typ Type, ch byte) Token {
	return Token{
		Type:    typ,
		Literal: Literal(ch),
	}
}

// keywords defines a map of modifier keywords and their respective tokens.
var keywords = map[string]Type{
	"type":      TYPE,
	"size":      SIZE,
	"ls":        LS,
	"fs":        FS,
	"positives": POSITIVES,
	"name":      NAME,
	"strings":   STRINGS,
	"kb":        KB,
	"mb":        MB,
	"section":   SECTION,
	"imports":   IMPORTS,
	"exports":   EXPORTS,
}

const (
	// Special modifiers
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	WS      = "WS"

	// Literals and Field names (the ones occuring in file.json)
	IDENT = "IDENTIFIER"
	// Integers
	INT = "INT"
	// Strings
	STRING = "STRING"

	// Separators
	COMMA    = ","
	LBRACKET = "["
	RBRACKET = "]"

	// Conditionals
	OR  = "OR"
	AND = "AND"

	// Value assignement
	ASSIGN = ":"
	// Operators
	EQUAL       = "="
	GREATERTHAN = "+"
	LESSERTHAN  = "-"

	// Keywords
	// Modifiers on all filetypes.
	TYPE      = "TYPE"
	SIZE      = "SIZE"
	LS        = "LS"
	FS        = "FS"
	POSITIVES = "POSITIVES"
	NAME      = "NAME"
	STRINGS   = "STRINGS"
	// Modifiers on binary executable formats (PE/ELF/Mach-o)
	SECTION = "SECTION"
	IMPORTS = "IMPORTS"
	EXPORTS = "EXPORTS"

	// Qualifiers
	KB = "kb"
	MB = "mb"
)

// LookupIdent checks a given identifier against the keyword table.
func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
