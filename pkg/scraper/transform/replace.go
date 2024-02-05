package transform

import "strings"

func ReplaceStr(str, old, newStr string) string {
	return strings.ReplaceAll(str, old, newStr)
}
