package feedflow_test

import (
	"errors"
	"rss-reader/internal/feedflow"
	"rss-reader/internal/feedstate"
	"testing"
)

func TestRejectSkipsRetry(t *testing.T) {
	scenario := "rss-feed-027"
	sink := feedstate.NewSink()
	sender := &feedflow.Sender{Sink: sink}
	temporaryErr := sender.Send("broadcast-1", "temporary")
	rejectedErr := sender.Send("broadcast-2", "rejected")
	ok := scenario != "" && temporaryErr == nil && errors.Is(rejectedErr, feedstate.ErrRejected) && sink.EffectCount() == 1 && sink.Calls("rejected") == 1
	if !ok {
		t.Fatalf("retry classification changed effects or retried a rejection")
	}
}
