package utils

import "strings"

func ArrayContains(array []string, testValue string) bool {
	for _, elemValue := range array {
		if elemValue == testValue {
			return true
		}
	}
	return false
}

func ArrayPrefixMatch(array []string, testValue string) bool {
	for _, elemValue := range array {
		if strings.HasPrefix(testValue, elemValue) {
			return true
		}
	}
	return false
}
