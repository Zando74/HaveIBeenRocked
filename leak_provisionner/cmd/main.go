package main

import (
	"github.com/Zando74/IHaveBeenRocked/leak_provisionner/internal/application/controller"
	"github.com/Zando74/IHaveBeenRocked/leak_provisionner/internal/config"
	"github.com/Zando74/IHaveBeenRocked/leak_provisionner/internal/logger"
)

func main() {
	cfg := config.ConfigSingleton.GetInstance()
	log := logger.LoggerSingleton.GetInstance()

	log.Info(logger.ApplicationStartedMessage)
	log.Info(logger.ConfigLoadedMessage, cfg)

	controller := &controller.Controller{}

	controller.RunCli()
}
