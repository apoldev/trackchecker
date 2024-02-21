package main

import (
	"os"

	"github.com/apoldev/trackchecker/internal/app/config"
	"github.com/apoldev/trackchecker/internal/pkg/app"
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

	trackCheckerApp := app.New(logger, cfg)
	runErr := trackCheckerApp.Run()
	if runErr != nil {
		logger.Fatal(runErr)
	}
}
