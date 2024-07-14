// +build gofuzz

package ql

// Fuzz function to be used by https://github.com/dvyukov/go-fuzz
func Fuzz(data []byte) int {

	l := NewLexer(string(data))
	p := NewParser(l)

	_, err := p.Parse()
	if err != nil {
		return 1
	}

	return 0
}
