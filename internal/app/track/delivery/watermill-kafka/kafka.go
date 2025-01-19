package watermillkafka

import (
	"context"

	"github.com/Shopify/sarama"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/apoldev/trackchecker/internal/app/config"
	"github.com/apoldev/trackchecker/internal/pkg/logger"
)

type KafkaBroker struct {
	publisher message.Publisher
	logger    logger.Logger
	cfg       config.Kafka
}

func NewKafkaBroker(cfg config.Kafka, logger logger.Logger) (*KafkaBroker, error) {
	kafkaPublisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   []string{cfg.Server},
			Marshaler: kafka.DefaultMarshaler{},
		},
		watermill.NewStdLogger(true, false),
	)
	if err != nil {
		return nil, err
	}
	logger.Infof("Configure watermill Kafka: %s %s", cfg.Topic, cfg.Server)

	return &KafkaBroker{
		publisher: kafkaPublisher,
		logger:    logger,
		cfg:       cfg,
	}, nil
}

func (b *KafkaBroker) Publish(_ context.Context, topic string, data []byte) error {
	return b.publisher.Publish(topic, message.NewMessage(
		watermill.NewUUID(),
		data,
	))
}

func (b *KafkaBroker) SubscribeQueue(
	ctx context.Context,
	topic,
	group string,
	handle func(ctx context.Context, message []byte) error,
) error {
	sub, err := b.createSubscriber(group)
	if err != nil {
		return err
	}

	for i := 0; i < b.cfg.WorkerCount; i++ {
		b.logger.Debugf("start kafka worker %d", i)
		ch, subErr := sub.Subscribe(ctx, topic)
		if subErr != nil {
			b.logger.Warnf("error subscribe to topic: %v", subErr)
			continue
		}

		go func(i int) {
			for {
				select {
				case msg := <-ch:
					b.logger.Debugf("Kafka worker %d got message %s\n", i, string(msg.Payload))
					handleErr := handle(ctx, msg.Payload)
					if handleErr != nil {
						b.logger.Warnf("Kafka worker err handle: %v", handleErr)
						continue
					}
					_ = msg.Ack()
				case <-ctx.Done():
					return
				}
			}
		}(i)
	}

	return nil
}

func (b *KafkaBroker) createSubscriber(group string) (*kafka.Subscriber, error) {
	saramaSubscriberConfig := kafka.DefaultSaramaSubscriberConfig()
	saramaSubscriberConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	return kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:               []string{b.cfg.Server},
			Unmarshaler:           kafka.DefaultMarshaler{},
			OverwriteSaramaConfig: saramaSubscriberConfig,
			ConsumerGroup:         group,
		},
		watermill.NewStdLogger(false, false),
	)
}
