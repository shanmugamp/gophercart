package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"log/slog"

	"github.com/shanmugamp/gophercart/services/product-service/internal/config"
	"github.com/shanmugamp/gophercart/services/product-service/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	logger := config.NewLogger(
		cfg.App.LogLevel,
	)

	logger.Info(
		"starting service",
		slog.String("service", cfg.App.Name),
		slog.String("environment", cfg.App.Environment),
	)

	httpServer := server.NewHTTPServer(
		cfg.HTTP.Host,
		cfg.HTTP.Port,
		logger,
	)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go func() {
		logger.Info(
			"http server started",
			slog.String(
				"address",
				httpServer.Server.Addr,
			),
		)

		if err := httpServer.Server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {

			logger.Error(
				"http server failed",
				slog.Any("error", err),
			)

			stop()
		}
	}()

	<-ctx.Done()

	logger.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		cfg.HTTP.ShutdownTimeout,
	)
	defer cancel()

	if err := httpServer.Server.Shutdown(
		shutdownCtx,
	); err != nil {
		logger.Error(
			"graceful shutdown failed",
			slog.Any("error", err),
		)

		return
	}

	logger.Info("service stopped")
}
