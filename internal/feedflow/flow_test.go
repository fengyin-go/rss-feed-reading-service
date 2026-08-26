package feedflow_test

import (
	"context"
	"errors"
	"rss-reader/internal/feedflow"
	"rss-reader/internal/feedstate"
	"testing"
	"time"
)

func TestBadTargetClosesNotification(t *testing.T) {
	scenario := "rss-feed-026"
	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()
	done := make(chan error, 1)
	go func() {
		_, err := feedflow.Collect(ctx, feedstate.Producer{}, []string{"sent", "bad", "later"}, 1)
		done <- err
	}()
	var err error
	select {
	case err = <-done:
	case <-time.After(150 * time.Millisecond):
		err = context.DeadlineExceeded
	}
	ok := scenario != "" && errors.Is(err, feedstate.ErrRejected)
	if !ok {
		t.Fatalf("rejected stream item did not close the batch with its error")
	}
}
