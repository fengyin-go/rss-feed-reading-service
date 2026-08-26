package service

import (
	"sort"
	"time"

	"rss-reader/internal/model"
	"rss-reader/internal/store"
	"rss-reader/pkg/idgen"
)

func (s *Service) CreateFetchLog(input model.FetchLog) (*model.FetchLog, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetFeed(input.FeedID); err != nil {
		if err == store.ErrNotFound {
			return nil, model.NewValidationError("feed_id", "所属订阅源不存在")
		}
		return nil, err
	}
	if _, err := s.store.GetFetchTask(input.TaskID); err != nil {
		if err == store.ErrNotFound {
			return nil, model.NewValidationError("task_id", "所属任务不存在")
		}
		return nil, err
	}
	now := time.Now()
	l := &model.FetchLog{
		ID:        idgen.Hex(),
		TaskID:    input.TaskID,
		FeedID:    input.FeedID,
		Level:     input.Level,
		Message:   input.Message,
		CreatedAt: now,
	}
	if err := s.store.CreateFetchLog(l); err != nil {
		return nil, err
	}
	return l, nil
}

func (s *Service) GetFetchLog(id string) (*model.FetchLog, error) {
	return s.store.GetFetchLog(id)
}

func (s *Service) ListFetchLogs(filter model.FetchLogFilter, page, size int) ([]*model.FetchLog, int, error) {
	all := s.store.ListFetchLogs()
	matched := make([]*model.FetchLog, 0, len(all))
	for _, l := range all {
		if filter.Match(l) {
			matched = append(matched, l)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.FetchLog{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) DeleteFetchLog(id string) error {
	return s.store.DeleteFetchLog(id)
}
