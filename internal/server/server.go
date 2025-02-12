package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/0x0FACED/merch-shop/config"
	"github.com/0x0FACED/merch-shop/internal/server/handler"
	"github.com/0x0FACED/merch-shop/internal/server/validator"
	"github.com/0x0FACED/merch-shop/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	echo     *echo.Echo
	config   config.ServerConfig
	logger   *logger.ZapLogger
	handler  *handler.Handler
	pprofSrv *http.Server // сервер для профилирвоания
}

func NewServer(cfg *config.ServiceConfig, h *handler.Handler) (*Server, error) {
	log := logger.New(cfg.Logger)

	e := echo.New()
	e.Server = &http.Server{
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}
	e.Validator = validator.NewAPIValidator()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	h.SetupRoutes(e)

	// Сервер для pprof
	pprofMux := http.NewServeMux()
	pprofMux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	pprofMux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	pprofMux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	pprofMux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	pprofMux.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))

	pprofSrv := &http.Server{
		Addr:    ":6060",
		Handler: pprofMux,
	}

	return &Server{
		echo:     e,
		config:   cfg.Server,
		logger:   log,
		handler:  h,
		pprofSrv: pprofSrv,
	}, nil
}

func (s Server) Echo() *echo.Echo {
	return s.echo
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("Starting server...")

	errChan := make(chan error, 1)

	go func() {
		addr := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)
		if err := s.echo.Start(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
		close(errChan)
	}()

	go func() {
		s.logger.Info("Starting pprof server on :6060")
		if err := s.pprofSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("Pprof server error", err)
		}
	}()

	select {
	case <-ctx.Done():
		s.logger.Info("Received shutdown signal...")
		return nil
	case err := <-errChan:
		return err
	}
}

func (s *Server) Shutdown() error {
	s.logger.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.echo.Shutdown(ctx); err != nil {
		return err
	}
	if err := s.pprofSrv.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
