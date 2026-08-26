package feedstate

import "errors"

var ErrStore = errors.New("stream store failed")
var ErrRollback = errors.New("stream rollback failed")

type Store struct {
	Saved       []string
	StateEvents []string
	Fail        bool
}

func (s *Store) Save(value string) error {
	if s.Fail {
		s.StateEvents = append(s.StateEvents, value)
		return ErrRollback
	}
	s.Saved = append(s.Saved, value)
	return nil
}
