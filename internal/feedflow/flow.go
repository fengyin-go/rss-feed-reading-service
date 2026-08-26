package feedflow

import "rss-reader/internal/feedstate"

type Observer struct{ Registry *feedstate.Registry }

func (o *Observer) Later(key string, ready chan<- struct{}, release <-chan struct{}) <-chan feedstate.Member {
	out := make(chan feedstate.Member, 1)
	go func() { ready <- struct{}{}; <-release; out <- *o.Registry.Snapshot(key) }()
	return out
}
