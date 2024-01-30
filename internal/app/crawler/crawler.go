package crawler

import (
	"errors"
	"github.com/apoldev/trackchecker/internal/app/models"
	"github.com/apoldev/trackchecker/pkg/scraper"
)

var (
	ErrNoSpiders = errors.New("no spiders")
)

// Crawler starts spiders for tracking packages and accumulate results.
type Crawler struct {
	trackingNumber *models.TrackingNumber
	spiders        []*models.Spider
	results        map[string]models.CrawlerResult
}

func NewCrawler(track *models.TrackingNumber, spiders []*models.Spider) *Crawler {
	return &Crawler{
		spiders:        spiders,
		trackingNumber: track,
		results:        make(map[string]models.CrawlerResult),
	}
}

// Start can starts spiders in
func (c *Crawler) Start() error {

	if len(c.spiders) == 0 {
		return ErrNoSpiders
	}

	for i := range c.spiders {
		args := scraper.NewArgs(scraper.Variables{
			"[track]": c.trackingNumber.Code,
		}, nil)

		err := c.spiders[i].Scrape(args)

		if err != nil {
			c.results[c.spiders[i].Code] = models.CrawlerResult{
				Err: err.Error(),
			}
			continue
		}

		c.results[c.spiders[i].Code] = models.CrawlerResult{
			Result:      args.ResultBuilder.GetData(),
			ExecuteTime: args.ExecuteTime.Seconds(),
		}
	}

	return nil
}

func (c *Crawler) GetResults() map[string]models.CrawlerResult {
	return c.results
}
