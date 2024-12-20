package controller

import (
	"fmt"
	"net"
	"time"

	grpc_controller "github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/application/controller/grpc"
	http_controller "github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/application/controller/http"
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/application/grpc_proto"
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/application/repository"
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/config"
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/domain/usecase"
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	cfg                      = config.ConfigSingleton.GetInstance()
	log                      = logger.LoggerSingleton.GetInstance()
	LeakedHashRepositoryImpl = repository.NewHashLeakedRepositoryImpl()
)

type Controller struct {
	grpcServer *grpc.Server
	httpServer *gin.Engine
}

func (c *Controller) Run() {
	go c.RunGRPCServer()
	go c.RunHTTPServer()
}

func (c *Controller) Shutdown() {
	log.Info(logger.ShutdownHTTPServerMessage)
	log.Info(logger.ShutdownGRPCServerMessage)
	c.grpcServer.Stop()
}

func (c *Controller) RunGRPCServer() {

	log.Info(logger.RunGRPCServerMessage)

	c.grpcServer = grpc.NewServer()
	reflection.Register(c.grpcServer)

	grpc_proto.RegisterRawPasswordListUploadServer(c.grpcServer, &grpc_controller.PasswordListUploadServer{
		ProcessPasswordBatchUseCase: &usecase.ProcessPasswordBatch{
			HashPasswordBatchUseCase: &usecase.HashPasswordBatch{
				LeakedHashRepository: LeakedHashRepositoryImpl,
			},
			HashLeakedRepository: LeakedHashRepositoryImpl,
		},
	})

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Grpc.Port))
	if err != nil {
		log.Fatal(err)
	}

	go func(g *grpc.Server) {
		if err := g.Serve(listen); err != nil {
			log.Fatal(err)
		}
	}(c.grpcServer)

}

func (c *Controller) RunHTTPServer() {

	log.Info(logger.RunHTTPServerMessage)
	gin.SetMode(gin.ReleaseMode)
	c.httpServer = gin.New()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"POST", "GET", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour

	c.httpServer.Use(cors.New(config))

	c.httpServer.Use(func(c *gin.Context) {
		log.Info(fmt.Sprintf("Request: %s %s", c.Request.Method, c.Request.URL.Path))
		c.Next()
	})
	c.httpServer.Use(gin.Recovery())

	httpController := &http_controller.CheckPasswordController{
		CheckPasswordUseCase: &usecase.CheckPasswordPresence{
			LeakedHashRepository: LeakedHashRepositoryImpl,
		},
	}
	c.httpServer.POST("/api/check", httpController.CheckPasswordPresence)

	go func(g *gin.Engine) {
		if err := c.httpServer.Run(fmt.Sprintf(":%d", cfg.Http.Port)); err != nil {
			log.Fatal(err)
		}
	}(c.httpServer)

}
