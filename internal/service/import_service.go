package service

import (
	"encoding/json"
	"rss-reader/internal/model"
	"rss-reader/internal/store"
	"rss-reader/pkg/idgen"
	"time"
)

// ImportFeedsResult 导入订阅源结果。
type ImportFeedsResult struct {
	Imported int      `json:"imported"`
	Skipped  int      `json:"skipped"`
	Errors   []string `json:"errors,omitempty"`
}

// ImportFeeds 从 JSON 数据导入订阅源。
func (s *Service) ImportFeeds(data []byte) (*ImportFeedsResult, error) {
	var inputs []model.Feed
	if err := json.Unmarshal(data, &inputs); err != nil {
		return nil, model.NewValidationError("data", "JSON 解析失败: "+err.Error())
	}
	result := &ImportFeedsResult{}
	for _, input := range inputs {
		if err := input.Validate(); err != nil {
			result.Errors = append(result.Errors, err.Error())
			continue
		}
		if _, err := s.store.GetFeedByURL(input.URL); err == nil {
			result.Skipped++
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
		result.Imported++
	}
	return result, nil
}

// ImportArticlesResult 导入文章结果。
type ImportArticlesResult struct {
	Imported int      `json:"imported"`
	Skipped  int      `json:"skipped"`
	Errors   []string `json:"errors,omitempty"`
}

// ImportArticles 从 JSON 数据导入文章。
func (s *Service) ImportArticles(feedID string, data []byte) (*ImportArticlesResult, error) {
	if _, err := s.store.GetFeed(feedID); err != nil {
		if err == store.ErrNotFound {
			return nil, model.NewValidationError("feed_id", "订阅源不存在")
		}
		return nil, err
	}
	var inputs []model.Article
	if err := json.Unmarshal(data, &inputs); err != nil {
		return nil, model.NewValidationError("data", "JSON 解析失败: "+err.Error())
	}
	result := &ImportArticlesResult{}
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
		result.Imported++
	}
	return result, nil
}

// ExportFeeds 导出全部订阅源。
func (s *Service) ExportFeeds() ([]*model.Feed, error) {
	feeds := s.store.ListFeeds()
	result := make([]*model.Feed, len(feeds))
	copy(result, feeds)
	return result, nil
}

// ExportAll 导出全部数据汇总。
func (s *Service) ExportAll() (map[string]interface{}, error) {
	result := make(map[string]interface{})
	result["feeds"] = s.store.ListFeeds()
	result["articles"] = s.store.ListArticles()
	result["categories"] = s.store.ListCategories()
	result["users"] = s.store.ListUsers()
	result["fetch_tasks"] = s.store.ListFetchTasks()
	result["fetch_logs"] = s.store.ListFetchLogs()
	result["subscriptions"] = s.store.ListSubscriptions()
	result["bookmarks"] = s.store.ListBookmarks()
	result["tags"] = s.store.ListTags()
	result["article_tags"] = s.store.ListArticleTags()
	result["notifications"] = s.store.ListNotifications()
	return result, nil
}
