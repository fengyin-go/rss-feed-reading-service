package store

import (
	"rss-reader/internal/model"
)

func (s *MemoryStore) CreateBookmark(b *model.Bookmark) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.bookmarks {
		if exist.UserID == b.UserID && exist.ArticleID == b.ArticleID {
			return ErrConflict
		}
	}
	s.bookmarks[b.ID] = b
	return nil
}

func (s *MemoryStore) GetBookmark(id string) (*model.Bookmark, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, ok := s.bookmarks[id]
	if !ok {
		return nil, ErrNotFound
	}
	return b, nil
}

func (s *MemoryStore) ListBookmarks() []*model.Bookmark {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Bookmark, 0, len(s.bookmarks))
	for _, b := range s.bookmarks {
		list = append(list, b)
	}
	return list
}

func (s *MemoryStore) DeleteBookmark(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.bookmarks[id]; !ok {
		return ErrNotFound
	}
	delete(s.bookmarks, id)
	return nil
}
