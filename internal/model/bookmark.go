package model

import (
	"time"
)

// Bookmark 表示用户对文章的收藏。
type Bookmark struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	ArticleID string    `json:"article_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (b *Bookmark) Validate() error {
	if b.UserID == "" {
		return NewValidationError("user_id", "用户 ID 不能为空")
	}
	if b.ArticleID == "" {
		return NewValidationError("article_id", "文章 ID 不能为空")
	}
	return nil
}

type BookmarkFilter struct {
	UserID    string
	ArticleID string
}

func (f BookmarkFilter) Match(b *Bookmark) bool {
	if f.UserID != "" && b.UserID != f.UserID {
		return false
	}
	if f.ArticleID != "" && b.ArticleID != f.ArticleID {
		return false
	}
	return true
}
