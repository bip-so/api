package redis

import (
	"strings"
)

/*
	Generate cache key generator method.
	Args:
		sender string
	 	uuids []string
	Creates a cache key string from the given inputs.
*/
func GenerateCacheKey(sender string, uuids []string) string {
	cacheKey := sender + ":"
	cacheKey += strings.Join(uuids[:], ":")
	return cacheKey
}
