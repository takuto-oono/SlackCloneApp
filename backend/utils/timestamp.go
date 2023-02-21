package utils

import (
	"time"
)

var TimeFormat = "2006-01-02 15:04:05.000000"

func GetCurrentTime() string {
	return time.Now().Format(TimeFormat)
}

func TimeFromString(dateString string) (time.Time, error) {
	return time.Parse(TimeFormat, dateString)
}
