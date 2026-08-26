package feedflow_test

import (
	"errors"
	"rss-reader/internal/feedflow"
	"rss-reader/internal/feedstate"
	"testing"
)

func TestMissingCaptionRuleFails(t *testing.T) {
	scenario := "rss-feed-028"
	missing := &feedflow.Checker{Policy: feedstate.LoadPolicy(false)}
	allowed, missingErr := missing.Check("caption")
	configured := &feedflow.Checker{Policy: feedstate.LoadPolicy(true)}
	addErr := configured.Add("caption")
	configuredAllowed, checkErr := configured.Check("caption")
	ok := scenario != "" && !allowed && errors.Is(missingErr, feedflow.ErrNoPolicy) && addErr == nil && checkErr == nil && configuredAllowed
	if !ok {
		t.Fatalf("missing policy bypassed checks or default policy could not accept a rule")
	}
}
