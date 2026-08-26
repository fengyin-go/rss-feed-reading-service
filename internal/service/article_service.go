package service

import (
	"sort"
	"time"

	"rss-reader/internal/model"
	"rss-reader/internal/store"
	"rss-reader/pkg/idgen"
)

func (s *Service) CreateArticle(input model.Article) (*model.Article, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetFeed(input.FeedID); err != nil {
		if err == store.ErrNotFound {
			return nil, model.NewValidationError("feed_id", "所属订阅源不存在")
		}
		return nil, err
	}
	if _, err := s.store.GetArticleByGUID(input.GUID); err == nil {
		return nil, store.ErrConflict
	}
	now := time.Now()
	a := &model.Article{
		ID:          idgen.Hex(),
		FeedID:      input.FeedID,
		GUID:        input.GUID,
		Title:       input.Title,
		URL:         input.URL,
		Summary:     input.Summary,
		Content:     input.Content,
		Author:      input.Author,
		PublishedAt: input.PublishedAt,
		CreatedAt:   now,
	}
	if err := s.store.CreateArticle(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) GetArticle(id string) (*model.Article, error) {
	return s.store.GetArticle(id)
}

func (s *Service) ListArticles(filter model.ArticleFilter, page, size int) ([]*model.Article, int, error) {
	all := s.store.ListArticles()
	matched := make([]*model.Article, 0, len(all))
	for _, a := range all {
		if filter.Match(a) {
			matched = append(matched, a)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].PublishedAt.After(matched[j].PublishedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Article{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateArticle(id string, input model.Article) (*model.Article, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	a, err := s.store.GetArticle(id)
	if err != nil {
		return nil, err
	}
	a.Title = input.Title
	a.URL = input.URL
	a.Summary = input.Summary
	a.Content = input.Content
	a.Author = input.Author
	a.PublishedAt = input.PublishedAt
	if err := s.store.UpdateArticle(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) DeleteArticle(id string) error {
	return s.store.DeleteArticle(id)
}

func (s *Service) MarkArticleRead(id string) (*model.Article, error) {
	a, err := s.store.GetArticle(id)
	if err != nil {
		return nil, err
	}
	a.IsRead = true
	if err := s.store.UpdateArticle(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) MarkArticleStarred(id string, starred bool) (*model.Article, error) {
	a, err := s.store.GetArticle(id)
	if err != nil {
		return nil, err
	}
	a.IsStarred = starred
	if err := s.store.UpdateArticle(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) ExportArticlesByFeed(feedID string) ([]*model.Article, error) {
	if _, err := s.store.GetFeed(feedID); err != nil {
		return nil, err
	}
	all := s.store.ListArticles()
	var result []*model.Article
	for _, a := range all {
		if a.FeedID == feedID {
			result = append(result, a)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].PublishedAt.After(result[j].PublishedAt)
	})
	return result, nil
}
