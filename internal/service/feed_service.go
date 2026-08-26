package service

import (
	"sort"
	"time"

	"rss-reader/internal/model"
	"rss-reader/pkg/idgen"
)

func (s *Service) CreateFeed(input model.Feed) (*model.Feed, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	f := &model.Feed{
		ID:            idgen.Hex(),
		Title:         input.Title,
		URL:           input.URL,
		Description:   input.Description,
		Category:      input.Category,
		Status:        input.Status,
		FetchInterval: input.FetchInterval,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := s.store.CreateFeed(f); err != nil {
		return nil, err
	}
	s.log.Infof("创建订阅源: %s", f.Title)
	return f, nil
}

func (s *Service) GetFeed(id string) (*model.Feed, error) {
	return s.store.GetFeed(id)
}

func (s *Service) ListFeeds(filter model.FeedFilter, page, size int) ([]*model.Feed, int, error) {
	all := s.store.ListFeeds()
	matched := make([]*model.Feed, 0, len(all))
	for _, f := range all {
		if filter.Match(f) {
			matched = append(matched, f)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Feed{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateFeed(id string, input model.Feed) (*model.Feed, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	f, err := s.store.GetFeed(id)
	if err != nil {
		return nil, err
	}
	f.Title = input.Title
	f.URL = input.URL
	f.Description = input.Description
	f.Category = input.Category
	f.FetchInterval = input.FetchInterval
	f.UpdatedAt = time.Now()
	if err := s.store.UpdateFeed(f); err != nil {
		return nil, err
	}
	return f, nil
}

func (s *Service) DeleteFeed(id string) error {
	return s.store.DeleteFeed(id)
}

func (s *Service) PauseFeed(id string) (*model.Feed, error) {
	f, err := s.store.GetFeed(id)
	if err != nil {
		return nil, err
	}
	if f.Status == model.FeedStatusPaused {
		return nil, model.NewValidationError("status", "订阅源已经是暂停状态")
	}
	if !model.CanTransitionFeed(f.Status, model.FeedStatusPaused) {
		return nil, model.NewValidationError("status", "状态流转不合法")
	}
	f.Status = model.FeedStatusPaused
	f.UpdatedAt = time.Now()
	if err := s.store.UpdateFeed(f); err != nil {
		return nil, err
	}
	s.log.Infof("暂停订阅源: %s", f.Title)
	return f, nil
}

func (s *Service) ResumeFeed(id string) (*model.Feed, error) {
	f, err := s.store.GetFeed(id)
	if err != nil {
		return nil, err
	}
	if f.Status == model.FeedStatusActive {
		return nil, model.NewValidationError("status", "订阅源已经是活跃状态")
	}
	if !model.CanTransitionFeed(f.Status, model.FeedStatusActive) {
		return nil, model.NewValidationError("status", "状态流转不合法")
	}
	f.Status = model.FeedStatusActive
	f.UpdatedAt = time.Now()
	if err := s.store.UpdateFeed(f); err != nil {
		return nil, err
	}
	s.log.Infof("恢复订阅源: %s", f.Title)
	return f, nil
}
