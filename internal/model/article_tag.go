package model

import (
	"time"
)

// ArticleTag 表示文章与标签的关联。
type ArticleTag struct {
	ID        string    `json:"id"`
	ArticleID string    `json:"article_id"`
	TagID     string    `json:"tag_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (at *ArticleTag) Validate() error {
	if at.ArticleID == "" {
		return NewValidationError("article_id", "文章 ID 不能为空")
	}
	if at.TagID == "" {
		return NewValidationError("tag_id", "标签 ID 不能为空")
	}
	return nil
}

type ArticleTagFilter struct {
	ArticleID string
	TagID     string
}

func (f ArticleTagFilter) Match(at *ArticleTag) bool {
	if f.ArticleID != "" && at.ArticleID != f.ArticleID {
		return false
	}
	if f.TagID != "" && at.TagID != f.TagID {
		return false
	}
	return true
}
