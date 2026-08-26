package store

import (
	"rss-reader/internal/model"
)

func (s *MemoryStore) CreateFetchTask(t *model.FetchTask) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.fetchTasks[t.ID] = t
	return nil
}

func (s *MemoryStore) GetFetchTask(id string) (*model.FetchTask, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.fetchTasks[id]
	if !ok {
		return nil, ErrNotFound
	}
	return t, nil
}

func (s *MemoryStore) ListFetchTasks() []*model.FetchTask {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.FetchTask, 0, len(s.fetchTasks))
	for _, t := range s.fetchTasks {
		list = append(list, t)
	}
	return list
}

func (s *MemoryStore) UpdateFetchTask(t *model.FetchTask) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.fetchTasks[t.ID]; !ok {
		return ErrNotFound
	}
	s.fetchTasks[t.ID] = t
	return nil
}

func (s *MemoryStore) DeleteFetchTask(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.fetchTasks[id]; !ok {
		return ErrNotFound
	}
	delete(s.fetchTasks, id)
	return nil
}
