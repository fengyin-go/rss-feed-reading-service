package feedstate

import "errors"

var ErrStore = errors.New("stream store failed")

type Store struct {
	Saved       []string
	StateEvents []string
	Fail        bool
}

func (s *Store) Save(value string) error {
	if s.Fail {
		// 保存失败时不留下内部状态事件，并返回原始失败。
		return ErrStore
	}
	s.Saved = append(s.Saved, value)
	return nil
}
