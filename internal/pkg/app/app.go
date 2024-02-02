package app

import (
	"fmt"
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
	"github.com/apoldev/trackchecker/internal/pkg/grpcserver"
	"github.com/apoldev/trackchecker/internal/pkg/httpserver"
	"github.com/apoldev/trackchecker/internal/pkg/logger"
	"github.com/redis/go-redis/v9"

	"net"
	"net/http"

	"github.com/nats-io/nats.go"
)

type TrackCheckerApp struct {
	config *config.Config
	logger logger.Logger

	redisClient *redis.Client
	natsConn    *nats.Conn
}

func New(logger logger.Logger, cfg *config.Config, redisClient *redis.Client, natsConn *nats.Conn) *TrackCheckerApp {
	return &TrackCheckerApp{
		config:      cfg,
		logger:      logger,
		natsConn:    natsConn,
		redisClient: redisClient,
	}
}

func (a *TrackCheckerApp) Run() error {
	js, err := a.natsConn.JetStream(nats.PublishAsyncMaxPending(a.config.Nats.JSMaxPending))
	if err != nil {
		a.logger.Fatal(err)
	}

	_, addErr := js.AddStream(&nats.StreamConfig{
		Name:     a.config.Nats.StreamName,
		Subjects: []string{a.config.Nats.Subject},
	})
	if addErr != nil {
		return addErr
	}

	repoSpider := repo2.NewSpiderRepo(a.logger)
	err = repoSpider.LoadSpiders(a.config.ConfigSpiders)
	a.logger.Infof("Loaded %d spiders", len(repoSpider.Spiders))

	natsPublisher := tracknats.NewTrackPublisher(a.natsConn, js, a.logger, a.config.Nats)
	trackRepo := repo.NewTrackRepo(a.redisClient, a.logger)

	httpClient := http.DefaultClient
	crawlerManager := usecase2.NewCrawlerManager(repoSpider, a.logger, httpClient)
	trackingUC := usecase.NewTracking(natsPublisher, a.logger, crawlerManager, trackRepo)
	natsConsumer := tracknats.NewTrackConsumer(a.natsConn, js, a.logger, a.config.Nats, trackingUC)
	trackHandler := trackhttp.NewTrackHandler(a.logger, trackingUC)
	restServer := httpserver.NewOpenAPIServer(a.logger, trackHandler, a.config.HTTPServer)
	grpcServer := grpcserver.NewGRPCServer(a.logger, trackingUC)

	// Start nats consumer
	go func() {
		err = natsConsumer.StartQueueReceiveMessages(a.config.Nats.Subject, a.config.Nats.DurableName)
		if err != nil {
			a.logger.Fatal(err)
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

	grpcServer.GracefulStop()
	shutdownErr := restServer.Shutdown()
	if shutdownErr != nil {
		return shutdownErr
	}

	return nil
}
