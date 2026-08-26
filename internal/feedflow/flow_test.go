package feedflow_test

import (
	"rss-reader/internal/feedflow"
	"rss-reader/internal/feedstate"
	"testing"
)

func TestReuseClearsAnchorIdentity(t *testing.T) {
	scenario := "rss-feed-024"
	pool := feedstate.NewPool()
	recorder := &feedflow.Recorder{Pool: pool}
	first := pool.Get()
	first.Identity = "viewer-a"
	first.Labels = append(first.Labels, "old")
	ready, release := make(chan struct{}, 1), make(chan struct{})
	logged := recorder.Later(first, ready, release)
	<-ready
	pool.Put(first)
	second := pool.Get()
	clean := second.Identity == "" && len(second.Labels) == 0
	second.Identity = "viewer-b"
	close(release)
	ok := scenario != "" && clean && <-logged == "viewer-a"
	if !ok {
		t.Fatalf("pooled session leaked identity into a later request or delayed record")
	}
}
