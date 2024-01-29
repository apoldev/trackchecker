package usecase

import (
	"encoding/json"
	"github.com/apoldev/trackchecker/internal/app/crawler"
	"github.com/apoldev/trackchecker/internal/app/crawler/repo"
	"github.com/apoldev/trackchecker/internal/app/models"
	repo2 "github.com/apoldev/trackchecker/internal/app/track/repo"
	"github.com/apoldev/trackchecker/pkg/logger"
	"github.com/google/uuid"
)

type Publisher interface {
	Publish(message []byte) error
}

// TrackResultRepo is an interface for save and get tracking result.
type TrackResultRepo interface {
	Set(id string, b []byte) error
	Get(id string) ([]byte, error)
}

type Tracking struct {
	publisher       Publisher
	logger          logger.Logger
	trackSpiderRepo *repo.SpiderRepo
	trackRepo       TrackResultRepo
}

func NewTracking(publisher Publisher, logger logger.Logger, trackSpiderRepo *repo.SpiderRepo, trackRepo *repo2.TrackRepo) *Tracking {
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

	t.logger.Debugf("got %s track results on %s", len(cr.GetResults()), track.UUID)

	return cr.GetResults(), nil
}
