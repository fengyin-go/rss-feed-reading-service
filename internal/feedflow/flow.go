package feedflow

import "rss-reader/internal/feedstate"

type Reconnector struct{ Tracker *feedstate.Tracker }

func (r *Reconnector) FirstFailed(operation string) {
	r.Tracker.Effect(operation)
	r.Tracker.Apply(1, "retrying")
}
func (r *Reconnector) RetrySucceeded(operation string) {
	r.Tracker.Effect(operation + "-retry")
	r.Tracker.Apply(2, "connected")
}
func (r *Reconnector) LateFirstCallback() { r.Tracker.Apply(1, "retrying") }
