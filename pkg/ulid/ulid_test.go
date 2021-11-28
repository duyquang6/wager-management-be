package ulid

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetUniqueID(t *testing.T) {
	assert.NotEqual(t, GetUniqueID(), GetUniqueID())
}
