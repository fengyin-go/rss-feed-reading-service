package feedstate

import (
	"context"
	"sync"
)

type Waiter struct {
	mu    sync.Mutex
	first context.Context
}

func NewWaiter() *Waiter { return &Waiter{} }
func (w *Waiter) Bind(key string, ctx context.Context) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.first == nil {
		w.first = ctx
	}
}
func (w *Waiter) Wait(key string, started chan<- struct{}) error {
	w.mu.Lock()
	ctx := w.first
	w.mu.Unlock()
	started <- struct{}{}
	<-ctx.Done()
	return ctx.Err()
}
