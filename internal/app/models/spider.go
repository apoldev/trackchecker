package models

import (
	"regexp"

	"github.com/apoldev/trackchecker/pkg/scraper"
)

// Spider can run scraper for tracking packages
// spider can be selected by regexps.
type Spider struct {
	scraper.Scraper
	Masks       []string `json:"masks"`
	RegexpMasks []*regexp.Regexp
}

func (s *Spider) Match(trackingNumber string) bool {
	for i := range s.RegexpMasks {
		if s.RegexpMasks[i].MatchString(trackingNumber) {
			return true
		}
	}

	return false
}
