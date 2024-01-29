package repo

import (
	"context"

	"github.com/apoldev/trackchecker/pkg/logger"
	"github.com/redis/go-redis/v9"
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

func (r *TrackRepo) Set(id string, data []byte) error {
	ctx := context.Background()
	err := r.redis.Set(ctx, "track:"+id, string(data), 0).Err()
	return err
}

func (r *TrackRepo) Get(id string) ([]byte, error) {
	ctx := context.Background()
	return r.redis.Get(ctx, "track:"+id).Bytes()
}
