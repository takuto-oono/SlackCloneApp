package utils

import (
	"time"
)

func CreateDefaultTime() time.Time {
	return time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
}

func GormTimeValidate(t time.Time) time.Time {
	if t.Year() == 0001 {
		return CreateDefaultTime()
	}
	return t
}
