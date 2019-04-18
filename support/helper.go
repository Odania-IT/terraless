package support

import (
	"github.com/sirupsen/logrus"
	"os"
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
	val = strings.Replace(val, ".", "-", -1)
	val = strings.Replace(val, "{", "-", -1)
	val = strings.Replace(val, "}", "-", -1)

	return val
}

func SanitizeSessionName(val string) string {
	re := regexp.MustCompile(`[^\w+=,.@-]`)
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

func RunningInAws() bool {
	codebuildId := os.Getenv("CODEBUILD_BUILD_ID")

	if codebuildId == "" {
		return false
	}

	logrus.Infof("Running in Codebuild with build id: %s\n", codebuildId)
	return true
}
