package nats

import (
	"context"
	"encoding/json"
	"time"

	"github.com/apoldev/trackchecker/internal/app/config"
	"github.com/apoldev/trackchecker/internal/app/models"
	"github.com/apoldev/trackchecker/internal/app/track/usecase"
	"github.com/apoldev/trackchecker/pkg/logger"
	"github.com/nats-io/nats.go"
)

const (
	fetchSize    = 1
	fetchMaxWait = 1 * time.Second
)

type TrackConsumer struct {
	nc       *nats.Conn
	js       nats.JetStream
	logger   logger.Logger
	cfg      config.Nats
	tracking *usecase.Tracking
}

func NewTrackConsumer(
	nc *nats.Conn,
	js nats.JetStream,
	logger logger.Logger,
	cfg config.Nats,
	tracking *usecase.Tracking,
) *TrackConsumer {
	return &TrackConsumer{
		nc:       nc,
		js:       js,
		logger:   logger,
		cfg:      cfg,
		tracking: tracking,
	}
}

func (c *TrackConsumer) worker(ctx context.Context, i int, ch chan *nats.Msg) {
	for {
		select {
		case msg := <-ch:
			c.logger.Debugf("worker %d got message %s\n", i, string(msg.Data))

			track := models.TrackingNumber{}
			err := json.Unmarshal(msg.Data, &track)

			if err != nil {
				c.logger.Warnf("worker err read message: %v", err)
				msg.Ack()
				continue
			}

			results, err := c.tracking.Tracking(&track)
			if err != nil {
				c.logger.Warnf("worker err tracking: %v", err)
				continue
			}

			err = c.tracking.SaveTrackingResult(&track, results)
			if err != nil {
				c.logger.Warnf("worker err save tracking result: %v", err)
				continue
			}

			msg.Ack()

		case <-ctx.Done():
			return
		}
	}
}

func (c *TrackConsumer) StartQueueReceiveMessages(subject, durable string) error {
	ctx := context.Background()

	sub, err := c.js.PullSubscribe(subject, durable)
	if err != nil {
		return err
	}

	ch := make(chan *nats.Msg, c.cfg.WorkerCount)

	for i := 0; i < c.cfg.WorkerCount; i++ {
		go c.worker(ctx, i, ch)
	}

	for {
		msgs, err := sub.Fetch(fetchSize, nats.MaxWait(fetchMaxWait))
		if err != nil {
			c.logger.Debug(err)
		}

		if len(msgs) > 0 {
			for _, msg := range msgs {
				ch <- msg
			}
		}
	}

	return nil
}
