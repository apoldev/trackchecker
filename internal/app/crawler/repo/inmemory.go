package repo

import (
	"encoding/json"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/apoldev/trackchecker/internal/pkg/logger"

	"github.com/apoldev/trackchecker/internal/app/models"
)

// SpiderRepo is a repository for spiders.
// can be loaded from local file.
type SpiderRepo struct {
	Spiders []models.Spider
	log     logger.Logger
}

func NewSpiderRepo(log logger.Logger) *SpiderRepo {
	return &SpiderRepo{
		log: log,
	}
}

// LoadSpiders loads spiders from json file and compiles regexp from mask string.
func (s *SpiderRepo) LoadSpiders(filename string) error {
	spiders := make([]models.Spider, 0)

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	err = s.readFile(file, &spiders)
	if err != nil {
		return err
	}

	s.compileRegexps(spiders)

	return nil
}

// FindSpidersByTrackingNumber returns scrapers that match the tracking number.
func (s *SpiderRepo) FindSpidersByTrackingNumber(trackingNumber string) []*models.Spider {
	spiders := make([]*models.Spider, 0)

	for i := range s.Spiders {
		if s.Spiders[i].Match(trackingNumber) {
			spiders = append(spiders, &s.Spiders[i])
		}
	}

	return spiders
}

func (s *SpiderRepo) compileRegexps(spiders []models.Spider) {
	var err error

	for i := range spiders {
		for j := range spiders[i].Masks {
			var re *regexp.Regexp

			regexpString := strings.Trim(spiders[i].Masks[j], "^$")
			re, err = regexp.Compile("^" + regexpString + "$")

			if err != nil {
				s.log.Warnf("spider regexp err: %s", err)
				continue
			}

			spiders[i].RegexpMasks = append(spiders[i].RegexpMasks, re)
		}

		s.Spiders = append(s.Spiders, spiders[i])
	}
}

func (s *SpiderRepo) readFile(reader io.Reader, spiders *[]models.Spider) error {
	dec := json.NewDecoder(reader)
	return dec.Decode(&spiders)
}
