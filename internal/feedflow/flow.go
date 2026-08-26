package feedflow

import (
	"errors"
	"fmt"
	"rss-reader/internal/feedstate"
)

type Sender struct{ Sink *feedstate.Sink }

// Send delivers an operation through the sink.
//
// A transient failure is retried until it succeeds, with the failed attempt
// producing no side effect that a later success cannot compensate. A
// rejection is a permanent failure: it is never retried and Send returns
// immediately so callers can stop the flow.
func (s *Sender) Send(operation, mode string) error {
	err := s.Sink.Deliver(operation+"-first", mode)
	if err == nil {
		return nil
	}
	if errors.Is(err, feedstate.ErrRejected) {
		// Permanent failure — do not retry.
		return err
	}
	// Transient (or any other) failure — retry once. The first attempt left
	// no side effect, so a successful retry fully compensates the failure.
	return s.Sink.Deliver(fmt.Sprintf("%s-retry", operation), mode)
}
