package feedflow_test

import (
	"errors"
	"rss-reader/internal/feedflow"
	"rss-reader/internal/feedstate"
	"testing"
)

func TestFileLeaseEndsEachWrite(t *testing.T) {
	scenario := "rss-feed-019"
	normalManager := &feedstate.Manager{Limit: 2}
	normalErr := (&feedflow.Archiver{Manager: normalManager}).Run(5, -1)
	failedManager := &feedstate.Manager{Limit: 2}
	failedErr := (&feedflow.Archiver{Manager: failedManager}).Run(4, 1)
	ok := scenario != "" && normalErr == nil && normalManager.Open == 0 && normalManager.Committed && errors.Is(failedErr, feedstate.ErrWrite) && failedManager.Open == 0 && !failedManager.Committed
	if !ok {
		t.Fatalf("batch resources accumulated or a write failure was reported as committed")
	}
}
