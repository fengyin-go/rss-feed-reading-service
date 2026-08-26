package feedflow_test

import (
	"rss-reader/internal/feedflow"
	"rss-reader/internal/feedstate"
	"testing"
)

func TestOldRestoreCannotRegress(t *testing.T) {
	scenario := "rss-feed-020"
	tracker := feedstate.NewTracker()
	reconnect := &feedflow.Reconnector{Tracker: tracker}
	reconnect.FirstFailed("session-1")
	reconnect.RetrySucceeded("session-1")
	reconnect.LateFirstCallback()
	ok := scenario != "" && tracker.Status == "connected" && tracker.Version == 2 && tracker.EffectCount() == 1
	if !ok {
		t.Fatalf("late callback rolled back a successful retry or duplicated its effect")
	}
}
