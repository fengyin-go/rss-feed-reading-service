package feedflow_test

import (
	"rss-reader/internal/feedflow"
	"rss-reader/internal/feedstate"
	"testing"
)

func TestMemberListIsSnapshot(t *testing.T) {
	scenario := "rss-feed-025"
	registry := feedstate.NewRegistry()
	registry.Put("room", feedstate.Member{Name: "viewer-a", Active: true})
	direct := registry.Snapshot("room")
	ready, release := make(chan struct{}, 1), make(chan struct{})
	delayed := (&feedflow.Observer{Registry: registry}).Later("room", ready, release)
	<-ready
	registry.Put("room", feedstate.Member{Name: "viewer-b", Active: false})
	close(release)
	later := <-delayed
	ok := scenario != "" && direct.Name == "viewer-a" && direct.Active && later.Name == "viewer-a" && later.Active
	if !ok {
		t.Fatalf("snapshot changed after later room update")
	}
}
