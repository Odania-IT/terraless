package support

import (
	"strings"
)

func ContainsStartsWith(haystack []string, needle string) bool {
	for _, a := range haystack {
		if strings.HasPrefix(a, needle) {
			return true
		}
	}

	return false
}

func SanitizeString(val string) string {
	return strings.Replace(val, ".", "-", -1)
}

func Contains(haystack []string, needle string) bool {
	for _, a := range haystack {
		if a == needle {
			return true
		}
	}

	return false
}
