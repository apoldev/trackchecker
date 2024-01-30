package nats

import (
	"github.com/apoldev/trackchecker/internal/app/config"
	"github.com/apoldev/trackchecker/pkg/logger"
	"github.com/nats-io/nats.go"
)

type TrackPublisher struct {
	nc     *nats.Conn
	js     nats.JetStream
	logger logger.Logger
	cfg    config.Nats
}

func NewTrackPublisher(
	nc *nats.Conn,
	js nats.JetStream,
	log logger.Logger,
	cfg config.Nats,
) *TrackPublisher {
	return &TrackPublisher{
		nc:     nc,
		js:     js,
		logger: log,
		cfg:    cfg,
	}
}

func (p *TrackPublisher) Publish(message []byte) error {
	_, err := p.js.Publish(p.cfg.Subject, message)
	return err
}
