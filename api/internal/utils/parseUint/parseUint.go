// Package parseuint provides a helper function for parsing strings into unsigned integers.
package parseuint

import (
	"fmt"
	"strconv"
)

// ParseUint attempts to parse the string `s` as an unsigned integer in base 10.
//
// Parameters:
//   - s: the string to parse
//   - field: a label used in the error message to indicate which field caused the error
//
// Returns:
//   - uint: the parsed unsigned integer
//   - error: a wrapped error if parsing fails
func ParseUint(s, field string) (uint, error) {
	sUint, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing %s: %q: %w", field, s, err)
	}
	return uint(sUint), nil
}
