package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/dyptan-io/log-management/v2/internal/platform/storage"
)

// ErrBadRequestID is an error when request ID is malformed.
var ErrBadRequestID = errors.New("ID cannot be empty")

type (
	// Repository is a struct that manipulates the Log entries.
	Repository struct {
		db *storage.InMemory[LogEntry]
	}

	// LogEntry is a struct that represents log entry data model on persistence level.
	// It expands Log model with record ID.
	LogEntry struct {
		Id         string
		Message    string
		Severity   string
		Timestamp  time.Time
		Attributes map[string]any
	}

	SearchOptions struct {
		From *time.Time
		To   *time.Time
	}
)

// ID returns a storage record ID for the LogEntry to comply with model constants.
func (l LogEntry) ID() storage.ID {
	return storage.ID(l.Id)
}

// NewRepository creates a new instance of Repository type.
func NewRepository(db *storage.InMemory[LogEntry]) Repository {
	return Repository{
		db: db,
	}
}

// GetByID returns a Log entry for the given ID.
func (r Repository) GetByID(id string) (LogEntry, error) {
	if id == "" {
		return LogEntry{}, ErrBadRequestID
	}

	entry, err := r.db.Get(storage.ID(id))
	if err != nil {
		return LogEntry{}, fmt.Errorf("getting log entry: %w", err)
	}

	return entry, nil
}

// Get returns Log entries by search criteria.
func (r Repository) Get(opts SearchOptions) ([]LogEntry, error) {
	res, err := r.db.Find(func(value LogEntry) bool {
		if opts.From != nil && opts.From.After(value.Timestamp) {
			return false
		}

		if opts.To != nil && opts.To.Before(value.Timestamp) {
			return false
		}

		return true
	})

	if err != nil {
		return []LogEntry{}, fmt.Errorf("getting all log entries: %w", err)
	}

	return res, nil
}

// Create creates a new Log entity in the storage.
func (r Repository) Create(entry LogEntry) (LogEntry, error) {
	if err := r.db.Insert(entry); err != nil {
		return LogEntry{}, fmt.Errorf("inserting log entry: %w", err)
	}

	return entry, nil
}
