package null

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"reflect"
)

// Uint support nullable uint type
// package go-sql not support uint type so this will override
type Uint struct {
	Uint  uint
	Valid bool // Valid is true if Uint is not NULL
}

func NewUint(value uint) Uint {
	return Uint{Uint: value, Valid: true}
}

// Scan implements the Scanner interface for NullUint
func (ni *Uint) Scan(value interface{}) error {
	var i sql.NullInt64
	if err := i.Scan(value); err != nil {
		return err
	}
	// if nil the make Valid false
	if reflect.TypeOf(value) == nil {
		*ni = Uint{}
	} else {
		*ni = Uint{uint(i.Int64), true}
	}
	return nil
}

// Value implements the driver Valuer interface.
func (ni Uint) Value() (driver.Value, error) {
	if !ni.Valid {
		return nil, nil
	}
	return int64(ni.Uint), nil
}

// MarshalJSON for Uint
func (ni Uint) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Uint)
}

// UnmarshalJSON for Uint
func (ni *Uint) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ni.Uint)
	ni.Valid = err == nil
	return err
}