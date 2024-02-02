package main

import (
	"os"

	"github.com/apoldev/trackchecker/internal/app/config"
	"github.com/apoldev/trackchecker/internal/pkg/app"
	nats2 "github.com/apoldev/trackchecker/internal/pkg/nats"
	"github.com/apoldev/trackchecker/internal/pkg/redis"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.DebugLevel)

	cfg, err := config.LoadConfig(os.Getenv("CONFIG_FILE"))
	if err != nil {
		logger.Fatal(err)
	}

	redisConn := redis.NewRedisConnection(&cfg.Redis)

	nc, err := nats2.NewNatsConn(cfg.Nats.Server)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("Connected to nats " + nc.ConnectedUrl())

	trackCheckerApp := app.New(logger, cfg, redisConn, nc)

	runErr := trackCheckerApp.Run()
	if runErr != nil {
		logger.Fatal(runErr)
	}
}
