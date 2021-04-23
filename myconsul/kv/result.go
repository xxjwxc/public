package kv

import (
	"encoding/json"
	"time"

	"github.com/tidwall/gjson"
)

// Result ...
type Result struct {
	g   gjson.Result
	k   string
	err error
}

// Err ...
func (r *Result) Err() error {
	return r.err
}

// Get ...
func (r *Result) Get(path string) *Result {
	r.g = r.g.Get(path)
	return r
}

// Scan ...
func (r *Result) Scan(x interface{}) error {
	return json.Unmarshal([]byte(r.g.Raw), x)
}

// Float ...
func (r *Result) Float(defaultValue ...float64) float64 {
	var df float64
	if len(defaultValue) != 0 {
		df = defaultValue[0]
	}

	if !r.g.Exists() {
		return df
	}

	return r.g.Float()
}

// Int ...
func (r *Result) Int(defaultValue ...int64) int64 {
	var df int64
	if len(defaultValue) != 0 {
		df = defaultValue[0]
	}

	if !r.g.Exists() {
		return df
	}

	return r.g.Int()
}

// Uint ...
func (r *Result) Uint(defaultValue ...uint64) uint64 {
	var df uint64
	if len(defaultValue) != 0 {
		df = defaultValue[0]
	}

	if !r.g.Exists() {
		return df
	}

	return r.g.Uint()
}

// Bool ...
func (r *Result) Bool(defaultValue ...bool) bool {
	var df bool
	if len(defaultValue) != 0 {
		df = defaultValue[0]
	}

	if !r.g.Exists() {
		return df
	}

	return r.g.Bool()
}

// Bytes ...
func (r *Result) Bytes(defaultValue ...[]byte) []byte {
	var df []byte
	if len(defaultValue) != 0 {
		df = defaultValue[0]
	}

	if !r.g.Exists() {
		return df
	}

	return []byte(r.g.Raw)
}

// String
func (r *Result) String(defaultValue ...string) string {
	var df string
	if len(defaultValue) != 0 {
		df = defaultValue[0]
	}

	if !r.g.Exists() {
		return df
	}

	return r.g.String()
}

// Time ...
func (r *Result) Time(defaultValue ...time.Time) time.Time {
	var df time.Time
	if len(defaultValue) != 0 {
		df = defaultValue[0]
	}

	if !r.g.Exists() {
		return df
	}

	return r.g.Time()
}

// Key ...
func (r *Result) Key() string {
	return r.k
}
