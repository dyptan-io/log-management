package fs

import (
	"bytes"
	"context"
	"io"
	"os"
	"path"
	"time"

	"golang.org/x/exp/slog"

	"github.com/diptanw/log-management/internal/platform/async"
)

// Watch traverses the watch directory on schedule and reads new log entries into a common buffer.
func Watch(watchDirs []string, watchInterval time.Duration, logger *slog.Logger) *bytes.Buffer {
	bytesRead := make(map[string]int64)
	buffer := bytes.NewBuffer(nil)

	async.Schedule(context.Background(), watchInterval, func(ctx context.Context) error {
		for _, dir := range watchDirs {
			if err := scanDir(dir, bytesRead, buffer); err != nil {
				return err
			}
		}

		return nil
	}, logger)

	return buffer
}

func scanDir(watchDir string, bytesRead map[string]int64, w io.Writer) error {
	dir, err := os.Open(watchDir)
	if err != nil {
		return err
	}

	defer dir.Close()

	entries, err := dir.ReadDir(-1)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		fi, err := e.Info()
		if err != nil {
			return err
		}

		name := path.Join(watchDir, e.Name())
		size, ok := bytesRead[name]
		if ok && size >= fi.Size() {
			continue
		}

		bytesRead[name] = fi.Size()

		return readFrom(name, size, w)
	}

	return nil
}

func readFrom(filePath string, pos int64, w io.Writer) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.Seek(pos, 0)
	if err != nil {
		return err
	}

	b, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	if _, err := w.Write(b); err != nil {
		return err
	}

	return nil
}
