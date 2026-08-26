package model

import (
	"strings"
	"time"
)

// FetchLog 表示抓取日志。
type FetchLog struct {
	ID        string    `json:"id"`
	TaskID    string    `json:"task_id"`
	FeedID    string    `json:"feed_id"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

func (l *FetchLog) Validate() error {
	l.Message = strings.TrimSpace(l.Message)
	l.Level = strings.TrimSpace(l.Level)
	if l.FeedID == "" {
		return NewValidationError("feed_id", "订阅源 ID 不能为空")
	}
	if l.TaskID == "" {
		return NewValidationError("task_id", "任务 ID 不能为空")
	}
	if l.Message == "" {
		return NewValidationError("message", "日志消息不能为空")
	}
	if l.Level == "" {
		l.Level = "info"
	}
	return nil
}

type FetchLogFilter struct {
	FeedID string
	TaskID string
	Level  string
}

func (f FetchLogFilter) Match(l *FetchLog) bool {
	if f.FeedID != "" && l.FeedID != f.FeedID {
		return false
	}
	if f.TaskID != "" && l.TaskID != f.TaskID {
		return false
	}
	if f.Level != "" && l.Level != f.Level {
		return false
	}
	return true
}
