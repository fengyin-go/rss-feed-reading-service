package store

import (
	"sync"

	"rss-reader/internal/model"
)

type MemoryStore struct {
	mu            sync.RWMutex
	feeds         map[string]*model.Feed
	articles      map[string]*model.Article
	categories    map[string]*model.Category
	fetchTasks    map[string]*model.FetchTask
	fetchLogs     map[string]*model.FetchLog
	users         map[string]*model.User
	subscriptions map[string]*model.Subscription
	bookmarks     map[string]*model.Bookmark
	notifications map[string]*model.Notification
	tags          map[string]*model.Tag
	articleTags   map[string]*model.ArticleTag
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		feeds:         make(map[string]*model.Feed),
		articles:      make(map[string]*model.Article),
		categories:    make(map[string]*model.Category),
		fetchTasks:    make(map[string]*model.FetchTask),
		fetchLogs:     make(map[string]*model.FetchLog),
		users:         make(map[string]*model.User),
		subscriptions: make(map[string]*model.Subscription),
		bookmarks:     make(map[string]*model.Bookmark),
		notifications: make(map[string]*model.Notification),
		tags:          make(map[string]*model.Tag),
		articleTags:   make(map[string]*model.ArticleTag),
	}
}

var _ Store = (*MemoryStore)(nil)
