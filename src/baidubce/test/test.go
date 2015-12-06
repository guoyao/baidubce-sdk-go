// Package test store some common data for tests and examples
package test

import (
	"fmt"
)

// Format returns a formatted string.
func Format(funcName, got, expected string) string {
	return fmt.Sprintf("%s failed. Got %s, expected %s", funcName, got, expected)
}
