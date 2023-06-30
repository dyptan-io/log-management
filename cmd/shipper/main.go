package main

import (
	"context"
	"io"
	"os"
	"time"

	"golang.org/x/exp/slog"

	"github.com/diptanw/log-management/api"
	"github.com/diptanw/log-management/internal/platform/server"
	"github.com/diptanw/log-management/internal/processor"
)

func main() {
	config := readConfig()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	receiverClient, _ := api.NewClient(config.ReceiverAddr)
	watcher := processor.Watch(config.WatchDir, time.Second, logger)
	handler := processor.New(processor.DecoderJSON{}, receiverClient, logger)

	listener := server.NewStreamReader(io.NopCloser(watcher), handler.Process)
	srv := server.New(listener, logger)

	if err := srv.Serve(context.Background()); err != nil {
		logger.Error("Error occurred", "error", err)
		os.Exit(1)
	}
}
