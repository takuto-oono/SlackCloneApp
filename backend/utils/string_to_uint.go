package utils

import "strconv"

func StringToUint(s string) (uint, error) {
	i, err := strconv.Atoi(s)
	return uint(i), err
}
