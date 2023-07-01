package utils

import (
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"strconv"
)

func String(num uint64) string {
	stringNum := strconv.FormatUint(num, 10)
	return stringNum
}

func Uint64(key string) uint64 {
	if key == "" {
		return 0
	}
	num, err := strconv.ParseUint(key, 10, 64)
	if err != nil {
		logger.Error(err.Error())
	}
	return num
}

func Remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
