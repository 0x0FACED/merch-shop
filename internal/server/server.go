package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/0x0FACED/merch-shop/config"
	"github.com/0x0FACED/merch-shop/internal/database"
	"github.com/0x0FACED/merch-shop/internal/server/handler"
	"github.com/0x0FACED/merch-shop/internal/server/validator"
	"github.com/0x0FACED/merch-shop/internal/service"
	"github.com/0x0FACED/merch-shop/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	echo   *echo.Echo
	config config.ServerConfig
	logger *logger.ZapLogger
	db     *database.Postgres
}

func NewServer(cfg *config.ServiceConfig) (*Server, error) {
	log := logger.New(cfg.Logger)

	log.Info("Config is loaded")

	db, err := database.New(cfg.Database, log)
	if err != nil {
		return nil, fmt.Errorf("error creating database instance: %w", err)
	}
	db.MustConnect(context.Background())

	userService := service.NewUserService(db, log)

	log.Info("UserService created")

	e := echo.New()
	e.Server = &http.Server{
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}
	e.Validator = validator.NewAPIValidator()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	h := handler.NewHandler(userService, log, &cfg.Server)
	h.SetupRoutes(e)

	return &Server{
		echo:   e,
		config: cfg.Server,
		logger: log,
		db:     db,
	}, nil
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("Starting server...")

	errChan := make(chan error, 1)
	go func() {
		addr := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)
		err := s.echo.Start(addr)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
		close(errChan)
	}()

	select {
	case <-ctx.Done():
		s.logger.Info("Interrupt received, shutting down...")
		return s.Shutdown()
	case err := <-errChan:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("server stopped with err: %w", err)
		}
		return nil
	}
}

func (s Server) Echo() *echo.Echo {
	return s.echo
}

func (s Server) Database() *pgxpool.Pool {
	return s.db.Pool()
}

func (s *Server) Shutdown() error {
	s.logger.Info("Shutting down the server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.echo.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	if err := s.db.Close(); err != nil {
		return fmt.Errorf("failed to close db: %w", err)
	}

	s.logger.Info("Server shutdown successfully")
	return nil
}
