package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/shanmugamp/gophercart/services/product-service/internal/handler"
	"github.com/shanmugamp/gophercart/services/product-service/internal/middleware"
)

type HTTPServer struct {
	Server *http.Server
}

func NewHTTPServer(
	host string,
	port int,
	logger *slog.Logger,
) *HTTPServer {
	mux := http.NewServeMux()

	mux.HandleFunc(
		"/health",
		handler.Health,
	)

	mux.HandleFunc(
		"/ready",
		handler.Readiness,
	)

	apiMux := http.NewServeMux()

	apiMux.Handle(
		"/api/v1/",
		http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			http.NotFound(w, r)
		}),
	)

	mux.Handle(
		"/api/v1/",
		apiMux,
	)

	server := &http.Server{
		Addr: fmt.Sprintf(
			"%s:%d",
			host,
			port,
		),
		Handler: middleware.RequestID(mux),

		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	return &HTTPServer{
		Server: server,
	}
}
