package transform

import (
	"regexp"
	"strings"
)

var regexpEmptySymbols = regexp.MustCompile(`(?is)[\t\s\r\n\x0B]+`)

func Clean(data string) string {
	result := data
	result = strings.ReplaceAll(result, "&nbsp;", " ")
	result = regexpEmptySymbols.ReplaceAllString(result, " ")
	result = strings.TrimSpace(result)
	result = strings.Trim(result, " !,.-;_/\\")
	return result
}
