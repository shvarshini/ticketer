package lock

import (
	"fmt"
	"sync"
)

type LockService interface {
	TryLock(showSeatID string, userID string) error
	Unlock(showSeatID string) error
}

type InMemoryLockService struct {
	mu    sync.Mutex
	locks map[string]string
}

func NewInMemoryLockService() *InMemoryLockService {
	return &InMemoryLockService{
		locks: make(map[string]string),
	}
}

func (s *InMemoryLockService) TryLock(showSeatID string, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.locks[showSeatID]; exists {
		return fmt.Errorf("showSeat %s is already locked", showSeatID)
	}

	s.locks[showSeatID] = userID
	return nil
}

func (s *InMemoryLockService) Unlock(showSeatID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.locks, showSeatID)
	return nil
}
