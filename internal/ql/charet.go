package ql

// chart.go implements several character classes and helper functions.

// isWhitespace checks if a given char is a whitespace
func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\n' || ch == '\t' || ch == '\r'
}

// isLetter checks if a given char is a letter in the range [aA-zZ]
func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

// isDigit checks if a given byte is a digit.
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// isOperator checks if a given character is an operator
func isOperator(ch byte) bool {
	return (ch == '+' || ch == '-')
}
