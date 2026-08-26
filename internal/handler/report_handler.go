package handler

import (
	"net/http"
	"strconv"

	"rss-reader/pkg/httpx"
)

func (s *Server) registerReportRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/reports/articles-by-date", s.articlesByDate)
	mux.HandleFunc("GET /api/reports/user-activity", s.userActivity)
	mux.HandleFunc("GET /api/reports/feed-health", s.feedHealth)
	mux.HandleFunc("GET /api/reports/unread-by-feed", s.unreadByFeed)
	mux.HandleFunc("GET /api/reports/task-duration", s.taskDurationStats)
}

func (s *Server) articlesByDate(w http.ResponseWriter, r *http.Request) {
	days := 30
	if v := r.URL.Query().Get("days"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			days = n
		}
	}
	result, err := s.svc.ArticlesByDate(days)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

func (s *Server) userActivity(w http.ResponseWriter, r *http.Request) {
	result, err := s.svc.UserActivity()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

func (s *Server) feedHealth(w http.ResponseWriter, r *http.Request) {
	result, err := s.svc.FeedHealth()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

func (s *Server) unreadByFeed(w http.ResponseWriter, r *http.Request) {
	result, err := s.svc.UnreadByFeed()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

func (s *Server) taskDurationStats(w http.ResponseWriter, r *http.Request) {
	result, err := s.svc.TaskDurationStats()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}
