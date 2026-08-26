package model

import (
	"strings"
	"time"
)

const (
	NotificationTypeFeedUpdate = "feed_update"
	NotificationTypeSystem     = "system"
)

// Notification 表示用户通知。
type Notification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

func (n *Notification) Validate() error {
	n.Title = strings.TrimSpace(n.Title)
	n.Content = strings.TrimSpace(n.Content)
	if n.UserID == "" {
		return NewValidationError("user_id", "用户 ID 不能为空")
	}
	if n.Title == "" {
		return NewValidationError("title", "通知标题不能为空")
	}
	if n.Type == "" {
		n.Type = NotificationTypeSystem
	}
	if n.Type != NotificationTypeFeedUpdate && n.Type != NotificationTypeSystem {
		return NewValidationError("type", "通知类型不合法")
	}
	return nil
}

type NotificationFilter struct {
	UserID string
	Type   string
	IsRead *bool
}

func (f NotificationFilter) Match(n *Notification) bool {
	if f.UserID != "" && n.UserID != f.UserID {
		return false
	}
	if f.Type != "" && n.Type != f.Type {
		return false
	}
	if f.IsRead != nil && n.IsRead != *f.IsRead {
		return false
	}
	return true
}
