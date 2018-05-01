package common

// Shared functions.

import (
	"strconv"
)

// Int64Or takes a string and if int64 conversion produces an error the value of
// parameter `or` will be returned by default.
func Int64Or(s string, or int64) int64 {
	i64, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return or
	}
	return i64
}
