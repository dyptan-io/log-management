package processor

import (
	"bytes"
	"context"
	"os"
	"path"
	"time"

	"golang.org/x/exp/slog"

	"github.com/diptanw/log-management/internal/platform/async"
)

// FileWatcher is a type for watching file changes in provided directory.
type FileWatcher struct {
	*bytes.Buffer
}

// Watch will traverse watch directory on schedule and will read new log entries into buffer.
func Watch(watchDir string, watchInterval time.Duration, logger *slog.Logger) FileWatcher {
	fileExists := make(map[string]struct{})
	buffer := bytes.NewBuffer(nil)
	scheduler := async.NewScheduler(logger)

	scheduler.Schedule(context.Background(), watchInterval, func(ctx context.Context) error {
		f, err := os.Open(watchDir)
		if err != nil {
			return err
		}

		defer f.Close()

		entries, err := f.ReadDir(-1)
		if err != nil {
			return err
		}

		for _, e := range entries {
			if e.IsDir() {
				continue
			}

			name := path.Join(watchDir, e.Name())

			if _, ok := fileExists[name]; ok {
				continue
			}

			fileExists[name] = struct{}{}

			b, err := os.ReadFile(name)
			if err != nil {
				return err
			}

			if _, err := buffer.Write(b); err != nil {
				return err
			}
		}

		return nil
	})

	return FileWatcher{
		Buffer: buffer,
	}
}
