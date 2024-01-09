package db_logger

import (
	"errors"
	"sync"
)

type Store struct {
	sync.RWMutex
	m map[string]string
}

// sentinel error
var ErrorNoSuchKey = errors.New("no such key")

func (s *Store) Put(key string, value string) error {
	s.Lock()
	s.m[key] = value
	s.Unlock()
	return nil
}

func (s *Store) Get(key string) (string, error) {
	s.RLock()
	value, ok := s.m[key]
	s.RUnlock()

	if !ok {
		return "", ErrorNoSuchKey
	}

	return value, nil
}

func (s *Store) Delete(key string) error {
	s.Lock()
	delete(s.m, key)
	s.Unlock()
	return nil
}

func NewStore() *Store {
	return &Store{m: make(map[string]string)}
}
