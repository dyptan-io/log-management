package server

import (
	"context"
	"os"
	"os/signal"
	"time"

	"golang.org/x/exp/slog"
)

// Listener is an interface for server listener.
type Listener interface {
	ListenAndServe() error
	Shutdown(context.Context) error
}

// Server is a wrapper for server listener like HTTP server. It handles process termination and
// shutdown.
type Server struct {
	srv    Listener
	logger *slog.Logger
}

// New returns a new instance of Server.
func New(srv Listener, log *slog.Logger) Server {
	return Server{srv: srv, logger: log}
}

// Serve starts a new server listener and handles interrupt and termination
// signals.
func (s Server) Serve(ctx context.Context) error {
	errsCh := make(chan error)

	go func() {
		s.logger.Info("Starting server...")
		errsCh <- s.srv.ListenAndServe()
	}()

	var err error

	defer func() {
		// Wait for completion before exiting.
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		s.logger.Info("Shutting down server...")
		err = s.srv.Shutdown(ctx)
	}()

	signalCtx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	// Wait until an error or interrupt signal is received.
	select {
	case err = <-errsCh:
		return err
	case <-signalCtx.Done():
		return err
	}
}
