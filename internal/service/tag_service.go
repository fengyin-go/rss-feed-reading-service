package service

import (
	"sort"
	"time"

	"rss-reader/internal/model"
	"rss-reader/pkg/idgen"
)

func (s *Service) CreateTag(input model.Tag) (*model.Tag, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	t := &model.Tag{
		ID:        idgen.Hex(),
		Name:      input.Name,
		Color:     input.Color,
		CreatedAt: now,
	}
	if err := s.store.CreateTag(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) GetTag(id string) (*model.Tag, error) {
	return s.store.GetTag(id)
}

func (s *Service) ListTags(filter model.TagFilter, page, size int) ([]*model.Tag, int, error) {
	all := s.store.ListTags()
	matched := make([]*model.Tag, 0, len(all))
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
		return []*model.Tag{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateTag(id string, input model.Tag) (*model.Tag, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	t, err := s.store.GetTag(id)
	if err != nil {
		return nil, err
	}
	t.Name = input.Name
	t.Color = input.Color
	if err := s.store.UpdateTag(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) DeleteTag(id string) error {
	return s.store.DeleteTag(id)
}
