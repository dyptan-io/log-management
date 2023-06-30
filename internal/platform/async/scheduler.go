package async

import (
	"context"
	"sync"
	"time"

	"golang.org/x/exp/slog"
)

type (
	// Scheduler is a struct that runs asynchronous jobs at a configured intervals.
	Scheduler struct {
		cancels   []context.CancelFunc
		cancelsMu sync.Mutex
		logger    *slog.Logger
	}
	asyncJobFn func(context.Context) error
)

// NewScheduler creates a new Scheduler.
func NewScheduler(log *slog.Logger) *Scheduler {
	return &Scheduler{
		logger: log,
	}
}

// Schedule runs provided function at a given interval.
func (p *Scheduler) Schedule(ctx context.Context, interval time.Duration, fn asyncJobFn) {
	p.cancelsMu.Lock()
	defer p.cancelsMu.Unlock()

	ctx, cancel := context.WithCancel(ctx)

	go func(ctx context.Context) {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := fn(ctx); err != nil {
					p.logger.Error(err.Error())
				}
			}
		}
	}(ctx)

	p.cancels = append(p.cancels, cancel)
}

// Close cancels all scheduled jobs.
func (p *Scheduler) Close() {
	p.cancelsMu.Lock()
	defer p.cancelsMu.Unlock()

	for _, cancel := range p.cancels {
		cancel()
	}

	p.cancels = nil
}
