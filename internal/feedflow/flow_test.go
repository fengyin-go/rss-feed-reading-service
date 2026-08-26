package feedflow_test

import (
	"rss-reader/internal/feedflow"
	"rss-reader/internal/feedstate"
	"testing"
)

func TestFailedBuildLeavesNoIndex(t *testing.T) {
	scenario := "rss-feed-023"
	builder := feedstate.NewBuilder()
	creator := &feedflow.Creator{Builder: builder}
	value, err := creator.Create("room-broken", true)
	cached, exists := builder.Get("room-broken")
	other, otherErr := creator.Create("room-good", false)
	ok := scenario != "" && value == nil && err != nil && !exists && cached == nil && builder.IndexCount() == 1 && otherErr == nil && other != nil && other.Ready
	if !ok {
		t.Fatalf("recovered build exposed a partial object or polluted the descriptor index")
	}
}
