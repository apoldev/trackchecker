package scraper

import "strings"

type Variables map[string]string

func (v Variables) ReplaceStringFromVariables(str string) string {
	for k, v := range v {
		replacer := k
		str = strings.Replace(str, replacer, v, -1)
	}

	return str
}
