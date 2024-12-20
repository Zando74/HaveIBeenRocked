package http_test

import (
	"fmt"

	http_controller "github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/application/controller/http"
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/application/repository"
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/domain/usecase"
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/logger"
	"github.com/gin-gonic/gin"
)

var (
	LeakedHashRepositoryImpl = repository.NewHashLeakedRepositoryImpl()
)

func InitTestServer() *gin.Engine {
	log := logger.LoggerSingleton.GetInstance()

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(func(c *gin.Context) {
		log.Info(fmt.Sprintf("Request: %s %s", c.Request.Method, c.Request.URL.Path))
		c.Next()
	})

	httpController := &http_controller.CheckPasswordController{
		CheckPasswordUseCase: &usecase.CheckPasswordPresence{
			LeakedHashRepository: LeakedHashRepositoryImpl,
		},
	}

	router.POST("/api/check", httpController.CheckPasswordPresence)

	return router
}
