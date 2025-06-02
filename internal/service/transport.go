package service

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/dyptan-io/log-management/v2/api"
	"github.com/dyptan-io/log-management/v2/internal/platform/storage"
)

// Server implements the api.ServerInterface.
type Server struct {
	repo   Repository
	logger *slog.Logger
}

// NewServer return a new instance of Server.
func NewServer(repo Repository, logger *slog.Logger) Server {
	return Server{
		repo:   repo,
		logger: logger,
	}
}

func (Server) Health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (s Server) ListLogs(w http.ResponseWriter, _ *http.Request, params api.ListLogsParams) {
	entries, err := s.repo.Get(SearchOptions{
		From: params.From,
		To:   params.To,
	})
	if err != nil {
		s.handleError(w, err)
		return
	}

	logs := make([]api.Log, len(entries))

	for i, entry := range entries {
		logs[i] = toDto(entry)
	}

	writeJSON(w, logs)
}

func (s Server) PostLog(w http.ResponseWriter, r *http.Request) {
	var logs []api.Log

	if err := readJSON(r, &logs); err != nil {
		s.handleError(w, err)
		return
	}

	for _, log := range logs {
		if _, err := s.repo.Create(fromDto(log)); err != nil {
			s.handleError(w, err)
			return
		}
	}
}

func (s Server) GetLogsById(w http.ResponseWriter, _ *http.Request, id string) {
	entry, err := s.repo.GetByID(id)
	if err != nil {
		s.handleError(w, err)
		return
	}

	writeJSON(w, toDto(entry))
}

func (s Server) handleError(w http.ResponseWriter, err error) {
	if status := toStatus(err); status != http.StatusOK {
		s.logger.Warn("request has failed", "error", err)

		w.WriteHeader(status)
		writeJSON(w, api.ErrorResponse{Errors: []string{err.Error()}})
	}
}

func toStatus(err error) int {
	if err == nil {
		return http.StatusOK
	}

	if isError(err, storage.ErrMissingID, ErrBadRequestID) {
		return http.StatusBadRequest
	}

	if isError(err, storage.ErrNotFound) {
		return http.StatusNotFound
	}

	return http.StatusInternalServerError
}

func isError(err error, targets ...error) bool {
	for _, t := range targets {
		if errors.Is(err, t) {
			return true
		}
	}

	return false
}

func toDto(entry LogEntry) api.Log {
	return api.Log{
		Id:         entry.Id,
		Message:    entry.Message,
		Severity:   entry.Severity,
		Attributes: entry.Attributes,
		Timestamp:  entry.Timestamp,
	}
}

func fromDto(entry api.Log) LogEntry {
	return LogEntry{
		Id:         entry.Id,
		Message:    entry.Message,
		Severity:   entry.Severity,
		Attributes: entry.Attributes,
		Timestamp:  entry.Timestamp,
	}
}
