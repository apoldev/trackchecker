package repo

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/apoldev/trackchecker/internal/app/models"

	"github.com/apoldev/trackchecker/pkg/logger"
	"github.com/redis/go-redis/v9"
)

const (
	lifeTime       = 60 * time.Minute
	prefixTracking = "tracking:"
	prefixTrack    = "track:"
)

var (
	ErrNotFound = errors.New("not found")
)

type TrackRepo struct {
	redis  *redis.Client
	logger logger.Logger
}

func NewTrackRepo(r *redis.Client, log logger.Logger) *TrackRepo {
	return &TrackRepo{
		redis:  r,
		logger: log,
	}
}

func (r *TrackRepo) Set(track *models.TrackingNumber, results *models.Crawler) error {
	var err error
	ctx := context.Background()
	b, err := json.Marshal(results)
	if err != nil {
		return err
	}
	err = r.redis.HSet(ctx, prefixTracking+track.RequestID, prefixTrack+track.UUID, string(b)).Err()

	if err != nil {
		return err
	}

	return r.redis.Expire(ctx, prefixTracking+track.RequestID, lifeTime).Err()
}

func (r *TrackRepo) Get(requestID string) ([]*models.Crawler, error) {
	ctx := context.Background()
	m, err := r.redis.HGetAll(ctx, "tracking:"+requestID).Result()
	if err != nil {
		return nil, ErrNotFound
	}

	results := make([]*models.Crawler, 0, len(m))
	for _, v := range m {
		result := &models.Crawler{}
		err = json.Unmarshal([]byte(v), result)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}
