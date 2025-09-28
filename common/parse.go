package common

import "strconv"

// TryParseInt64 try to parse a string to int64
func TryParseInt64(s string) (bool, int64) {
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return false, 0
	}
	return true, val
}
