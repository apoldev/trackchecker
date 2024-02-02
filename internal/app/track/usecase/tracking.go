package usecase

import (
	"context"
	"encoding/json"

	"github.com/apoldev/trackchecker/internal/pkg/logger"

	"github.com/apoldev/trackchecker/internal/app/models"
	"github.com/google/uuid"
)

// Publisher is an interface for publish message to queue.
//
//go:generate go run github.com/vektra/mockery/v2@v2.40.1 --name Publisher
type Publisher interface {
	Publish(ctx context.Context, message []byte) error
}

// TrackResultRepo is an interface for save and get tracking result.
//
//go:generate go run github.com/vektra/mockery/v2@v2.40.1 --name TrackResultRepo
type TrackResultRepo interface {
	Set(ctx context.Context, track *models.TrackingNumber, crawler *models.Crawler) error
	Get(ctx context.Context, requestID string) ([]*models.Crawler, error)
}

// Crawler is an interface for start crawler.
//
//go:generate go run github.com/vektra/mockery/v2@v2.40.1 --name Crawler
type Crawler interface {
	Start(number *models.TrackingNumber) (*models.Crawler, error)
}

type Tracking struct {
	publisher Publisher
	logger    logger.Logger
	crawler   Crawler
	trackRepo TrackResultRepo
}

func NewTracking(
	publisher Publisher,
	logger logger.Logger,
	crawler Crawler,
	trackRepo TrackResultRepo,
) *Tracking {
	return &Tracking{
		publisher: publisher,
		logger:    logger,
		crawler:   crawler,
		trackRepo: trackRepo,
	}
}

func (t *Tracking) PublishTrackingNumbersToQueue(
	ctx context.Context,
	id string,
	trackingNumbers []string,
) ([]models.TrackingNumber, error) {
	tracks := make([]models.TrackingNumber, 0, len(trackingNumbers))
	for i := range trackingNumbers {
		track := models.TrackingNumber{
			Code:      trackingNumbers[i],
			UUID:      uuid.NewString(),
			RequestID: id,
		}

		b, err := json.Marshal(&track)
		if err != nil {
			t.logger.Warnf("error marshal tracking number: %v", err)
			continue
		}
		err = t.publisher.Publish(ctx, b)
		if err != nil {
			t.logger.Warnf("error publish tracking number to queue: %v", err)
			return nil, err
		}

		tracks = append(tracks, track)
	}

	return tracks, nil
}

func (t *Tracking) GetTrackingResult(ctx context.Context, requestID string) ([]*models.Crawler, error) {
	return t.trackRepo.Get(ctx, requestID)
}

func (t *Tracking) SaveTrackingResult(
	ctx context.Context,
	track *models.TrackingNumber,
	results *models.Crawler,
) error {
	return t.trackRepo.Set(ctx, track, results)
}

// Tracking selected spiders for tracking number and starts Crawler.
func (t *Tracking) Tracking(track *models.TrackingNumber) (*models.Crawler, error) {
	crawler, err := t.crawler.Start(track)
	if err != nil {
		t.logger.Warnf("crawler err: %v", err)
		return nil, err
	}

	t.logger.Debugf("got %d track results on %s", len(crawler.Results), track.UUID)

	return crawler, nil
}
