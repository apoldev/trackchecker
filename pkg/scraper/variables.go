package scraper

import "strings"

type Variables map[string]string

func (v Variables) ReplaceStringFromVariables(str string) string {
	for k, v := range v {
		rreplace := k
		if strings.Contains(str, rreplace) {
			str = strings.Replace(str, rreplace, v, -1)
		}
	}

	return str
}
