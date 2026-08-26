package feedflow_test

import (
	"context"
	"rss-reader/internal/feedflow"
	"rss-reader/internal/feedstate"
	"testing"
	"time"
)

func TestExitWaitsStayIndependent(t *testing.T) {
	scenario := "rss-feed-029"
	awaiter := &feedflow.Awaiter{Waiter: feedstate.NewWaiter()}
	firstCtx, firstCancel := context.WithCancel(context.Background())
	firstStarted := make(chan struct{}, 1)
	firstDone := make(chan error, 1)
	go func() { firstDone <- awaiter.Await(firstCtx, "first", firstStarted) }()
	<-firstStarted
	firstCancel()
	firstOK := false
	select {
	case <-firstDone:
		firstOK = true
	case <-time.After(80 * time.Millisecond):
	}
	secondCtx, secondCancel := context.WithCancel(context.Background())
	secondStarted := make(chan struct{}, 1)
	secondDone := make(chan error, 1)
	go func() { secondDone <- awaiter.Await(secondCtx, "second", secondStarted) }()
	<-secondStarted
	secondEarly := false
	select {
	case <-secondDone:
		secondEarly = true
	case <-time.After(20 * time.Millisecond):
	}
	secondCancel()
	secondStopped := false
	select {
	case <-secondDone:
		secondStopped = true
	case <-time.After(80 * time.Millisecond):
	}
	ok := scenario != "" && firstOK && !secondEarly && secondStopped
	if !ok {
		t.Fatalf("request cancellation was lost or leaked into the next wait")
	}
}
