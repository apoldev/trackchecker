package transform

import (
	"regexp"
	"time"

	"github.com/araddon/dateparse"
)

var (
	regexpDateUnix = regexp.MustCompile(`([0-9]{10,13})`)
)

func Date(data string) string {
	str := checkUnixTime(data)
	parsedDate, err := dateparse.ParseAny(str)
	if err != nil {
		return data
	}
	date := parsedDate.UTC()
	return date.Format(time.RFC3339)
}

func checkUnixTime(str string) string {
	subs := regexpDateUnix.FindStringSubmatch(str)
	if len(subs) > 1 {
		return subs[1]
	}
	return str
}
