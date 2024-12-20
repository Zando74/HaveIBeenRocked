package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/application/controller"
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/application/repository/postgresql"
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/config"
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/logger"
)

func main() {
	cfg := config.ConfigSingleton.GetInstance()
	log := logger.LoggerSingleton.GetInstance()

	log.Info(logger.ApplicationStartedMessage)
	log.Info(logger.ConfigLoadedMessage, cfg)

	postgresql.Init()

	controller := &controller.Controller{}

	controller.Run()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh
	log.Info(logger.ApplicationShutdownMessage)
	controller.Shutdown()
	postgresql.DBConnectionSingleton.Close()
	os.Exit(0)
}
