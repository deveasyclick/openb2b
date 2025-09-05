// Package numbergen provides utilities for generating unique, prefixed identifiers
// (e.g., order numbers, invoice numbers) using ULIDs.
package numbergen

import (
	"crypto/rand"
	"fmt"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
)

// Generate creates a unique identifier string with a specified prefix.
// The generated ID is a ULID (Universally Unique Lexicographically Sortable Identifier)
// that is:
//
//   - Based on the current UTC timestamp (for sortability)
//   - Randomized with cryptographic entropy
//   - Returned as an uppercase string
//
// Example output:
//
//	Generate("ord") -> "ORD-01J9Z2HZ8G8Y3X9J6HQX5V3X4K"
//
// Parameters:
//   - prefix: a string that will be prepended (converted to uppercase) to the ULID.
//
// Returns:
//   - A unique identifier in the format "<PREFIX>-<ULID>"
func Generate(prefix string) string {
	t := time.Now().UTC()
	entropy := ulid.Monotonic(rand.Reader, 0)
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return fmt.Sprintf("%s-%s", strings.ToUpper(prefix), id.String())
}
