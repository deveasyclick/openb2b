package order

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
)

func generateOrderNumber() string {
	t := time.Now().UTC()
	entropy := ulid.Monotonic(rand.Reader, 0)
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return fmt.Sprintf("ORD-%s", id.String()) // e.g. ORD-01J9Z2HZ8G8Y3X9J6HQX5V3X4K
}
