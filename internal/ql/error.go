package ql

import (
	"fmt"
)

// newSyntaxError creates a readable error string for syntax errors.
func newSyntaxError(ctx string, expected string, got string) error {
	err := fmt.Errorf(
		"[%s]syntax error : expected %s got %s", ctx, expected, got,
	)
	return err
}
