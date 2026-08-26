package service

import (
	"sort"
	"time"

	"rss-reader/internal/model"
	"rss-reader/pkg/idgen"
)

func (s *Service) CreateUser(input model.User) (*model.User, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	u := &model.User{
		ID:        idgen.Hex(),
		Username:  input.Username,
		Email:     input.Email,
		CreatedAt: now,
	}
	if err := s.store.CreateUser(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Service) GetUser(id string) (*model.User, error) {
	return s.store.GetUser(id)
}

func (s *Service) ListUsers(filter model.UserFilter, page, size int) ([]*model.User, int, error) {
	all := s.store.ListUsers()
	matched := make([]*model.User, 0, len(all))
	for _, u := range all {
		if filter.Match(u) {
			matched = append(matched, u)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.User{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateUser(id string, input model.User) (*model.User, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	u, err := s.store.GetUser(id)
	if err != nil {
		return nil, err
	}
	u.Username = input.Username
	u.Email = input.Email
	if err := s.store.UpdateUser(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Service) DeleteUser(id string) error {
	return s.store.DeleteUser(id)
}
