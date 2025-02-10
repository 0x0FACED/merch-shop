package server

import (
	"context"
	"net/http"

	"github.com/0x0FACED/merch-shop/config"
	"github.com/0x0FACED/merch-shop/internal/database"
	"github.com/0x0FACED/merch-shop/internal/server/handler"
	"github.com/0x0FACED/merch-shop/internal/server/validator"
	"github.com/0x0FACED/merch-shop/internal/service"
	"github.com/0x0FACED/merch-shop/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	echoInstance *echo.Echo
	handler      *handler.Handler

	logger *logger.ZapLogger

	config config.ServerConfig
}

func newServer(e *echo.Echo, h *handler.Handler, l *logger.ZapLogger, cfg config.ServerConfig) *Server {
	return &Server{
		echoInstance: e,
		handler:      h,
		logger:       l,
		config:       cfg,
	}
}

func StartHTTP() error {
	server, err := prepareServer()
	if err != nil {
		return err
	}

	addr := addr(server.config.Host, server.config.Port)
	return server.echoInstance.Start(addr)
}

func addr(host, port string) string {
	return host + ":" + port
}

func prepareServer() (*Server, error) {
	ctx := context.Background()

	cfg := config.MustLoad()

	logger := logger.New(cfg.Logger)

	logger.Info("Config successfully loaded")

	logger.Info("Logger created")

	db, err := database.New(cfg.Database, logger)
	if err != nil {
		return nil, err
	}

	// Паникует, если не удалось установить коннект.
	// Все же коннект к базе супер важен, поэтому надо падать с паникой
	db.MustConnect(ctx)

	u := service.NewUserService(db, logger)

	logger.Info("User Service created")

	e := echoInstanceWithConfig(cfg.Server)

	logger.Info("Echo instance created")

	h := handler.NewHandler(u, logger, &cfg.Server)

	logger.Info("Handler is created")

	server := newServer(e, h, logger, cfg.Server)

	logger.Info("Setting up API Handlers")

	server.setupAPIRoutes()

	return server, nil
}

func (s *Server) setupAPIRoutes() {
	s.handler.SetupRoutes(s.echoInstance)
}

func echoInstanceWithConfig(cfg config.ServerConfig) *echo.Echo {
	e := echo.New()
	e.Server = &http.Server{ // пока что такие настройки
		ReadTimeout:  cfg.ReadTimeout,  // из конфига в секундах
		WriteTimeout: cfg.WriteTimeout, // из конфига в секундах
		IdleTimeout:  cfg.IdleTimeout,  // из конфига в секундах
	}

	e.Validator = validator.NewAPIValidator()

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	return e
}
