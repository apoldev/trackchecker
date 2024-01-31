package repo

import (
	"encoding/json"
	"os"
	"regexp"
	"strings"

	"github.com/apoldev/trackchecker/internal/app/models"
	"github.com/apoldev/trackchecker/pkg/logger"
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

func (s *SpiderRepo) LoadSpiders(filename string) error {
	spiders := make([]models.Spider, 0)

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	dec := json.NewDecoder(file)
	err = dec.Decode(&spiders)
	if err != nil {
		return err
	}

	for i := range spiders {
		for j := range spiders[i].Masks {
			res := strings.TrimPrefix(spiders[i].Masks[j], "^")
			res = strings.TrimSuffix(res, "^")
			re, err := regexp.Compile("^" + spiders[i].Masks[j] + "$")

			if err != nil {
				s.log.Warnf("spider regexp err: %s", err)
				continue
			}

			spiders[i].RegexpMasks = append(spiders[i].RegexpMasks, re)
		}

		s.Spiders = append(s.Spiders, spiders[i])
	}

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
