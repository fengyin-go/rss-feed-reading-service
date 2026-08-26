package service

import (
	"sort"
	"strings"

	"rss-reader/internal/model"
)

// SearchResult 搜索结果。
type SearchResult struct {
	Articles []*model.Article `json:"articles"`
	Feeds    []*model.Feed    `json:"feeds"`
}

func (s *Service) Search(keyword string, page, size int) (*SearchResult, error) {
	keyword = strings.ToLower(strings.TrimSpace(keyword))
	result := &SearchResult{}

	if keyword != "" {
		allArticles := s.store.ListArticles()
		for _, a := range allArticles {
			if strings.Contains(strings.ToLower(a.Title), keyword) ||
				strings.Contains(strings.ToLower(a.Summary), keyword) ||
				strings.Contains(strings.ToLower(a.Content), keyword) ||
				strings.Contains(strings.ToLower(a.Author), keyword) {
				result.Articles = append(result.Articles, a)
			}
		}
		sort.Slice(result.Articles, func(i, j int) bool {
			return result.Articles[i].PublishedAt.After(result.Articles[j].PublishedAt)
		})

		allFeeds := s.store.ListFeeds()
		for _, f := range allFeeds {
			if strings.Contains(strings.ToLower(f.Title), keyword) ||
				strings.Contains(strings.ToLower(f.Description), keyword) ||
				strings.Contains(strings.ToLower(f.URL), keyword) ||
				strings.Contains(strings.ToLower(f.Category), keyword) {
				result.Feeds = append(result.Feeds, f)
			}
		}
		sort.Slice(result.Feeds, func(i, j int) bool {
			return result.Feeds[i].CreatedAt.After(result.Feeds[j].CreatedAt)
		})
	}

	start := (page - 1) * size
	if start < len(result.Articles) {
		end := start + size
		if end > len(result.Articles) {
			end = len(result.Articles)
		}
		result.Articles = result.Articles[start:end]
	} else {
		result.Articles = []*model.Article{}
	}

	feedStart := (page - 1) * size
	if feedStart < len(result.Feeds) {
		feedEnd := feedStart + size
		if feedEnd > len(result.Feeds) {
			feedEnd = len(result.Feeds)
		}
		result.Feeds = result.Feeds[feedStart:feedEnd]
	} else {
		result.Feeds = []*model.Feed{}
	}

	return result, nil
}

// RecentArticles 获取最近文章。
func (s *Service) RecentArticles(limit int) ([]*model.Article, error) {
	all := s.store.ListArticles()
	sort.Slice(all, func(i, j int) bool {
		return all[i].PublishedAt.After(all[j].PublishedAt)
	})
	if limit > 0 && limit < len(all) {
		return all[:limit], nil
	}
	return all, nil
}

// UnreadArticles 获取未读文章。
func (s *Service) UnreadArticles(page, size int) ([]*model.Article, int, error) {
	all := s.store.ListArticles()
	var matched []*model.Article
	for _, a := range all {
		if !a.IsRead {
			matched = append(matched, a)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].PublishedAt.After(matched[j].PublishedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Article{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// StarredArticles 获取收藏文章。
func (s *Service) StarredArticles(page, size int) ([]*model.Article, int, error) {
	all := s.store.ListArticles()
	var matched []*model.Article
	for _, a := range all {
		if a.IsStarred {
			matched = append(matched, a)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].PublishedAt.After(matched[j].PublishedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Article{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}
