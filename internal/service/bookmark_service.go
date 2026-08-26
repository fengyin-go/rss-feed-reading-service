package service

import (
	"sort"
	"time"

	"rss-reader/internal/model"
	"rss-reader/internal/store"
	"rss-reader/pkg/idgen"
)

func (s *Service) CreateBookmark(input model.Bookmark) (*model.Bookmark, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetUser(input.UserID); err != nil {
		if err == store.ErrNotFound {
			return nil, model.NewValidationError("user_id", "用户不存在")
		}
		return nil, err
	}
	if _, err := s.store.GetArticle(input.ArticleID); err != nil {
		if err == store.ErrNotFound {
			return nil, model.NewValidationError("article_id", "文章不存在")
		}
		return nil, err
	}
	now := time.Now()
	b := &model.Bookmark{
		ID:        idgen.Hex(),
		UserID:    input.UserID,
		ArticleID: input.ArticleID,
		CreatedAt: now,
	}
	if err := s.store.CreateBookmark(b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *Service) GetBookmark(id string) (*model.Bookmark, error) {
	return s.store.GetBookmark(id)
}

func (s *Service) ListBookmarks(filter model.BookmarkFilter, page, size int) ([]*model.Bookmark, int, error) {
	all := s.store.ListBookmarks()
	matched := make([]*model.Bookmark, 0, len(all))
	for _, b := range all {
		if filter.Match(b) {
			matched = append(matched, b)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Bookmark{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) DeleteBookmark(id string) error {
	return s.store.DeleteBookmark(id)
}
