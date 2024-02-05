package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	watermillkafka "github.com/apoldev/trackchecker/internal/app/track/delivery/watermill-kafka"
	"github.com/apoldev/trackchecker/internal/pkg/redis"

	"github.com/apoldev/trackchecker/internal/app/config"
	usecaseCrawler "github.com/apoldev/trackchecker/internal/app/crawler"
	repo2 "github.com/apoldev/trackchecker/internal/app/crawler/repo"
	trackhttp "github.com/apoldev/trackchecker/internal/app/track/delivery/http"
	tracknats "github.com/apoldev/trackchecker/internal/app/track/delivery/nats"
	"github.com/apoldev/trackchecker/internal/app/track/repo"
	usecaseTrack "github.com/apoldev/trackchecker/internal/app/track/usecase"
	"github.com/apoldev/trackchecker/internal/pkg/grpcserver"
	"github.com/apoldev/trackchecker/internal/pkg/httpserver"
	"github.com/apoldev/trackchecker/internal/pkg/logger"

	"net"
	"net/http"
)

var (
	ErrEmptyTopicOrGroup = errors.New("topic or group name is empty")
	ErrUnknownBroker     = errors.New("unknown broker")
)

type TrackCheckerApp struct {
	config *config.Config
	logger logger.Logger

	crawlerManager *usecaseCrawler.Manager
	trackingUC     *usecaseTrack.Tracking

	broker Broker
}

type Broker interface {
	Publish(ctx context.Context, topic string, message []byte) error
	SubscribeQueue(
		ctx context.Context,
		topic,
		group string,
		handle func(ctx context.Context, message []byte) error,
	) error
}

func New(logger logger.Logger, cfg *config.Config) *TrackCheckerApp {
	return &TrackCheckerApp{
		config: cfg,
		logger: logger,
	}
}

func (a *TrackCheckerApp) Run() error {
	var err error
	repoSpider := repo2.NewSpiderRepo(a.logger)
	loadErr := repoSpider.LoadSpiders(a.config.ConfigSpiders)
	if loadErr != nil {
		return loadErr
	}
	a.logger.Infof("Loaded %d spiders", len(repoSpider.Spiders))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// configure Broker. Nats or Kafka
	var pubTopic string
	var groupName string
	if a.config.Broker == config.BrokerNats {
		pubTopic = a.config.Nats.Subject
		groupName = a.config.Nats.DurableName
	} else if a.config.Broker == config.BrokerKafka {
		pubTopic = a.config.Kafka.Topic
		groupName = a.config.Kafka.ConsumerGroup
	}
	if pubTopic == "" || groupName == "" {
		return ErrEmptyTopicOrGroup
	}

	a.broker, err = a.configureBroker()
	if err != nil {
		return err
	}

	// todo replace with Transport with proxy
	httpClient := http.DefaultClient

	redisConn := redis.NewRedisConnection(&a.config.Redis)
	trackRepo := repo.NewTrackRepo(redisConn, a.logger)

	a.crawlerManager = usecaseCrawler.NewCrawlerManager(repoSpider, a.logger, httpClient)
	a.trackingUC = usecaseTrack.NewTracking(a.broker, pubTopic, a.logger, a.crawlerManager, trackRepo)

	trackHandler := trackhttp.NewTrackHandler(a.logger, a.trackingUC)
	restServer := httpserver.NewOpenAPIServer(a.logger, trackHandler, a.config.HTTPServer)

	grpcServer := grpcserver.NewGRPCServer(a.logger, a.trackingUC)

	// Start consumer for queue
	go func() {
		queueErr := a.broker.SubscribeQueue(
			ctx,
			pubTopic,
			groupName,
			a.trackingUC.MessageHandle,
		)
		if queueErr != nil {
			a.logger.Fatal(queueErr)
		}
	}()

	// Start rest server
	go func() {
		servErr := restServer.Serve()
		if servErr != nil {
			a.logger.Fatal(servErr)
		}
	}()

	// Start grpc server
	go func() {
		grpcListener, listenErr := net.Listen("tcp", fmt.Sprintf(":%d", a.config.GRPCServer.Port))
		if listenErr != nil {
			a.logger.Fatal(listenErr)
		}
		_ = grpcServer.Serve(grpcListener)
		defer grpcListener.Close()
	}()

	// Graceful shutdown
	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, syscall.SIGINT, syscall.SIGTERM)
	<-chSignal
	cancel()

	grpcServer.GracefulStop()
	shutdownErr := restServer.Shutdown()
	if shutdownErr != nil {
		return shutdownErr
	}

	return nil
}

func (a *TrackCheckerApp) configureBroker() (Broker, error) {
	if a.config.Broker == config.BrokerNats {
		return tracknats.NewNatsBroker(a.config.Nats, a.logger)
	} else if a.config.Broker == config.BrokerKafka {
		return watermillkafka.NewKafkaBroker(a.config.Kafka, a.logger)
	}
	return nil, ErrUnknownBroker
}
