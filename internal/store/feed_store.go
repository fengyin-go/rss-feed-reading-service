package store

import (
	"rss-reader/internal/model"
)

func (s *MemoryStore) CreateFeed(f *model.Feed) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.feeds {
		if exist.URL == f.URL {
			return ErrConflict
		}
	}
	s.feeds[f.ID] = f
	return nil
}

func (s *MemoryStore) GetFeed(id string) (*model.Feed, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	f, ok := s.feeds[id]
	if !ok {
		return nil, ErrNotFound
	}
	return f, nil
}

func (s *MemoryStore) GetFeedByURL(url string) (*model.Feed, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, f := range s.feeds {
		if f.URL == url {
			return f, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListFeeds() []*model.Feed {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Feed, 0, len(s.feeds))
	for _, f := range s.feeds {
		list = append(list, f)
	}
	return list
}

func (s *MemoryStore) UpdateFeed(f *model.Feed) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.feeds[f.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.feeds {
		if exist.ID != f.ID && exist.URL == f.URL {
			return ErrConflict
		}
	}
	s.feeds[f.ID] = f
	return nil
}

func (s *MemoryStore) DeleteFeed(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.feeds[id]; !ok {
		return ErrNotFound
	}
	delete(s.feeds, id)
	return nil
}
