package repo_test

import (
	"encoding/json"
	"os"
	"regexp"
	"testing"

	"github.com/apoldev/trackchecker/internal/app/crawler/repo"
	"github.com/apoldev/trackchecker/internal/app/models"
	"github.com/apoldev/trackchecker/pkg/scraper"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

var tmpSpiders = []models.Spider{
	{
		Scraper: scraper.Scraper{
			Code: "gtt",
		},
		Masks:       []string{"[A-Z]{2}[0-9]{9}[A-Z]{2}"},
		RegexpMasks: []*regexp.Regexp{regexp.MustCompile("^[A-Z]{2}[0-9]{9}[A-Z]{2}$")},
	},
	{
		Scraper: scraper.Scraper{
			Code: "russian-post",
		},
		Masks: []string{"^[0-9]{14}$", "[A-Z]{2}[0-9]{9}RU"},
		RegexpMasks: []*regexp.Regexp{
			regexp.MustCompile("^[0-9]{14}$"),
			regexp.MustCompile("^[A-Z]{2}[0-9]{9}RU$"),
		},
	},
	{
		Scraper: scraper.Scraper{
			Code: "usps",
		},
		Masks: []string{"[0-9]{20}", "[A-Z]{2}[0-9]{9}US"},
		RegexpMasks: []*regexp.Regexp{
			regexp.MustCompile("^[0-9]{20}$"),
			regexp.MustCompile("^[A-Z]{2}[0-9]{9}US$"),
		},
	},
	{
		Scraper: scraper.Scraper{
			Code: "fedex",
		},
		Masks: []string{"[0-9]{12}"},
		RegexpMasks: []*regexp.Regexp{
			regexp.MustCompile("^[0-9]{12}$"),
		},
	},

	{
		Scraper: scraper.Scraper{
			Code: "broken_regexp",
		},
		Masks:       []string{"[0-9"},
		RegexpMasks: []*regexp.Regexp{},
	},
}

var tmpFile = "spiders.json"

func createTestFile(t *testing.T) func() {
	file, err := os.Create(tmpFile)
	require.NoError(t, err)

	defer file.Close()

	enc := json.NewEncoder(file)
	enc.Encode(&tmpSpiders)

	return func() {
		err = os.Remove(tmpFile)
		require.NoError(t, err)
	}
}

func TestSpiderRepo_LoadSpiders(t *testing.T) {
	logger := logrus.New()
	var err error

	// Check errors if file not exists
	spiderRepo := repo.NewSpiderRepo(logger)
	err = spiderRepo.LoadSpiders(tmpFile)
	require.ErrorContains(t, err, "no such file or directory")

	// OK
	removeFile := createTestFile(t)
	defer removeFile()

	spiderRepo = repo.NewSpiderRepo(logger)
	err = spiderRepo.LoadSpiders(tmpFile)

	require.NoError(t, err)
	require.Len(t, spiderRepo.Spiders, len(tmpSpiders))

	for i := range spiderRepo.Spiders {
		require.Equal(t, tmpSpiders[i].Code, spiderRepo.Spiders[i].Code)
		require.Len(t, spiderRepo.Spiders[i].RegexpMasks, len(tmpSpiders[i].RegexpMasks))
	}
}

func TestSpiderRepo_FindSpidersByTrackingNumber(t *testing.T) {
	logger := logrus.New()

	removeFile := createTestFile(t)
	defer removeFile()

	spiderRepo := repo.NewSpiderRepo(logger)
	spiderRepo.LoadSpiders(tmpFile)

	cases := []struct {
		name                 string
		trackingNumber       string
		expectedSpidersCount int
	}{
		{
			name:                 "russian-post digit format",
			trackingNumber:       "12345678901234",
			expectedSpidersCount: 1,
		},

		{
			name:                 "russian-post upu format",
			trackingNumber:       "HH000222333RU",
			expectedSpidersCount: 1,
		},
		{
			name:                 "usps",
			trackingNumber:       "HH000222333US",
			expectedSpidersCount: 2,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			spiders := spiderRepo.FindSpidersByTrackingNumber(c.trackingNumber)
			require.Len(t, spiders, c.expectedSpidersCount)
		})
	}
}
