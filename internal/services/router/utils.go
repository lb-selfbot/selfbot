package router

import "strings"

// stringHasPrefix checks whether or not the string contains one of the given prefixes and returns the string without the prefix
func stringHasPrefix(str string, prefix string, ignoreCase bool) (bool, string) {
	var hasPrefix bool

	if ignoreCase {
		hasPrefix = strings.HasPrefix(strings.ToLower(str), strings.ToLower(prefix))
	} else {
		hasPrefix = strings.HasPrefix(str, prefix)
	}

	if hasPrefix {
		return true, string(str[len(prefix):])
	}

	return false, str
}

// stringArrayContains checks whether or not the given string array contains the given string
func stringArrayContains(array []string, str string, ignoreCase bool) bool {
	if ignoreCase {
		str = strings.ToLower(str)
	}
	for _, value := range array {
		if ignoreCase {
			value = strings.ToLower(value)
		}
		if value == str {
			return true
		}
	}
	return false
}
