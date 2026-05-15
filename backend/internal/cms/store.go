package cms

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Store struct {
	mu      sync.RWMutex
	path    string
	content Content
}

func NewStore(path string) *Store {
	return &Store{path: path}
}

func (s *Store) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := os.MkdirAll(filepath.Dir(s.path), 0755); err != nil {
		return err
	}

	now := time.Now().UTC()
	if _, err := os.Stat(s.path); errors.Is(err, os.ErrNotExist) {
		s.content = DefaultContent(now)
		return s.saveLocked()
	}

	data, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &s.content); err != nil {
		return err
	}

	s.content.Normalize(now)
	return s.content.Validate()
}

func (s *Store) Snapshot() Content {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return cloneContent(s.content)
}

func (s *Store) Replace(next Content) (Content, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	next.Normalize(time.Now().UTC())
	if err := next.Validate(); err != nil {
		return Content{}, err
	}
	s.content = next
	if err := s.saveLocked(); err != nil {
		return Content{}, err
	}
	return cloneContent(s.content), nil
}

func (s *Store) Update(mutator func(*Content) error) (Content, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	next := cloneContent(s.content)
	if err := mutator(&next); err != nil {
		return Content{}, err
	}
	next.Normalize(time.Now().UTC())
	if err := next.Validate(); err != nil {
		return Content{}, err
	}
	s.content = next
	if err := s.saveLocked(); err != nil {
		return Content{}, err
	}
	return cloneContent(s.content), nil
}

func (s *Store) saveLocked() error {
	data, err := json.MarshalIndent(s.content, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')

	if err := os.MkdirAll(filepath.Dir(s.path), 0755); err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}

func cloneContent(content Content) Content {
	data, err := json.Marshal(content)
	if err != nil {
		return content
	}
	var clone Content
	if err := json.Unmarshal(data, &clone); err != nil {
		return content
	}
	return clone
}
