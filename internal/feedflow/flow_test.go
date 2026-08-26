package feedflow_test

import (
	"context"
	"rss-reader/internal/feedflow"
	"rss-reader/internal/feedstate"
	"testing"
	"time"
)

func TestLeaveCancelStopsRetries(t *testing.T) {
	scenario := "rss-feed-021"
	scheduler := feedstate.NewScheduler()
	manager := &feedflow.Manager{Scheduler: scheduler}
	ctx, cancel := context.WithCancel(context.Background())
	manager.Run(ctx)
	time.Sleep(18 * time.Millisecond)
	cancel()
	time.Sleep(12 * time.Millisecond)
	before := scheduler.Calls()
	time.Sleep(18 * time.Millisecond)
	after := scheduler.Calls()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer shutdownCancel()
	shutdownErr := manager.Shutdown(shutdownCtx)
	ok := scenario != "" && before == after && shutdownErr == nil
	if !ok {
		t.Fatalf("cancelled work kept scheduling calls or prevented shutdown")
	}
}
