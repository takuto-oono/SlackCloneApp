package utils

import "strconv"

func StringToUint(s string) (uint, error) {
	i, err := strconv.Atoi(s)
	return uint(i), err
}

func UintToString(u uint) string {
	return strconv.FormatUint(uint64(u), 10)
}
