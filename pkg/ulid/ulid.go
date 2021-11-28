// Package ulid generate sorted Unique ID
package ulid

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

var t = time.Now()
var entropy = ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)

// GetUniqueID generate unique ID
func GetUniqueID() string {
	id, err := ulid.New(ulid.Timestamp(t), entropy)
	if err != nil {
		return ""
	}
	return id.String()
}
