package helper

import (
	"strconv"
	"time"
)

func StrToInt64(numstr string, defaultVal int64) int64 {
	val, err := strconv.ParseInt(numstr, 10, 64)
	if err != nil {
		return defaultVal
	}
	return val
}

func UnixTimeToStr(unixTime int64) string {
	datetime := time.Unix(unixTime, 0)
	return datetime.Format(time.RFC822)
}
