package feedflow

import (
	"context"
	"rss-reader/internal/feedstate"
)

type Manager struct{ Scheduler *feedstate.Scheduler }

func (m *Manager) Run(ctx context.Context) { m.Scheduler.Start(context.Background()) }
func (m *Manager) Shutdown(ctx context.Context) error {
	select {
	case <-m.Scheduler.Done():
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
