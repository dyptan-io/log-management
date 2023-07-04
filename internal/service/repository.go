package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/diptanw/log-management/internal/platform/storage"
)

// ErrBadRequestID is an error when request ID is malformed.
var ErrBadRequestID = errors.New("ID cannot be empty")

// Repository is a struct that manipulates the Log entries.
type Repository struct {
	db *storage.InMemory[LogEntry]
}

// LogEntry is a struct that represents log entry data model on persistence level.
// It expands Log model with record ID.
type LogEntry struct {
	Id         string
	Message    string
	Severity   string
	Timestamp  time.Time
	Attributes map[string]any
}

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

// GetAll returns all received Log entries.
func (r Repository) GetAll() ([]LogEntry, error) {
	res, err := r.db.Find(func(value LogEntry) bool {
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
