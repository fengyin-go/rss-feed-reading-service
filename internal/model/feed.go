package model

import (
	"strings"
	"time"
)

const (
	FeedStatusActive = "active"
	FeedStatusPaused = "paused"
)

// Feed 表示一个 RSS 订阅源。
type Feed struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	URL           string    `json:"url"`
	Description   string    `json:"description"`
	Category      string    `json:"category"`
	Status        string    `json:"status"`
	FetchInterval int       `json:"fetch_interval"`
	LastFetchedAt time.Time `json:"last_fetched_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (f *Feed) Validate() error {
	f.Title = strings.TrimSpace(f.Title)
	f.URL = strings.TrimSpace(f.URL)
	f.Category = strings.TrimSpace(f.Category)
	if f.Title == "" {
		return NewValidationError("title", "订阅源标题不能为空")
	}
	if f.URL == "" {
		return NewValidationError("url", "订阅源 URL 不能为空")
	}
	if f.Status == "" {
		f.Status = FeedStatusActive
	}
	if f.Status != FeedStatusActive && f.Status != FeedStatusPaused {
		return NewValidationError("status", "订阅源状态不合法")
	}
	if f.FetchInterval < 1 {
		f.FetchInterval = 60
	}
	return nil
}

type FeedFilter struct {
	Category string
	Status   string
	Keyword  string
}

func (f FeedFilter) Match(feed *Feed) bool {
	if f.Category != "" && feed.Category != f.Category {
		return false
	}
	if f.Status != "" && feed.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(feed.Title), k) &&
			!strings.Contains(strings.ToLower(feed.URL), k) {
			return false
		}
	}
	return true
}

func CanTransitionFeed(from, to string) bool {
	if from == FeedStatusActive && to == FeedStatusPaused {
		return true
	}
	if from == FeedStatusPaused && to == FeedStatusActive {
		return true
	}
	return false
}
