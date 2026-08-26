package service

import (
	"rss-reader/internal/model"
	"rss-reader/internal/store"
	"rss-reader/pkg/idgen"
	"time"
)

// BatchCreateArticlesResult 批量创建文章结果。
type BatchCreateArticlesResult struct {
	Created    int      `json:"created"`
	Skipped    int      `json:"skipped"`
	Errors     []string `json:"errors,omitempty"`
	ArticleIDs []string `json:"article_ids,omitempty"`
}

func (s *Service) BatchCreateArticles(feedID string, inputs []model.Article) (*BatchCreateArticlesResult, error) {
	if _, err := s.store.GetFeed(feedID); err != nil {
		if err == store.ErrNotFound {
			return nil, model.NewValidationError("feed_id", "订阅源不存在")
		}
		return nil, err
	}
	result := &BatchCreateArticlesResult{}
	for _, input := range inputs {
		input.FeedID = feedID
		if err := input.Validate(); err != nil {
			result.Errors = append(result.Errors, err.Error())
			continue
		}
		if _, err := s.store.GetArticleByGUID(input.GUID); err == nil {
			result.Skipped++
			continue
		}
		now := time.Now()
		a := &model.Article{
			ID:          idgen.Hex(),
			FeedID:      feedID,
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
			result.Errors = append(result.Errors, err.Error())
			continue
		}
		result.Created++
		result.ArticleIDs = append(result.ArticleIDs, a.ID)
	}
	return result, nil
}

// BatchMarkArticlesReadResult 批量标记已读结果。
type BatchMarkArticlesReadResult struct {
	Marked int      `json:"marked"`
	Errors []string `json:"errors,omitempty"`
}

func (s *Service) BatchMarkArticlesRead(articleIDs []string) (*BatchMarkArticlesReadResult, error) {
	result := &BatchMarkArticlesReadResult{}
	for _, id := range articleIDs {
		a, err := s.store.GetArticle(id)
		if err != nil {
			result.Errors = append(result.Errors, id+": "+err.Error())
			continue
		}
		a.IsRead = true
		if err := s.store.UpdateArticle(a); err != nil {
			result.Errors = append(result.Errors, id+": "+err.Error())
			continue
		}
		result.Marked++
	}
	return result, nil
}

// BatchStarArticlesResult 批量收藏结果。
type BatchStarArticlesResult struct {
	Starred int      `json:"starred"`
	Errors  []string `json:"errors,omitempty"`
}

func (s *Service) BatchStarArticles(articleIDs []string, starred bool) (*BatchStarArticlesResult, error) {
	result := &BatchStarArticlesResult{}
	for _, id := range articleIDs {
		a, err := s.store.GetArticle(id)
		if err != nil {
			result.Errors = append(result.Errors, id+": "+err.Error())
			continue
		}
		a.IsStarred = starred
		if err := s.store.UpdateArticle(a); err != nil {
			result.Errors = append(result.Errors, id+": "+err.Error())
			continue
		}
		result.Starred++
	}
	return result, nil
}

// BatchDeleteArticlesResult 批量删除结果。
type BatchDeleteArticlesResult struct {
	Deleted int      `json:"deleted"`
	Errors  []string `json:"errors,omitempty"`
}

func (s *Service) BatchDeleteArticles(articleIDs []string) (*BatchDeleteArticlesResult, error) {
	result := &BatchDeleteArticlesResult{}
	for _, id := range articleIDs {
		if err := s.store.DeleteArticle(id); err != nil {
			result.Errors = append(result.Errors, id+": "+err.Error())
			continue
		}
		result.Deleted++
	}
	return result, nil
}

// BatchCreateFeedsResult 批量创建订阅源结果。
type BatchCreateFeedsResult struct {
	Created int      `json:"created"`
	Errors  []string `json:"errors,omitempty"`
	FeedIDs []string `json:"feed_ids,omitempty"`
}

func (s *Service) BatchCreateFeeds(inputs []model.Feed) (*BatchCreateFeedsResult, error) {
	result := &BatchCreateFeedsResult{}
	for _, input := range inputs {
		if err := input.Validate(); err != nil {
			result.Errors = append(result.Errors, err.Error())
			continue
		}
		if _, err := s.store.GetFeedByURL(input.URL); err == nil {
			result.Errors = append(result.Errors, "URL 已存在: "+input.URL)
			continue
		}
		now := time.Now()
		f := &model.Feed{
			ID:            idgen.Hex(),
			Title:         input.Title,
			URL:           input.URL,
			Description:   input.Description,
			Category:      input.Category,
			Status:        input.Status,
			FetchInterval: input.FetchInterval,
			CreatedAt:     now,
			UpdatedAt:     now,
		}
		if err := s.store.CreateFeed(f); err != nil {
			result.Errors = append(result.Errors, err.Error())
			continue
		}
		result.Created++
		result.FeedIDs = append(result.FeedIDs, f.ID)
	}
	return result, nil
}
