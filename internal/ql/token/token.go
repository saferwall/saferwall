package token

// Kind represents the modifier type of a token.
type Kind int

// Literal value of a token.
type Literal string

// Token represents a token in a search query.
type Token struct {
	Literal
	Kind
}

// New creates a new token instance, attaching a literal to its kind.
func New(kind Kind, lit string) Token {
	return Token{
		Kind:    kind,
		Literal: Literal(lit),
	}
}

// Modifiers defines a map of modifier Modifiers and their respective tokens.
var Modifiers = map[string]Kind{
	"size":        FileSize,
	"content":     FileContent,
	"type":        FileType,
	"extension":   FileExtension,
	"name":        FileName,
	"positives":   Positives,
	"trid":        Trid,
	"packer":      Packer,
	"magic":       FileMagic,
	"tag":         Tag,
	"fs":          FirstSeen,
	"ls":          LastScanned,
	"crc32":       CRC32,
	"avast":       Avast,
	"avira":       Avira,
	"bitdefender": Bitdefender,
	"clamav":      Clamav,
	"comodo":      Comodo,
	"drweb":       DrWeb,
	"eset":        Eset,
	"fsecure":     FSecure,
	"kaspersky":   Kaspersky,
	"mcafee":      McAfee,
	"sophos":      Sophos,
	"symantec":    Symantec,
	"trendmicro":  TrendMicro,
	"windefender": Windefender,
	"md5":         MD5,
	"sha1":        SHA1,
	"sha256":      SHA256,
	"sha512":      SHA512,
	"ssdeep":      SSDeep,
}

const (
	// Special modifiers
	Unknown Kind = iota
	EOF
	// Literal values.
	//
	// Hashes, we consider hashes a first citizens.
	HashLiteral
	// Number literals (in either decimal or hexadecimal representation).
	IntegerLiteral
	// Strings.
	StringLiteral
	// Date literal in ISO 8601
	DateLiteral

	// File type literals.
	Pe

	//

	// Separators and grouping symbols.
	Comma
	Colon
	LBracket
	RBracket

	// Logical operators
	LogicalOr
	LogicalAnd

	// Characteristic suffix that act as operators.
	Plus
	Minus

	// Qualifiers
	KB
	MB

	// File metadata related to platform submission.
	FirstSeen
	LastScanned

	// Platform level modifiers.
	Tag

	// First class modifiers.
	//
	// Search modifiers for file metadata.
	FileContent
	FileExtension
	FileMagic
	FileName
	FileType
	FileSize

	// Commonly used for file specific detection.
	Trid
	Packer

	// Antivirus detections.
	Positives
	Avast
	Avira
	Bitdefender
	Clamav
	Comodo
	DrWeb
	Eset
	FSecure
	Kaspersky
	McAfee
	Sophos
	Symantec
	TrendMicro
	Windefender

	// Checksums and hashes.
	CRC32
	MD5
	SHA1
	SHA256
	SHA512
	SSDeep

	// Modifiers on binary executable formats (PE/ELF/Mach-o)
	Sections
	Imports
	Exports
)

// GetModifier checks a given literal against the modifiers table, returns
// true and the modifier token kind if found. False otherwise.
func GetModifier(ident string) (Kind, bool) {
	if tok, ok := Modifiers[ident]; ok {
		return tok, true
	}
	return Unknown, false
}
