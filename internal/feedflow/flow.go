package feedflow

import (
	"fmt"
	"rss-reader/internal/feedstate"
)

type Sender struct{ Sink *feedstate.Sink }

func (s *Sender) Send(operation, mode string) error {
	err := s.Sink.Deliver(operation+"-first", mode)
	if err != nil {
		return s.Sink.Deliver(fmt.Sprintf("%s-retry", operation), mode)
	}
	return nil
}
