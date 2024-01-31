package main

import (
	"net/http"
	"os"

	"github.com/apoldev/trackchecker/internal/app/config"
	usecase2 "github.com/apoldev/trackchecker/internal/app/crawler"
	repo2 "github.com/apoldev/trackchecker/internal/app/crawler/repo"
	trackhttp "github.com/apoldev/trackchecker/internal/app/track/delivery/http"
	tracknats "github.com/apoldev/trackchecker/internal/app/track/delivery/nats"
	"github.com/apoldev/trackchecker/internal/app/track/repo"
	"github.com/apoldev/trackchecker/internal/app/track/usecase"
	nats2 "github.com/apoldev/trackchecker/internal/pkg/nats"
	redis2 "github.com/apoldev/trackchecker/internal/pkg/redis"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{})
	logger.SetLevel(logrus.DebugLevel)

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

	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		logger.Fatal(err)
	}

	js.AddStream(&nats.StreamConfig{
		Name:     cfg.Nats.StreamName,
		Subjects: []string{cfg.Nats.Subject},
	})

	repoSpider := repo2.NewSpiderRepo(logger)
	err = repoSpider.LoadSpiders(cfg.ConfigSpiders)
	logger.Infof("Loaded %d spiders", len(repoSpider.Spiders))

	natsPublisher := tracknats.NewTrackPublisher(nc, js, logger, cfg.Nats)
	trackRepo := repo.NewTrackRepo(redisClient, logger)

	httpClient := http.DefaultClient
	crawlerManager := usecase2.NewCrawlerManager(repoSpider, logger, httpClient)
	trackingUC := usecase.NewTracking(natsPublisher, logger, crawlerManager, trackRepo)
	trackHandler := trackhttp.NewTrackHandler(logger, trackingUC)
	natsConsumer := tracknats.NewTrackConsumer(nc, js, logger, cfg.Nats, trackingUC)

	go func() {
		err := natsConsumer.StartQueueReceiveMessages(cfg.Nats.Subject, cfg.Nats.DurableName) //nolint:govet
		if err != nil {
			logger.Fatal(err)
		}
	}()

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.GET("/api/track", trackHandler.GetTrackingNumberResultHandler)
	r.POST("/api/track", trackHandler.TrackingNumberCrawlerHandler)

	r.Run(":" + cfg.HTTPServer.Port)
}
