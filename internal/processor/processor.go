package processor

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/exp/slog"

	"github.com/diptanw/log-management/api"
	"github.com/diptanw/log-management/internal/platform/server"
)

type (
	// Decoder is an interface for Log decoder.
	Decoder interface {
		Decode(b []byte) (api.Log, error)
	}

	Source interface {
	}

	// Processor is a struct that processes and send log entries to receiver.
	Processor struct {
		decoder Decoder
		client  *api.Client
		logger  *slog.Logger
	}
)

// New returns a new instance of Processor.
func New(encoder Decoder, client *api.Client, logger *slog.Logger) Processor {
	return Processor{
		decoder: encoder,
		client:  client,
		logger:  logger,
	}
}

// Process decodes raw log entries sends them to logs receiver.
func (p Processor) Process(m server.Message) error {
	log, err := p.decoder.Decode(m.Data)
	if err != nil {
		return fmt.Errorf("decodig raw log entry: %w", err)
	}

	// For optimal performance, logs need to be sent in batches.
	resp, err := p.client.PostLog(context.Background(), []api.Log{log})
	if err != nil {
		return fmt.Errorf("sending entries to receiver: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		p.logger.Error("Server responded with unsuccessful status code", slog.Int("code", resp.StatusCode))
	}

	return nil
}
