package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/dyptan-io/log-management/v2/api"
	"github.com/dyptan-io/log-management/v2/internal/platform/server"
	"github.com/dyptan-io/log-management/v2/internal/platform/storage"
	"github.com/dyptan-io/log-management/v2/internal/service"
)

func main() {
	config := readConfig()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	store := storage.NewInMemory[service.LogEntry]()
	repository := service.NewRepository(store)
	router := http.NewServeMux()

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
