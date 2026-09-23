package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shanmugamp/gophercart/services/product-service/internal/handler"
)

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/health", handler.HealthHandler)
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, nil),
	)

	logger.Info("starting product service")

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go func() {
		logger.Info("product service listening",
			slog.String("address", server.Addr),
		)

		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {

			logger.Error("server failed",
				slog.Any("error", err),
			)

			stop()
		}
	}()

	<-ctx.Done()

	logger.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed",
			slog.Any("error", err),
		)
	}

	logger.Info("product service stopped")
}
