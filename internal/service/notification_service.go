package service

import (
	"sort"
	"time"

	"rss-reader/internal/model"
	"rss-reader/internal/store"
	"rss-reader/pkg/idgen"
)

func (s *Service) CreateNotification(input model.Notification) (*model.Notification, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetUser(input.UserID); err != nil {
		if err == store.ErrNotFound {
			return nil, model.NewValidationError("user_id", "用户不存在")
		}
		return nil, err
	}
	now := time.Now()
	n := &model.Notification{
		ID:        idgen.Hex(),
		UserID:    input.UserID,
		Type:      input.Type,
		Title:     input.Title,
		Content:   input.Content,
		CreatedAt: now,
	}
	if err := s.store.CreateNotification(n); err != nil {
		return nil, err
	}
	return n, nil
}

func (s *Service) GetNotification(id string) (*model.Notification, error) {
	return s.store.GetNotification(id)
}

func (s *Service) ListNotifications(filter model.NotificationFilter, page, size int) ([]*model.Notification, int, error) {
	all := s.store.ListNotifications()
	matched := make([]*model.Notification, 0, len(all))
	for _, n := range all {
		if filter.Match(n) {
			matched = append(matched, n)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Notification{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) MarkNotificationRead(id string) (*model.Notification, error) {
	n, err := s.store.GetNotification(id)
	if err != nil {
		return nil, err
	}
	n.IsRead = true
	if err := s.store.UpdateNotification(n); err != nil {
		return nil, err
	}
	return n, nil
}

func (s *Service) DeleteNotification(id string) error {
	return s.store.DeleteNotification(id)
}
