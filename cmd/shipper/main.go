package main

import (
	"context"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/dyptan-io/log-management/v2/api"
	"github.com/dyptan-io/log-management/v2/internal/platform/fs"
	"github.com/dyptan-io/log-management/v2/internal/platform/server"
	"github.com/dyptan-io/log-management/v2/internal/processor"
)

func main() {
	config := readConfig()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	receiverClient, _ := api.NewClient(config.ReceiverAddr)
	reader := fs.Watch(config.WatchDirs, time.Second, logger)
	handler := processor.New(processor.DecoderJSON{}, receiverClient)
	listener := server.NewStreamReader(io.NopCloser(reader), handler.Process)

	if err := server.New(listener, logger).Serve(context.Background()); err != nil {
		logger.Error("Error occurred", "error", err)
		os.Exit(1)
	}
}
