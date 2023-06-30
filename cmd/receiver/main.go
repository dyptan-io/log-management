package main

import (
	"context"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/slog"

	"github.com/diptanw/log-management/api"
	"github.com/diptanw/log-management/internal/platform/server"
	"github.com/diptanw/log-management/internal/platform/storage"
	"github.com/diptanw/log-management/internal/service"
)

func main() {
	config := readConfig()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	store := storage.NewInMemory[service.LogEntry]()
	repository := service.NewRepository(store)
	router := chi.NewRouter()

	api.HandlerFromMux(service.NewServer(repository, logger), router)

	srv := server.New(&http.Server{
		Addr:    config.HTTPAddr,
		Handler: router,
	}, logger)

	if err := srv.Serve(context.Background()); err != nil {
		logger.Error("fatal error occurred", "error", err)
		os.Exit(1)
	}
}
