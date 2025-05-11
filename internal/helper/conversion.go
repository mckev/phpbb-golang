package helper

import "strconv"

func StrToInt64(numstr string, defaultVal int64) int64 {
	val, err := strconv.ParseInt(numstr, 10, 64)
	if err != nil {
		return defaultVal
	}
	return val
}
