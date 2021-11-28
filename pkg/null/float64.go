package null

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"reflect"
)

type Float64 sql.NullFloat64

func NewFloat64(value float64) Float64 {
	return Float64{Float64: value, Valid: true}
}

// Scan implements the Scanner interface for NullFloat64
func (ni *Float64) Scan(value interface{}) error {
	var i sql.NullFloat64
	if err := i.Scan(value); err != nil {
		return err
	}
	// if nil the make Valid false
	if reflect.TypeOf(value) == nil {
		*ni = Float64{}
	} else {
		*ni = Float64{Float64: i.Float64, Valid: true}
	}
	return nil
}

// Value implements the driver Valuer interface.
func (ni Float64) Value() (driver.Value, error) {
	if !ni.Valid {
		return nil, nil
	}
	return ni.Float64, nil
}

// MarshalJSON for Uint
func (ni Float64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Float64)
}

// UnmarshalJSON for Uint
func (ni *Float64) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ni.Float64)
	ni.Valid = err == nil
	return err
}
