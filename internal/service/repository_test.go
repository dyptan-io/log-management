package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/dyptan-io/log-management/v2/internal/platform/storage"
)

func TestRepository_GetByID(t *testing.T) {
	store := storage.NewInMemory[LogEntry]()
	repo := NewRepository(store)

	testEntry := LogEntry{
		Id:         "test-123",
		Message:    "Test message",
		Severity:   "INFO",
		Timestamp:  time.Now(),
		Attributes: map[string]any{"user": "tester"},
	}

	require.NoError(t, store.Insert(testEntry))

	tests := map[string]struct {
		giveID    string
		wantEntry LogEntry
		wantErr   error
	}{
		"get existing entry": {
			giveID:    "test-123",
			wantEntry: testEntry,
		},
		"empty ID": {
			giveID:  "",
			wantErr: ErrBadRequestID,
		},
		"non-existent ID": {
			giveID:  "non-existent",
			wantErr: storage.ErrNotFound,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			entry, err := repo.GetByID(test.giveID)

			require.ErrorIs(t, err, test.wantErr)
			require.Equal(t, test.wantEntry.Id, entry.Id)
			require.Equal(t, test.wantEntry.Message, entry.Message)
			require.Equal(t, test.wantEntry.Severity, entry.Severity)
			require.Equal(t, test.wantEntry.Timestamp, entry.Timestamp)
			require.Equal(t, test.wantEntry.Attributes, entry.Attributes)
		})
	}
}

func TestRepository_Get(t *testing.T) {
	store := storage.NewInMemory[LogEntry]()
	repo := NewRepository(store)

	now := time.Now()
	pastTime := now.Add(-24 * time.Hour)
	futureTime := now.Add(24 * time.Hour)

	testEntries := []LogEntry{
		{
			Id:        "past",
			Message:   "Past message",
			Severity:  "INFO",
			Timestamp: pastTime,
		},
		{
			Id:        "present",
			Message:   "Present message",
			Severity:  "WARN",
			Timestamp: now,
		},
		{
			Id:        "future",
			Message:   "Future message",
			Severity:  "ERROR",
			Timestamp: futureTime,
		},
	}

	for _, entry := range testEntries {
		require.NoError(t, store.Insert(entry))
	}

	tests := map[string]struct {
		giveOpts    SearchOptions
		wantEntries []LogEntry
	}{
		"get all entries": {
			wantEntries: testEntries,
		},
		"filter by from time": {
			giveOpts: SearchOptions{
				From: &now,
			},
			wantEntries: []LogEntry{
				testEntries[1], // "present"
				testEntries[2], // "future"
			},
		},
		"filter by to time": {
			giveOpts: SearchOptions{
				To: &now,
			},
			wantEntries: []LogEntry{
				testEntries[0], // "past"
				testEntries[1], // "present"
			},
		},
		"filter by from and to time": {
			giveOpts: SearchOptions{
				From: &pastTime,
				To:   &futureTime,
			},
			wantEntries: []LogEntry{
				testEntries[0], // "past"
				testEntries[1], // "present"
				testEntries[2], // "future"
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			entries, err := repo.Get(test.giveOpts)

			require.NoError(t, err)
			require.ElementsMatch(t, test.wantEntries, entries)
		})
	}
}

func TestRepository_Create(t *testing.T) {
	store := storage.NewInMemory[LogEntry]()
	repo := NewRepository(store)

	tests := map[string]struct {
		giveEntry LogEntry
		wantErr   error
	}{
		"create valid entry": {
			giveEntry: LogEntry{
				Id:        "test-123",
				Message:   "Test message",
				Severity:  "INFO",
				Timestamp: time.Now(),
			},
		},
		"create entry with empty ID": {
			giveEntry: LogEntry{
				Message:   "No ID",
				Severity:  "WARN",
				Timestamp: time.Now(),
			},
			wantErr: storage.ErrMissingID,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			entry, err := repo.Create(test.giveEntry)

			require.ErrorIs(t, err, test.wantErr)
			require.Equal(t, test.giveEntry.Id, entry.Id)

			if test.wantErr != nil {
				return
			}

			storedEntry, err := repo.GetByID(entry.Id)

			require.NoError(t, err)
			require.Equal(t, entry.Id, storedEntry.Id)
		})
	}
}
