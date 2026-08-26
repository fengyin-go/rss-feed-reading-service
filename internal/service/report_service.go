package service

import (
	"sort"
	"time"

	"rss-reader/internal/model"
)

// DailyArticleCount 按日期统计文章数。
type DailyArticleCount struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

func (s *Service) ArticlesByDate(days int) ([]DailyArticleCount, error) {
	if days <= 0 {
		days = 30
	}
	articles := s.store.ListArticles()
	counts := make(map[string]int)
	now := time.Now()
	cutoff := now.AddDate(0, 0, -days)
	for _, a := range articles {
		if a.PublishedAt.After(cutoff) || a.PublishedAt.Equal(cutoff) {
			date := a.PublishedAt.Format("2006-01-02")
			counts[date]++
		}
	}
	var result []DailyArticleCount
	for d := 0; d < days; d++ {
		date := now.AddDate(0, 0, -d).Format("2006-01-02")
		result = append(result, DailyArticleCount{Date: date, Count: counts[date]})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Date < result[j].Date
	})
	return result, nil
}

// UserActivity 用户活跃度统计。
type UserActivity struct {
	UserID            string `json:"user_id"`
	Username          string `json:"username"`
	SubscriptionCount int    `json:"subscription_count"`
	BookmarkCount     int    `json:"bookmark_count"`
}

func (s *Service) UserActivity() ([]UserActivity, error) {
	users := s.store.ListUsers()
	subscriptions := s.store.ListSubscriptions()
	bookmarks := s.store.ListBookmarks()

	subCounts := make(map[string]int)
	for _, sub := range subscriptions {
		subCounts[sub.UserID]++
	}
	bmCounts := make(map[string]int)
	for _, b := range bookmarks {
		bmCounts[b.UserID]++
	}

	var result []UserActivity
	for _, u := range users {
		result = append(result, UserActivity{
			UserID:            u.ID,
			Username:          u.Username,
			SubscriptionCount: subCounts[u.ID],
			BookmarkCount:     bmCounts[u.ID],
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].SubscriptionCount+result[i].BookmarkCount > result[j].SubscriptionCount+result[j].BookmarkCount
	})
	return result, nil
}

// FeedHealth 订阅源健康度。
type FeedHealth struct {
	FeedID       string    `json:"feed_id"`
	FeedTitle    string    `json:"feed_title"`
	ArticleCount int       `json:"article_count"`
	LastFetched  time.Time `json:"last_fetched_at"`
	IsHealthy    bool      `json:"is_healthy"`
}

func (s *Service) FeedHealth() ([]FeedHealth, error) {
	feeds := s.store.ListFeeds()
	articles := s.store.ListArticles()
	tasks := s.store.ListFetchTasks()

	feedArticleCount := make(map[string]int)
	for _, a := range articles {
		feedArticleCount[a.FeedID]++
	}

	feedLastTask := make(map[string]model.FetchTask)
	for _, t := range tasks {
		if existing, ok := feedLastTask[t.FeedID]; !ok || t.CreatedAt.After(existing.CreatedAt) {
			feedLastTask[t.FeedID] = *t
		}
	}

	now := time.Now()
	var result []FeedHealth
	for _, f := range feeds {
		fh := FeedHealth{
			FeedID:       f.ID,
			FeedTitle:    f.Title,
			ArticleCount: feedArticleCount[f.ID],
			LastFetched:  f.LastFetchedAt,
			IsHealthy:    f.Status == model.FeedStatusActive,
		}
		if lastTask, ok := feedLastTask[f.ID]; ok {
			if lastTask.Status == model.FetchTaskFailed {
				fh.IsHealthy = false
			}
			if lastTask.CreatedAt.After(fh.LastFetched) {
				fh.LastFetched = lastTask.CreatedAt
			}
		}
		if !fh.LastFetched.IsZero() && now.Sub(fh.LastFetched) > 24*time.Hour {
			fh.IsHealthy = false
		}
		result = append(result, fh)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].FeedTitle < result[j].FeedTitle
	})
	return result, nil
}

// UnreadByFeed 每订阅源未读文章数。
type UnreadByFeed struct {
	FeedID    string `json:"feed_id"`
	FeedTitle string `json:"feed_title"`
	Unread    int    `json:"unread"`
}

func (s *Service) UnreadByFeed() ([]UnreadByFeed, error) {
	feeds := s.store.ListFeeds()
	articles := s.store.ListArticles()

	feedTitle := make(map[string]string)
	for _, f := range feeds {
		feedTitle[f.ID] = f.Title
	}

	counts := make(map[string]int)
	for _, a := range articles {
		if !a.IsRead {
			counts[a.FeedID]++
		}
	}

	var result []UnreadByFeed
	for feedID, cnt := range counts {
		result = append(result, UnreadByFeed{
			FeedID:    feedID,
			FeedTitle: feedTitle[feedID],
			Unread:    cnt,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Unread > result[j].Unread
	})
	return result, nil
}

// TaskDurationStats 任务耗时统计。
type TaskDurationStats struct {
	AvgDurationSec float64 `json:"avg_duration_sec"`
	MaxDurationSec float64 `json:"max_duration_sec"`
	MinDurationSec float64 `json:"min_duration_sec"`
}

func (s *Service) TaskDurationStats() (*TaskDurationStats, error) {
	tasks := s.store.ListFetchTasks()
	stats := &TaskDurationStats{}
	var total float64
	var count int
	for _, t := range tasks {
		if !t.StartedAt.IsZero() && !t.FinishedAt.IsZero() {
			dur := t.FinishedAt.Sub(t.StartedAt).Seconds()
			if count == 0 || dur > stats.MaxDurationSec {
				stats.MaxDurationSec = dur
			}
			if count == 0 || dur < stats.MinDurationSec {
				stats.MinDurationSec = dur
			}
			total += dur
			count++
		}
	}
	if count > 0 {
		stats.AvgDurationSec = total / float64(count)
	}
	return stats, nil
}
