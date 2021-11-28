package null

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	uintJSON    = []byte(`4294967294`)
	nullJSON    = []byte(`null`)
	boolJSON    = []byte(`true`)
	invalidJSON = []byte(`:)`)
	float64JSON = []byte(`1.2345`)
)

func TestNewUint(t *testing.T) {
	assert.Equal(t, Uint{100, true}, NewUint(100))
}

func TestUnmarshalUint(t *testing.T) {
	var i Uint
	err := json.Unmarshal(uintJSON, &i)
	assert.NoError(t, err)
	assertUint(t, i, "uint json")

	var null Uint
	err = json.Unmarshal(nullJSON, &null)
	assert.NoError(t, err)

	var badType Uint
	err = json.Unmarshal(boolJSON, &badType)
	assert.Error(t, err)
	assertNullUint(t, badType, "wrong type json")

	var invalid Uint
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullUint(t, invalid, "invalid json")
}

func TestUnmarshalNonUinteger(t *testing.T) {
	var i Uint
	err := json.Unmarshal(float64JSON, &i)
	assert.Error(t, err)
}

func TestMarshalUint(t *testing.T) {
	i := uint(4294967294)
	data, err := json.Marshal(i)
	assert.NoError(t, err)
	assertJSONEquals(t, data, "4294967294", "non-empty json marshal")

	// invalid values should be encoded as null
	null := Uint{0, false}
	data, err = json.Marshal(null)
	assert.NoError(t, err)
	assertJSONEquals(t, data, "null", "null json marshal")
}

func TestUintScan(t *testing.T) {
	var i Uint
	err := i.Scan(4294967294)
	assert.NoError(t, err)
	assertUint(t, i, "scanned uint32")

	var null Uint
	err = null.Scan(nil)
	assert.NoError(t, err)
	assertNullUint(t, null, "scanned null")
}

func assertUint(t *testing.T, i Uint, from string) {
	if i.Uint != 4294967294 {
		t.Errorf("bad %s uint: %d ≠ %d\n", from, i.Uint, 4294967294)
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullUint(t *testing.T, i Uint, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func assertJSONEquals(t *testing.T, data []byte, cmp string, from string) {
	if string(data) != cmp {
		t.Errorf("bad %s data: %s ≠ %s\n", from, data, cmp)
	}
}
