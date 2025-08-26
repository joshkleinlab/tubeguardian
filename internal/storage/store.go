package storage

import (
	"encoding/json"
	"os"
	"sync"
)

type Store struct {
	mu       sync.Mutex
	filePath string
	LastSeen string `json:"last_seen_comment_id"`
}

// NewStore initializes the store with a JSON file
func NewStore(path string) *Store {
	s := &Store{filePath: path}
	_ = s.load() // Try to load existing file
	return s
}

// load reads the last seen ID from disk
func (s *Store) load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.Open(s.filePath)
	if err != nil {
		// no file yet, ignore
		return nil
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(s)
}

// Save writes the current state to disk
func (s *Store) Save() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.Create(s.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(s)
}

// GetLastSeen returns the last seen comment ID
func (s *Store) GetLastSeen() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.LastSeen
}

// SetLastSeen updates the last seen comment ID and persists it
func (s *Store) SetLastSeen(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.LastSeen = id
	return s.Save()
}
