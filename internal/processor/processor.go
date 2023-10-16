package processor

import (
	"context"
	"fmt"
	"net/http"

	"github.com/diptanw/log-management/v2/api"
	"github.com/diptanw/log-management/v2/internal/platform/server"
)

type (
	// SourceDecoder is an interface for Log decoder.
	SourceDecoder interface {
		Decode(b []byte) (api.Log, error)
	}

	// Processor is a struct that processes and send log entries to receiver.
	Processor struct {
		decoder SourceDecoder
		client  *api.Client
	}
)

// New returns a new instance of Processor.
func New(encoder SourceDecoder, client *api.Client) Processor {
	return Processor{
		decoder: encoder,
		client:  client,
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
		return fmt.Errorf("sending entry to receiver: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server responded with unsuccessful status code: %d", resp.StatusCode)
	}

	return nil
}
