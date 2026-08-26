package service

import (
	"sort"
	"time"

	"rss-reader/internal/model"
	"rss-reader/internal/store"
	"rss-reader/pkg/idgen"
)

func (s *Service) CreateFetchTask(input model.FetchTask) (*model.FetchTask, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetFeed(input.FeedID); err != nil {
		if err == store.ErrNotFound {
			return nil, model.NewValidationError("feed_id", "所属订阅源不存在")
		}
		return nil, err
	}
	now := time.Now()
	tk := &model.FetchTask{
		ID:        idgen.Hex(),
		FeedID:    input.FeedID,
		Status:    model.FetchTaskPending,
		CreatedAt: now,
	}
	if err := s.store.CreateFetchTask(tk); err != nil {
		return nil, err
	}
	return tk, nil
}

func (s *Service) GetFetchTask(id string) (*model.FetchTask, error) {
	return s.store.GetFetchTask(id)
}

func (s *Service) ListFetchTasks(filter model.FetchTaskFilter, page, size int) ([]*model.FetchTask, int, error) {
	all := s.store.ListFetchTasks()
	matched := make([]*model.FetchTask, 0, len(all))
	for _, t := range all {
		if filter.Match(t) {
			matched = append(matched, t)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.FetchTask{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) StartFetchTask(id string) (*model.FetchTask, error) {
	tk, err := s.store.GetFetchTask(id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransitionFetchTask(tk.Status, model.FetchTaskRunning) {
		return nil, store.ErrConflict
	}
	tk.Status = model.FetchTaskRunning
	tk.StartedAt = time.Now()
	if err := s.store.UpdateFetchTask(tk); err != nil {
		return nil, err
	}
	return tk, nil
}

func (s *Service) CompleteFetchTask(id string, success bool, fetchedCount int, errMsg string) (*model.FetchTask, error) {
	tk, err := s.store.GetFetchTask(id)
	if err != nil {
		return nil, err
	}
	targetStatus := model.FetchTaskSuccess
	if !success {
		targetStatus = model.FetchTaskFailed
	}
	if !model.CanTransitionFetchTask(tk.Status, targetStatus) {
		return nil, store.ErrConflict
	}
	tk.Status = targetStatus
	tk.FinishedAt = time.Now()
	tk.FetchedCount = fetchedCount
	tk.Error = errMsg
	if err := s.store.UpdateFetchTask(tk); err != nil {
		return nil, err
	}
	return tk, nil
}

func (s *Service) DeleteFetchTask(id string) error {
	return s.store.DeleteFetchTask(id)
}
