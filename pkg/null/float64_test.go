package null

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewFloat64(t *testing.T) {
	assert.Equal(t, Float64{Float64: 100.32, Valid: true}, NewFloat64(100.32))
}

func TestUnmarshalFloat64(t *testing.T) {
	var i Float64
	err := json.Unmarshal(float64JSON, &i)
	assert.NoError(t, err)
	assertFloat64(t, i, "float64 json")

	var null Float64
	err = json.Unmarshal(nullJSON, &null)
	assert.NoError(t, err)

	var badType Float64
	err = json.Unmarshal(boolJSON, &badType)
	assert.Error(t, err)
	assertNullFloat64(t, badType, "wrong type json")

	var invalid Float64
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullFloat64(t, invalid, "invalid json")
}

func TestMarshalFloat64(t *testing.T) {
	i := 1.2345
	data, err := json.Marshal(i)
	assert.NoError(t, err)
	assertJSONEquals(t, data, "1.2345", "non-empty json marshal")

	// invalid values should be encoded as null
	null := Uint{0, false}
	data, err = json.Marshal(null)
	assert.NoError(t, err)
	assertJSONEquals(t, data, "null", "null json marshal")
}

func TestFloat64Scan(t *testing.T) {
	var i Float64
	err := i.Scan(1.2345)
	assert.NoError(t, err)
	assertFloat64(t, i, "scanned float64")

	var null Float64
	err = null.Scan(nil)
	assert.NoError(t, err)
	assertNullFloat64(t, null, "scanned null")
}


func assertFloat64(t *testing.T, i Float64, from string) {
	if i.Float64 != 1.2345 {
		t.Errorf("bad %s float64: %f ≠ %f\n", from, i.Float64, 1.2345)
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullFloat64(t *testing.T, i Float64, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}
