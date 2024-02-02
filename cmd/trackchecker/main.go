package main

import (
	"context"
	"fmt"
	"github.com/apoldev/trackchecker/internal/pkg/grpcserver"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/apoldev/trackchecker/internal/app/config"
	usecase2 "github.com/apoldev/trackchecker/internal/app/crawler"
	repo2 "github.com/apoldev/trackchecker/internal/app/crawler/repo"
	trackhttp "github.com/apoldev/trackchecker/internal/app/track/delivery/http"
	tracknats "github.com/apoldev/trackchecker/internal/app/track/delivery/nats"
	"github.com/apoldev/trackchecker/internal/app/track/repo"
	"github.com/apoldev/trackchecker/internal/app/track/usecase"
	"github.com/apoldev/trackchecker/internal/pkg/httpserver"
	nats2 "github.com/apoldev/trackchecker/internal/pkg/nats"
	redis2 "github.com/apoldev/trackchecker/internal/pkg/redis"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{})
	logger.SetLevel(logrus.DebugLevel)

	_, cancel := context.WithCancel(context.Background())

	cfg, err := config.LoadConfig(os.Getenv("CONFIG_FILE"))
	if err != nil {
		logger.Fatal(err)
	}

	// Redis
	redisClient := redis2.NewRedisConnection(&cfg.Redis)

	// NATS
	nc, err := nats2.NewNatsConn(cfg.Nats.Server)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("Connected to " + nc.ConnectedUrl())

	js, err := nc.JetStream(nats.PublishAsyncMaxPending(cfg.Nats.JSMaxPending))
	if err != nil {
		logger.Fatal(err)
	}

	_, addErr := js.AddStream(&nats.StreamConfig{
		Name:     cfg.Nats.StreamName,
		Subjects: []string{cfg.Nats.Subject},
	})
	if addErr != nil {
		logger.Fatal(addErr)
	}

	repoSpider := repo2.NewSpiderRepo(logger)
	err = repoSpider.LoadSpiders(cfg.ConfigSpiders)
	logger.Infof("Loaded %d spiders", len(repoSpider.Spiders))

	natsPublisher := tracknats.NewTrackPublisher(nc, js, logger, cfg.Nats)
	trackRepo := repo.NewTrackRepo(redisClient, logger)

	httpClient := http.DefaultClient
	crawlerManager := usecase2.NewCrawlerManager(repoSpider, logger, httpClient)
	trackingUC := usecase.NewTracking(natsPublisher, logger, crawlerManager, trackRepo)

	natsConsumer := tracknats.NewTrackConsumer(nc, js, logger, cfg.Nats, trackingUC)

	go func() {
		err = natsConsumer.StartQueueReceiveMessages(cfg.Nats.Subject, cfg.Nats.DurableName)
		if err != nil {
			logger.Fatal(err)
		}
	}()

	trackHandler := trackhttp.NewTrackHandler(logger, trackingUC)
	server := httpserver.NewOpenAPIServer(logger, trackHandler, cfg.HTTPServer)

	defer server.Shutdown() //nolint: errcheck // ignore error and generated code from go-swagger

	go func() {
		servErr := server.Serve()
		if servErr != nil {
			logger.Fatal(servErr)
		}
	}()

	go func() {
		grpcServer := grpcserver.NewGRPCServer(logger, cfg.GRPCServer, trackingUC)
		grpcListener, listenErr := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCServer.Port))
		if listenErr != nil {
			logger.Fatal(listenErr)
		}
		grpcServer.Serve(grpcListener)
		defer grpcListener.Close()
	}()

	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, syscall.SIGINT, syscall.SIGTERM)
	<-chSignal
	cancel()

	logger.Info("Shutting down...")
}
