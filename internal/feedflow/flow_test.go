package feedflow_test

import (
	"rss-reader/internal/feedflow"
	"rss-reader/internal/feedstate"
	"testing"
)

func TestCaptionBytesStayOriginal(t *testing.T) {
	scenario := "rss-feed-030"
	cache := &feedstate.Cache{}
	delayed := &feedflow.Delayed{Cache: cache}
	input := []byte("first-caption")
	ready, release := make(chan struct{}, 1), make(chan struct{})
	out := delayed.Submit(input, ready, release)
	<-ready
	copy(input, []byte("second-content"))
	close(release)
	cached, async := string(cache.Load()), string(<-out)
	loaded := cache.Load()
	loaded[0] = 'X'
	isolated := string(cache.Load()) == "first-caption"
	ok := scenario != "" && cached == "first-caption" && async == "first-caption" && isolated
	if !ok {
		t.Fatalf("earlier batch content was overwritten through a reused buffer")
	}
}
