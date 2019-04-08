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

func ProcessString(check string, arguments Arguments, settings TerralessSettings) string {
	r, _ := regexp.Compile(string(`\$\{[a-zA-Z0-9_]+\}`))
	matches := r.FindAllString(check, -1)

	if len(matches) > 0 {
		for _, match := range matches {
			if "${environment}" == match {
				check = strings.Replace(check, match, arguments.Environment, -1)
			} else {
				found := false

				for key, val := range arguments.Variables {
					str := "${" + key + "}"
					if str == match {
						check = strings.Replace(check, match, val, -1)
						found = true
					}
				}

				for key, val := range settings.Variables {
					str := "${" + key + "}"
					if str == match {
						check = strings.Replace(check, match, val, -1)
						found = true
					}
				}

				if !found {
					logrus.Fatalf("Variable %s not found!\n", match)
				}
			}
		}
	}

	return check
}

func ProcessData(data map[string]string, arguments Arguments, settings TerralessSettings) map[string]string {
	for key, val := range data {
		data[key] = ProcessString(val, arguments, settings)
	}

	return data
}
