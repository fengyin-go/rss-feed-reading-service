package model

import (
	"time"
)

const (
	FetchTaskPending = "pending"
	FetchTaskRunning = "running"
	FetchTaskSuccess = "success"
	FetchTaskFailed  = "failed"
)

// FetchTask 表示一次抓取任务。
type FetchTask struct {
	ID           string    `json:"id"`
	FeedID       string    `json:"feed_id"`
	Status       string    `json:"status"`
	StartedAt    time.Time `json:"started_at"`
	FinishedAt   time.Time `json:"finished_at"`
	FetchedCount int       `json:"fetched_count"`
	Error        string    `json:"error"`
	CreatedAt    time.Time `json:"created_at"`
}

func (t *FetchTask) Validate() error {
	if t.FeedID == "" {
		return NewValidationError("feed_id", "订阅源 ID 不能为空")
	}
	if t.Status == "" {
		t.Status = FetchTaskPending
	}
	if t.Status != FetchTaskPending && t.Status != FetchTaskRunning &&
		t.Status != FetchTaskSuccess && t.Status != FetchTaskFailed {
		return NewValidationError("status", "任务状态不合法")
	}
	return nil
}

type FetchTaskFilter struct {
	FeedID string
	Status string
}

func (f FetchTaskFilter) Match(t *FetchTask) bool {
	if f.FeedID != "" && t.FeedID != f.FeedID {
		return false
	}
	if f.Status != "" && t.Status != f.Status {
		return false
	}
	return true
}

var fetchTaskTransitions = map[string]map[string]bool{
	FetchTaskPending: {FetchTaskRunning: true},
	FetchTaskRunning: {FetchTaskSuccess: true, FetchTaskFailed: true},
}

func CanTransitionFetchTask(from, to string) bool {
	if m, ok := fetchTaskTransitions[from]; ok {
		return m[to]
	}
	return false
}
