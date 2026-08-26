// Package app 负责依赖装配。
package app

import (
	"net/http"

	"rss-reader/internal/config"
	"rss-reader/internal/handler"
	"rss-reader/internal/service"
	"rss-reader/internal/store"
	"rss-reader/pkg/logger"
)

type App struct {
	server *handler.Server
}

func New(cfg *config.Config, log *logger.Logger) (*App, error) {
	st := store.NewMemoryStore()
	svc := service.New(st, log, cfg)
	server := handler.NewServer(svc, log, cfg)
	log.Infof("应用装配完成，配置：%s", cfg.String())
	return &App{server: server}, nil
}

func (a *App) Routes() http.Handler { return a.server.Routes() }
