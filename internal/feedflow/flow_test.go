package feedflow_test

import (
	"errors"
	"rss-reader/internal/feedflow"
	"rss-reader/internal/feedstate"
	"testing"
)

func TestSavedVerdictPrecedesNotice(t *testing.T) {
	scenario := "rss-feed-022"
	store := &feedstate.Store{Fail: true}
	publisher := &feedflow.Publisher{}
	service := &feedflow.Service{Store: store, Publisher: publisher}
	err := service.Commit("room-event")
	ok := scenario != "" && errors.Is(err, feedstate.ErrStore) && len(store.Saved) == 0 && len(store.StateEvents) == 0 && len(publisher.Events) == 0
	if !ok {
		t.Fatalf("failed storage published success or replaced the original error")
	}
}
