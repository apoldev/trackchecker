package nats

import (
	"context"
	"time"

	"github.com/apoldev/trackchecker/internal/app/config"
	"github.com/apoldev/trackchecker/internal/pkg/logger"
	natsapp "github.com/apoldev/trackchecker/internal/pkg/nats"
	"github.com/nats-io/nats.go"
)

const (
	fetchSize    = 10
	fetchMaxWait = 10 * time.Second
)

type Broker struct {
	nc     *nats.Conn
	js     nats.JetStream
	logger logger.Logger
	cfg    config.Nats
}

func NewNatsBroker(cfg config.Nats, logger logger.Logger) (*Broker, error) {
	nc, err := natsapp.NewNatsConn(cfg.Server)
	if err != nil {
		return nil, err
	}

	js, err := nc.JetStream(nats.PublishAsyncMaxPending(cfg.JSMaxPending))
	if err != nil {
		return nil, err
	}

	_, addErr := js.AddStream(&nats.StreamConfig{
		Name:     cfg.StreamName,
		Subjects: []string{cfg.Subject},
	})
	if addErr != nil {
		return nil, addErr
	}

	return &Broker{
		nc:     nc,
		js:     js,
		logger: logger,
		cfg:    cfg,
	}, nil
}

func (b *Broker) Publish(_ context.Context, topic string, message []byte) error {
	_, err := b.js.Publish(topic, message)
	return err
}

func (b *Broker) SubscribeQueue(
	ctx context.Context,
	topic,
	group string,
	handle func(ctx context.Context, message []byte) error,
) error {
	sub, err := b.js.PullSubscribe(topic, group)
	if err != nil {
		return err
	}
	b.logger.Infof("subscribed to %s", topic)

	ch := make(chan *nats.Msg, b.cfg.WorkerCount)
	for i := 0; i < b.cfg.WorkerCount; i++ {
		go b.worker(ctx, i, ch, handle)
	}

	for {
		msgs, fetchErr := sub.Fetch(fetchSize, nats.MaxWait(fetchMaxWait))
		if fetchErr != nil {
			b.logger.Debug(fetchErr)
		}

		if len(msgs) > 0 {
			for _, msg := range msgs {
				ch <- msg
			}
		}
	}

	return nil //nolint:govet // cuz todo
}

func (b *Broker) worker(
	ctx context.Context,
	i int,
	ch <-chan *nats.Msg,
	handle func(ctx context.Context, message []byte) error,
) {
	for {
		select {
		case msg := <-ch:
			b.logger.Debugf("NATS worker %d got message %s\n", i, string(msg.Data))
			handleErr := handle(ctx, msg.Data)
			if handleErr != nil {
				b.logger.Warnf("NATS worker err handle: %v", handleErr)
				_ = msg.Nak()
				continue
			}
			_ = msg.Ack()
		case <-ctx.Done():
			return
		}
	}
}
