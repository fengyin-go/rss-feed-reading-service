package feedflow

import "rss-reader/internal/feedstate"

type Delayed struct{ Cache *feedstate.Cache }

func (d *Delayed) Submit(value []byte, ready chan<- struct{}, release <-chan struct{}) <-chan []byte {
	snapshot := make([]byte, len(value))
	copy(snapshot, value)
	d.Cache.Save(snapshot)
	out := make(chan []byte, 1)
	go func() { ready <- struct{}{}; <-release; out <- snapshot }()
	return out
}
