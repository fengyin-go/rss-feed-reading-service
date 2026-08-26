package model

import (
	"time"
)

// Subscription 表示用户对订阅源的订阅关系。
type Subscription struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	FeedID    string    `json:"feed_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *Subscription) Validate() error {
	if s.UserID == "" {
		return NewValidationError("user_id", "用户 ID 不能为空")
	}
	if s.FeedID == "" {
		return NewValidationError("feed_id", "订阅源 ID 不能为空")
	}
	return nil
}

type SubscriptionFilter struct {
	UserID string
	FeedID string
}

func (f SubscriptionFilter) Match(s *Subscription) bool {
	if f.UserID != "" && s.UserID != f.UserID {
		return false
	}
	if f.FeedID != "" && s.FeedID != f.FeedID {
		return false
	}
	return true
}
