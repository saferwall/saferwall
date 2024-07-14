package ql

import "strconv"

// mbToBytes converts a MB value to its bytes value.
func mbToBytes(v int) int {
	return v * 1e6
}

// kbToBytes converts a KB value to its bytes value.
func kbToBytes(v int) int {
	return v * 1e3
}

// mbToBytesStr converts a MB value to its bytes value for strings.
func mbToBytesStr(v string) string {
	n, _ := strconv.Atoi(v)
	return strconv.Itoa(mbToBytes(n))
}

// kbToBytesStr converts a MB value to its bytes value for strings.
func kbToBytesStr(v string) string {
	n, _ := strconv.Atoi(v)
	return strconv.Itoa(kbToBytes(n))
}

// convOp converts operator from +/-/= to >= <= or =.
func convOp(op string) string {
	switch op {
	case "+":
		return ">="
	case "-":
		return "<="
	case "=":
		return "="
	}
	return "="
}
