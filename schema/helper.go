package schema

import (
	"github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

func EnrichWithData(data map[string]string, override map[string]string) map[string]string {
	result := map[string]string{}
	for key, val := range data {
		result[key] = val
	}

	for key, val := range override {
		result[key] = val
	}

	return result
}

func ProcessString(check string, arguments Arguments) string {
	r, _ := regexp.Compile(string(`\$\{[a-z]+\}`))
	matches := r.FindStringSubmatch(check)

	if len(matches) > 0 {
		for _, match := range matches {
			if "${environment}" == match {
				check = strings.Replace(check, match, arguments.Environment, -1)
			} else {
				logrus.Fatal("Not implemented! Currently only environment can be replaced!")
			}
		}
	}

	return check
}

func ProcessData(data map[string]string, arguments Arguments) map[string]string {
	for key, val := range data {
		data[key] = ProcessString(val, arguments)
	}

	return data
}
