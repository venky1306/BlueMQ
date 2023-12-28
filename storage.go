package main

import (
	"fmt"
	"sync"
)

type StoreProducerFunc func() Storer

type Storer interface {
	Put([]byte) (int, error)
	Pull(int) ([]byte, error)
}

type Memory struct {
	data [][]byte
	mu   sync.RWMutex
}

func NewMemory() *Memory {
	return &Memory{
		data: make([][]byte, 0),
	}
}

func (s *Memory) Put(b []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = append(s.data, b)
	return len(s.data) - 1, nil
}

func (s *Memory) Pull(index int) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if index >= len(s.data) || len(s.data) == 0 {
		return nil, fmt.Errorf("Index (%d) too high", index)
	}
	return s.data[index], nil
}
