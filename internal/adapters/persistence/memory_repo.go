package persistence

import (
	"sync"
	"time"

	"github.com/anderson-oki/awdl0-disabler/internal/core/domain"
)

type MemoryEventRepo struct {
	events []domain.Event
	mu     sync.RWMutex
}

func NewMemoryEventRepo() *MemoryEventRepo {
	return &MemoryEventRepo{
		events: make([]domain.Event, 0),
	}
}

func (r *MemoryEventRepo) Add(event domain.Event) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events = append(r.events, event)
}

func (r *MemoryEventRepo) GetRecent(duration time.Duration) []domain.Event {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cutoff := time.Now().Add(-duration)
	var recent []domain.Event

	for _, e := range r.events {
		if e.Timestamp.After(cutoff) {
			recent = append(recent, e)
		}
	}
	return recent
}
