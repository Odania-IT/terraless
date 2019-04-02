package support

import (
	"regexp"
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

func SanitizeSessionName(val string) string {
	re := regexp.MustCompile("[^\\w+=,.@-]")
	return re.ReplaceAllString(val, "-")
}

func Contains(haystack []string, needle string) bool {
	for _, a := range haystack {
		if a == needle {
			return true
		}
	}

	return false
}
