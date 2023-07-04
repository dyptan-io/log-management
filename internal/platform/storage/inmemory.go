// Package storage implements data storage types for the persistence layer
package storage

import (
	"errors"
	"sync"
)

var (
	// ErrMissingID is an error when the entity ID is missing.
	ErrMissingID = errors.New("missing record ID")
	// ErrNotFound is an error when record is not found.
	ErrNotFound = errors.New("record not found")
)

type (
	// ID is a type that represents record's unique identifier.
	ID string

	// Matcher is a func that filters dataset results.
	Matcher[T Record] func(value T) bool

	// Record is an interface that defines record's constraints.
	Record interface {
		ID() ID
	}

	// InMemory is a simple in-memory storage.
	InMemory[T Record] struct {
		records sync.Map
	}
)

// NewInMemory return a new instance on InMemory storage for a given type.
func NewInMemory[T Record]() *InMemory[T] {
	return &InMemory[T]{
		records: sync.Map{},
	}
}

// Get returns a single record matching the provided ID.
func (s *InMemory[T]) Get(id ID) (t T, err error) {
	if id == "" {
		return t, ErrMissingID
	}

	r, ok := s.records.Load(id)
	if !ok {
		return t, ErrNotFound
	}

	return r.(T), nil
}

// Find return a result of filtered records.
func (s *InMemory[T]) Find(match Matcher[T]) ([]T, error) {
	var res []T

	// Match by scanning the whole set (suboptimal).
	s.records.Range(func(k, v any) bool {
		if t, ok := v.(T); ok && match(t) {
			res = append(res, t)
		}

		return true
	})

	return res, nil
}

// Insert inserts an entry and overrides existing one.
func (s *InMemory[T]) Insert(r T) error {
	if r.ID() == "" {
		return ErrMissingID
	}

	s.records.Store(r.ID(), r)

	return nil
}
