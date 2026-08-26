package feedstate

import (
	"context"
	"sync"
	"time"
)

type Scheduler struct {
	mu     sync.Mutex
	calls  int
	done   chan struct{}
	cancel context.CancelFunc
}

func NewScheduler() *Scheduler { return &Scheduler{done: make(chan struct{})} }

func (s *Scheduler) Start(ctx context.Context) {
	s.mu.Lock()
	// 已在运行则忽略重复启动，避免泄漏后台 goroutine。
	if s.cancel != nil {
		s.mu.Unlock()
		return
	}
	ctx, cancel := context.WithCancel(ctx)
	s.cancel = cancel
	s.mu.Unlock()

	go func() {
		ticker := time.NewTicker(5 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				s.finish()
				return
			case <-ticker.C:
				s.mu.Lock()
				s.calls++
				s.mu.Unlock()
			}
		}
	}()
}

// finish 标记调度已停止并关闭 done，幂等。
func (s *Scheduler) finish() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.cancel == nil {
		return
	}
	s.cancel = nil
	close(s.done)
}

// Stop 主动取消后台调度并在 ctx 允许的范围内等待其退出。
// 取消（离开超时）后不再产生清理调用；未启动或已停止时立即返回 nil。
func (s *Scheduler) Stop(ctx context.Context) error {
	s.mu.Lock()
	cancel := s.cancel
	s.mu.Unlock()
	if cancel == nil {
		return nil
	}
	cancel()
	select {
	case <-s.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *Scheduler) Calls() int            { s.mu.Lock(); defer s.mu.Unlock(); return s.calls }
func (s *Scheduler) Done() <-chan struct{} { return s.done }
