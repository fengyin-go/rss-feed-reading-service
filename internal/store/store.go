// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"rss-reader/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法，便于测试时替换实现。
type Store interface {
	// Feed
	CreateFeed(f *model.Feed) error
	GetFeed(id string) (*model.Feed, error)
	GetFeedByURL(url string) (*model.Feed, error)
	ListFeeds() []*model.Feed
	UpdateFeed(f *model.Feed) error
	DeleteFeed(id string) error

	// Article
	CreateArticle(a *model.Article) error
	GetArticle(id string) (*model.Article, error)
	GetArticleByGUID(guid string) (*model.Article, error)
	ListArticles() []*model.Article
	UpdateArticle(a *model.Article) error
	DeleteArticle(id string) error

	// Category
	CreateCategory(c *model.Category) error
	GetCategory(id string) (*model.Category, error)
	GetCategoryByName(name string) (*model.Category, error)
	ListCategories() []*model.Category
	UpdateCategory(c *model.Category) error
	DeleteCategory(id string) error

	// FetchTask
	CreateFetchTask(t *model.FetchTask) error
	GetFetchTask(id string) (*model.FetchTask, error)
	ListFetchTasks() []*model.FetchTask
	UpdateFetchTask(t *model.FetchTask) error
	DeleteFetchTask(id string) error

	// FetchLog
	CreateFetchLog(l *model.FetchLog) error
	GetFetchLog(id string) (*model.FetchLog, error)
	ListFetchLogs() []*model.FetchLog
	DeleteFetchLog(id string) error

	// User
	CreateUser(u *model.User) error
	GetUser(id string) (*model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	ListUsers() []*model.User
	UpdateUser(u *model.User) error
	DeleteUser(id string) error

	// Subscription
	CreateSubscription(s *model.Subscription) error
	GetSubscription(id string) (*model.Subscription, error)
	ListSubscriptions() []*model.Subscription
	DeleteSubscription(id string) error

	// Bookmark
	CreateBookmark(b *model.Bookmark) error
	GetBookmark(id string) (*model.Bookmark, error)
	ListBookmarks() []*model.Bookmark
	DeleteBookmark(id string) error

	// Notification
	CreateNotification(n *model.Notification) error
	GetNotification(id string) (*model.Notification, error)
	ListNotifications() []*model.Notification
	UpdateNotification(n *model.Notification) error
	DeleteNotification(id string) error

	// Tag
	CreateTag(t *model.Tag) error
	GetTag(id string) (*model.Tag, error)
	GetTagByName(name string) (*model.Tag, error)
	ListTags() []*model.Tag
	UpdateTag(t *model.Tag) error
	DeleteTag(id string) error

	// ArticleTag
	CreateArticleTag(at *model.ArticleTag) error
	GetArticleTag(id string) (*model.ArticleTag, error)
	ListArticleTags() []*model.ArticleTag
	DeleteArticleTag(id string) error
}
