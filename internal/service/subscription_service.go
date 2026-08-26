package service

import (
	"sort"
	"time"

	"rss-reader/internal/model"
	"rss-reader/internal/store"
	"rss-reader/pkg/idgen"
)

func (s *Service) CreateSubscription(input model.Subscription) (*model.Subscription, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetUser(input.UserID); err != nil {
		if err == store.ErrNotFound {
			return nil, model.NewValidationError("user_id", "用户不存在")
		}
		return nil, err
	}
	if _, err := s.store.GetFeed(input.FeedID); err != nil {
		if err == store.ErrNotFound {
			return nil, model.NewValidationError("feed_id", "订阅源不存在")
		}
		return nil, err
	}
	now := time.Now()
	sub := &model.Subscription{
		ID:        idgen.Hex(),
		UserID:    input.UserID,
		FeedID:    input.FeedID,
		CreatedAt: now,
	}
	if err := s.store.CreateSubscription(sub); err != nil {
		return nil, err
	}
	return sub, nil
}

func (s *Service) GetSubscription(id string) (*model.Subscription, error) {
	return s.store.GetSubscription(id)
}

func (s *Service) ListSubscriptions(filter model.SubscriptionFilter, page, size int) ([]*model.Subscription, int, error) {
	all := s.store.ListSubscriptions()
	matched := make([]*model.Subscription, 0, len(all))
	for _, sub := range all {
		if filter.Match(sub) {
			matched = append(matched, sub)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Subscription{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) DeleteSubscription(id string) error {
	return s.store.DeleteSubscription(id)
}
