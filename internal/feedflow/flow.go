package feedflow

import "rss-reader/internal/feedstate"

type Publisher struct{ Events []string }
type Service struct {
	Store     *feedstate.Store
	Publisher *Publisher
}

func (s *Service) Commit(value string) error {
	// 先保存结论；保存失败时不得发送成功提醒，并向上返回原始失败。
	if err := s.Store.Save(value); err != nil {
		return err
	}
	s.Publisher.Events = append(s.Publisher.Events, value)
	return nil
}
