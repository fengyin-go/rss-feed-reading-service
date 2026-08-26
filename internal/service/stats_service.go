package service

import (
	"sort"

	"rss-reader/internal/model"
)

// GlobalStats 全局统计。
type GlobalStats struct {
	FeedCount           int     `json:"feed_count"`
	ArticleCount        int     `json:"article_count"`
	CategoryCount       int     `json:"category_count"`
	UserCount           int     `json:"user_count"`
	TaskCount           int     `json:"task_count"`
	TaskSuccessRate     float64 `json:"task_success_rate"`
	UnreadArticleCount  int     `json:"unread_article_count"`
	StarredArticleCount int     `json:"starred_article_count"`
}

func (s *Service) GlobalStats() (*GlobalStats, error) {
	feeds := s.store.ListFeeds()
	articles := s.store.ListArticles()
	categories := s.store.ListCategories()
	users := s.store.ListUsers()
	tasks := s.store.ListFetchTasks()

	stats := &GlobalStats{
		FeedCount:     len(feeds),
		ArticleCount:  len(articles),
		CategoryCount: len(categories),
		UserCount:     len(users),
		TaskCount:     len(tasks),
	}

	if len(tasks) > 0 {
		successCount := 0
		for _, t := range tasks {
			if t.Status == model.FetchTaskSuccess {
				successCount++
			}
		}
		stats.TaskSuccessRate = float64(successCount) / float64(len(tasks))
	}

	for _, a := range articles {
		if !a.IsRead {
			stats.UnreadArticleCount++
		}
		if a.IsStarred {
			stats.StarredArticleCount++
		}
	}

	return stats, nil
}

// CategoryArticleCount 按分类统计文章数。
type CategoryArticleCount struct {
	Category string `json:"category"`
	Count    int    `json:"count"`
}

func (s *Service) ArticlesByCategory() ([]CategoryArticleCount, error) {
	feeds := s.store.ListFeeds()
	articles := s.store.ListArticles()

	feedCategory := make(map[string]string)
	for _, f := range feeds {
		feedCategory[f.ID] = f.Category
	}

	counts := make(map[string]int)
	for _, a := range articles {
		cat := feedCategory[a.FeedID]
		if cat == "" {
			cat = "未分类"
		}
		counts[cat]++
	}

	var result []CategoryArticleCount
	for cat, cnt := range counts {
		result = append(result, CategoryArticleCount{Category: cat, Count: cnt})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})
	return result, nil
}

// FeedArticleCount 订阅源文章数统计。
type FeedArticleCount struct {
	FeedID    string `json:"feed_id"`
	FeedTitle string `json:"feed_title"`
	Count     int    `json:"count"`
}

func (s *Service) ArticlesByFeedTopN(n int) ([]FeedArticleCount, error) {
	feeds := s.store.ListFeeds()
	articles := s.store.ListArticles()

	feedTitle := make(map[string]string)
	for _, f := range feeds {
		feedTitle[f.ID] = f.Title
	}

	counts := make(map[string]int)
	for _, a := range articles {
		counts[a.FeedID]++
	}

	var result []FeedArticleCount
	for feedID, cnt := range counts {
		result = append(result, FeedArticleCount{FeedID: feedID, FeedTitle: feedTitle[feedID], Count: cnt})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})
	if n > 0 && n < len(result) {
		result = result[:n]
	}
	return result, nil
}

// FeedFetchRate 每订阅源抓取成功率。
type FeedFetchRate struct {
	FeedID    string  `json:"feed_id"`
	FeedTitle string  `json:"feed_title"`
	Total     int     `json:"total"`
	Success   int     `json:"success"`
	Rate      float64 `json:"rate"`
}

func (s *Service) FetchRateByFeed() ([]FeedFetchRate, error) {
	feeds := s.store.ListFeeds()
	tasks := s.store.ListFetchTasks()

	feedTitle := make(map[string]string)
	for _, f := range feeds {
		feedTitle[f.ID] = f.Title
	}

	type pair struct{ total, success int }
	counts := make(map[string]pair)
	for _, t := range tasks {
		p := counts[t.FeedID]
		p.total++
		if t.Status == model.FetchTaskSuccess {
			p.success++
		}
		counts[t.FeedID] = p
	}

	var result []FeedFetchRate
	for feedID, p := range counts {
		rate := 0.0
		if p.total > 0 {
			rate = float64(p.success) / float64(p.total)
		}
		result = append(result, FeedFetchRate{
			FeedID:    feedID,
			FeedTitle: feedTitle[feedID],
			Total:     p.total,
			Success:   p.success,
			Rate:      rate,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Rate > result[j].Rate
	})
	return result, nil
}
