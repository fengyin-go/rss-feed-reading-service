package model

import (
	"strings"
	"time"
)

// Article 表示一篇 RSS 文章。
type Article struct {
	ID          string    `json:"id"`
	FeedID      string    `json:"feed_id"`
	GUID        string    `json:"guid"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Summary     string    `json:"summary"`
	Content     string    `json:"content"`
	Author      string    `json:"author"`
	PublishedAt time.Time `json:"published_at"`
	IsRead      bool      `json:"is_read"`
	IsStarred   bool      `json:"is_starred"`
	CreatedAt   time.Time `json:"created_at"`
}

func (a *Article) Validate() error {
	a.Title = strings.TrimSpace(a.Title)
	a.URL = strings.TrimSpace(a.URL)
	a.GUID = strings.TrimSpace(a.GUID)
	if a.FeedID == "" {
		return NewValidationError("feed_id", "所属订阅源不能为空")
	}
	if a.GUID == "" {
		return NewValidationError("guid", "文章 GUID 不能为空")
	}
	if a.Title == "" {
		return NewValidationError("title", "文章标题不能为空")
	}
	if a.URL == "" {
		return NewValidationError("url", "文章 URL 不能为空")
	}
	return nil
}

type ArticleFilter struct {
	FeedID    string
	Author    string
	IsRead    *bool
	IsStarred *bool
	Keyword   string
}

func (f ArticleFilter) Match(a *Article) bool {
	if f.FeedID != "" && a.FeedID != f.FeedID {
		return false
	}
	if f.Author != "" && a.Author != f.Author {
		return false
	}
	if f.IsRead != nil && a.IsRead != *f.IsRead {
		return false
	}
	if f.IsStarred != nil && a.IsStarred != *f.IsStarred {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(a.Title), k) &&
			!strings.Contains(strings.ToLower(a.Summary), k) {
			return false
		}
	}
	return true
}
