package usecase

import (
	"encoding/json"
	"github.com/apoldev/trackchecker/internal/app/crawler"
	"github.com/apoldev/trackchecker/internal/app/models"
	"github.com/apoldev/trackchecker/pkg/logger"
	"github.com/google/uuid"
)

// Publisher is an interface for publish message to queue.
//
//go:generate go run github.com/vektra/mockery/v2@v2.40.1 --name Publisher
type Publisher interface {
	Publish(message []byte) error
}

// TrackResultRepo is an interface for save and get tracking result.
//
//go:generate go run github.com/vektra/mockery/v2@v2.40.1 --name TrackResultRepo
type TrackResultRepo interface {
	Set(id string, b []byte) error
	Get(id string) ([]byte, error)
}

// SpiderRepo is an interface for get spiders by tracking number.
//
//go:generate go run github.com/vektra/mockery/v2@v2.40.1 --name SpiderRepo
type SpiderRepo interface {
	FindSpidersByTrackingNumber(trackingNumber string) []*models.Spider
}

type Tracking struct {
	publisher       Publisher
	logger          logger.Logger
	trackSpiderRepo SpiderRepo
	trackRepo       TrackResultRepo
}

func NewTracking(
	publisher Publisher,
	logger logger.Logger,
	trackSpiderRepo SpiderRepo,
	trackRepo TrackResultRepo,
) *Tracking {
	return &Tracking{
		publisher:       publisher,
		logger:          logger,
		trackSpiderRepo: trackSpiderRepo,
		trackRepo:       trackRepo,
	}
}

func (t *Tracking) PublishTrackingNumberToQueue(trackingNumber string) (models.TrackingNumber, error) {
	track := models.TrackingNumber{
		Code: trackingNumber,
		UUID: uuid.NewString(),
	}

	b, err := json.Marshal(&track)
	if err != nil {
		return models.TrackingNumber{}, err
	}

	err = t.publisher.Publish(b)
	if err != nil {
		return models.TrackingNumber{}, err
	}

	return track, nil
}

func (t *Tracking) GetTrackingResult(id string) ([]byte, error) {
	return t.trackRepo.Get(id)
}

func (t *Tracking) SaveTrackingResult(track *models.TrackingNumber, results map[string]models.CrawlerResult) error {
	b, err := json.Marshal(results)
	if err != nil {
		return err
	}

	return t.trackRepo.Set(track.UUID, b)
}

func (t *Tracking) Tracking(track *models.TrackingNumber) (map[string]models.CrawlerResult, error) {
	spiders := t.trackSpiderRepo.FindSpidersByTrackingNumber(track.Code)

	cr := crawler.NewCrawler(track, spiders)
	err := cr.Start()
	if err != nil {
		t.logger.Warnf("crawler err: %v", err)
		return nil, err
	}

	t.logger.Debugf("got %d track results on %s", len(cr.GetResults()), track.UUID)

	return cr.GetResults(), nil
}
