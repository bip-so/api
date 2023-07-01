package utils

import (
	"os"
)

// AppRoot app root path
var AppRoot, _ = os.Getwd()

// GetAbsURL get absolute URL from request, refer: https://stackoverflow.com/questions/6899069/why-are-request-url-host-and-scheme-blank-in-the-development-server

func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
