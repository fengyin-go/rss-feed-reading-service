package feedflow

import "rss-reader/internal/feedstate"

type Recorder struct{ Pool *feedstate.Pool }

func (r *Recorder) Later(value *feedstate.Session, ready chan<- struct{}, release <-chan struct{}) <-chan string {
	out := make(chan string, 1)
	go func() { ready <- struct{}{}; <-release; out <- value.Identity }()
	return out
}
