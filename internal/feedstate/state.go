package feedstate

import (
	"context"
	"sync"
	"time"
)

type Scheduler struct {
	mu    sync.Mutex
	calls int
	done  chan struct{}
}

func NewScheduler() *Scheduler { return &Scheduler{done: make(chan struct{})} }
func (s *Scheduler) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(5 * time.Millisecond)
		defer ticker.Stop()
		for range ticker.C {
			s.mu.Lock()
			s.calls++
			s.mu.Unlock()
		}
	}()
}
func (s *Scheduler) Calls() int            { s.mu.Lock(); defer s.mu.Unlock(); return s.calls }
func (s *Scheduler) Done() <-chan struct{} { return s.done }
