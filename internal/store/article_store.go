package store

import (
	"rss-reader/internal/model"
)

func (s *MemoryStore) CreateArticle(a *model.Article) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.articles {
		if exist.GUID == a.GUID {
			return ErrConflict
		}
	}
	s.articles[a.ID] = a
	return nil
}

func (s *MemoryStore) GetArticle(id string) (*model.Article, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.articles[id]
	if !ok {
		return nil, ErrNotFound
	}
	return a, nil
}

func (s *MemoryStore) GetArticleByGUID(guid string) (*model.Article, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, a := range s.articles {
		if a.GUID == guid {
			return a, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListArticles() []*model.Article {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Article, 0, len(s.articles))
	for _, a := range s.articles {
		list = append(list, a)
	}
	return list
}

func (s *MemoryStore) UpdateArticle(a *model.Article) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.articles[a.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.articles {
		if exist.ID != a.ID && exist.GUID == a.GUID {
			return ErrConflict
		}
	}
	s.articles[a.ID] = a
	return nil
}

func (s *MemoryStore) DeleteArticle(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.articles[id]; !ok {
		return ErrNotFound
	}
	delete(s.articles, id)
	return nil
}
