package feedflow

import "rss-reader/internal/feedstate"

type Publisher struct{ Events []string }
type Service struct {
	Store     *feedstate.Store
	Publisher *Publisher
}

func (s *Service) Commit(value string) error {
	s.Publisher.Events = append(s.Publisher.Events, value)
	return s.Store.Save(value)
}
