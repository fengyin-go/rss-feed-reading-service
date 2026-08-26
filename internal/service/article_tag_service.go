package service

import (
	"sort"
	"time"

	"rss-reader/internal/model"
	"rss-reader/internal/store"
	"rss-reader/pkg/idgen"
)

func (s *Service) CreateArticleTag(input model.ArticleTag) (*model.ArticleTag, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetArticle(input.ArticleID); err != nil {
		if err == store.ErrNotFound {
			return nil, model.NewValidationError("article_id", "文章不存在")
		}
		return nil, err
	}
	if _, err := s.store.GetTag(input.TagID); err != nil {
		if err == store.ErrNotFound {
			return nil, model.NewValidationError("tag_id", "标签不存在")
		}
		return nil, err
	}
	now := time.Now()
	at := &model.ArticleTag{
		ID:        idgen.Hex(),
		ArticleID: input.ArticleID,
		TagID:     input.TagID,
		CreatedAt: now,
	}
	if err := s.store.CreateArticleTag(at); err != nil {
		return nil, err
	}
	return at, nil
}

func (s *Service) GetArticleTag(id string) (*model.ArticleTag, error) {
	return s.store.GetArticleTag(id)
}

func (s *Service) ListArticleTags(filter model.ArticleTagFilter, page, size int) ([]*model.ArticleTag, int, error) {
	all := s.store.ListArticleTags()
	matched := make([]*model.ArticleTag, 0, len(all))
	for _, at := range all {
		if filter.Match(at) {
			matched = append(matched, at)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.ArticleTag{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) DeleteArticleTag(id string) error {
	return s.store.DeleteArticleTag(id)
}

// ArticlesByTag 获取某标签下的文章。
func (s *Service) ArticlesByTag(tagID string) ([]*model.Article, error) {
	if _, err := s.store.GetTag(tagID); err != nil {
		return nil, err
	}
	all := s.store.ListArticleTags()
	var articleIDs []string
	for _, at := range all {
		if at.TagID == tagID {
			articleIDs = append(articleIDs, at.ArticleID)
		}
	}
	var result []*model.Article
	for _, id := range articleIDs {
		a, err := s.store.GetArticle(id)
		if err != nil {
			continue
		}
		result = append(result, a)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].PublishedAt.After(result[j].PublishedAt)
	})
	return result, nil
}

// TagsByArticle 获取某文章的标签。
func (s *Service) TagsByArticle(articleID string) ([]*model.Tag, error) {
	if _, err := s.store.GetArticle(articleID); err != nil {
		return nil, err
	}
	all := s.store.ListArticleTags()
	var tagIDs []string
	for _, at := range all {
		if at.ArticleID == articleID {
			tagIDs = append(tagIDs, at.TagID)
		}
	}
	var result []*model.Tag
	for _, id := range tagIDs {
		t, err := s.store.GetTag(id)
		if err != nil {
			continue
		}
		result = append(result, t)
	}
	return result, nil
}
