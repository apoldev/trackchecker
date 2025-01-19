package transform

import (
	"regexp"
	"strings"
)

func ReplaceStr(str, old, newStr string) string {
	return strings.ReplaceAll(str, old, newStr)
}

func ReplaceRegexp(str, groups, newStr string) string {
	re, err := regexp.Compile(groups)
	if err != nil {
		return str
	}
	return re.ReplaceAllString(str, newStr)
}
