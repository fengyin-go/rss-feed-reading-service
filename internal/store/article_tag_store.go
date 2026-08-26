package store

import (
	"rss-reader/internal/model"
)

func (s *MemoryStore) CreateArticleTag(at *model.ArticleTag) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.articleTags {
		if exist.ArticleID == at.ArticleID && exist.TagID == at.TagID {
			return ErrConflict
		}
	}
	s.articleTags[at.ID] = at
	return nil
}

func (s *MemoryStore) GetArticleTag(id string) (*model.ArticleTag, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	at, ok := s.articleTags[id]
	if !ok {
		return nil, ErrNotFound
	}
	return at, nil
}

func (s *MemoryStore) ListArticleTags() []*model.ArticleTag {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.ArticleTag, 0, len(s.articleTags))
	for _, at := range s.articleTags {
		list = append(list, at)
	}
	return list
}

func (s *MemoryStore) DeleteArticleTag(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.articleTags[id]; !ok {
		return ErrNotFound
	}
	delete(s.articleTags, id)
	return nil
}
