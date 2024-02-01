package crawler

import (
	"errors"
	"net/http"
	"sync"

	"github.com/apoldev/trackchecker/internal/app/models"
	"github.com/apoldev/trackchecker/pkg/logger"
	"github.com/apoldev/trackchecker/pkg/scraper"
)

const (
	StatusFinish    = "finish"
	StatusRunning   = "running"
	StatusNoSpiders = "no_spiders"
)

var (
	ErrTrackIsNil = errors.New("tracking number is nil")
)

// SpiderFinder can found spiders by tracking number regexp.
type SpiderFinder interface {
	FindSpidersByTrackingNumber(trackingNumber string) []*models.Spider
}

// Manager creates Crawler instances for starts spiders.
// Crawler executes N tasks in goroutines and collect results.
type Manager struct {
	SpiderFinder SpiderFinder
	logger       logger.Logger
	client       *http.Client
}

func NewCrawlerManager(repo SpiderFinder, log logger.Logger, client *http.Client) *Manager {
	return &Manager{
		SpiderFinder: repo,
		logger:       log,
		client:       client,
	}
}

type crawlerTask struct {
	spider         *models.Spider
	trackingNumber string
}

// Start can start spiders in parallel and wait for all spiders to finish.
// For each spider we have to create new args for scraper with tracking number
// for replace in url or body or headers, etc...
//
// After scraping, we can get result from scraper and save it to results.
// If we have error, we can save it to results too.
//
// One spider can find another tracking number in result and start new spider for this tracking number.
func (c *Manager) Start(track *models.TrackingNumber) (*models.Crawler, error) {
	if track == nil {
		return nil, ErrTrackIsNil
	}

	crawler := models.Crawler{
		TrackingNumber: *track,
		Results:        make([]models.CrawlerResult, 0),
		Status:         StatusRunning,
	}
	spiders := c.SpiderFinder.FindSpidersByTrackingNumber(track.Code)
	if len(spiders) == 0 {
		crawler.Status = StatusNoSpiders
		return &crawler, nil
	}
	wg := sync.WaitGroup{}
	chTasks := make(chan crawlerTask)
	resultMutex := sync.Mutex{}
	mu := sync.Mutex{}
	usedSpidersWithTrackingNumber := make(map[string]struct{})

	wg.Add(len(spiders))
	go func() {
		for i := range spiders {
			chTasks <- crawlerTask{
				spider:         spiders[i],
				trackingNumber: track.Code,
			}
		}
	}()

	go func() {
		for task := range chTasks {
			go c.runCrawlerTask(&crawler, task, &wg, &resultMutex, &mu, usedSpidersWithTrackingNumber)
		}
	}()

	wg.Wait()
	close(chTasks)
	crawler.Status = StatusFinish
	return &crawler, nil
}

func (c *Manager) runCrawlerTask(
	crawler *models.Crawler,
	task crawlerTask,
	wg *sync.WaitGroup,
	resultMutex *sync.Mutex,
	mu *sync.Mutex,
	usedSpidersWithTrackingNumber map[string]struct{},
) {
	defer wg.Done()
	taskKey := task.trackingNumber + ":" + task.spider.Code
	mu.Lock()
	if _, ok := usedSpidersWithTrackingNumber[taskKey]; ok {
		defer mu.Unlock()
		return
	}
	usedSpidersWithTrackingNumber[taskKey] = struct{}{}
	mu.Unlock()

	var result models.CrawlerResult
	args := scraper.NewArgs(scraper.Variables{
		"[track]": task.trackingNumber,
	}, nil)

	// Start Scrape
	err := task.spider.Scrape(args)
	if err != nil {
		result = models.CrawlerResult{
			Spider:         task.spider.Code,
			TrackingNumber: task.trackingNumber,
			Err:            err.Error(),
			ExecuteTime:    args.ExecuteTime.Seconds(),
		}
	} else {
		result = models.CrawlerResult{
			Spider:         task.spider.Code,
			TrackingNumber: task.trackingNumber,
			ExecuteTime:    args.ExecuteTime.Seconds(),
			Result:         args.ResultBuilder.GetData(),
		}
	}
	// todo add new task if exist new tracking number or countryTo
	resultMutex.Lock()
	defer resultMutex.Unlock()
	crawler.Results = append(crawler.Results, result)
}
