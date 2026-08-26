package store

import (
	"rss-reader/internal/model"
)

func (s *MemoryStore) CreateFetchLog(l *model.FetchLog) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.fetchLogs[l.ID] = l
	return nil
}

func (s *MemoryStore) GetFetchLog(id string) (*model.FetchLog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	l, ok := s.fetchLogs[id]
	if !ok {
		return nil, ErrNotFound
	}
	return l, nil
}

func (s *MemoryStore) ListFetchLogs() []*model.FetchLog {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.FetchLog, 0, len(s.fetchLogs))
	for _, l := range s.fetchLogs {
		list = append(list, l)
	}
	return list
}

func (s *MemoryStore) DeleteFetchLog(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.fetchLogs[id]; !ok {
		return ErrNotFound
	}
	delete(s.fetchLogs, id)
	return nil
}
