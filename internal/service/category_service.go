package service

import (
	"sort"
	"time"

	"rss-reader/internal/model"
	"rss-reader/pkg/idgen"
)

func (s *Service) CreateCategory(input model.Category) (*model.Category, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	c := &model.Category{
		ID:          idgen.Hex(),
		Name:        input.Name,
		Description: input.Description,
		CreatedAt:   now,
	}
	if err := s.store.CreateCategory(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) GetCategory(id string) (*model.Category, error) {
	return s.store.GetCategory(id)
}

func (s *Service) ListCategories(filter model.CategoryFilter, page, size int) ([]*model.Category, int, error) {
	all := s.store.ListCategories()
	matched := make([]*model.Category, 0, len(all))
	for _, c := range all {
		if filter.Match(c) {
			matched = append(matched, c)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Category{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateCategory(id string, input model.Category) (*model.Category, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	c, err := s.store.GetCategory(id)
	if err != nil {
		return nil, err
	}
	c.Name = input.Name
	c.Description = input.Description
	if err := s.store.UpdateCategory(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) DeleteCategory(id string) error {
	return s.store.DeleteCategory(id)
}
