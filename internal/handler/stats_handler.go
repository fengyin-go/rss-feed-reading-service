package handler

import (
	"net/http"

	"rss-reader/pkg/httpx"
)

func (s *Server) registerStatsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/stats/global", s.globalStats)
	mux.HandleFunc("GET /api/stats/articles-by-category", s.articlesByCategory)
	mux.HandleFunc("GET /api/stats/articles-by-feed", s.articlesByFeed)
	mux.HandleFunc("GET /api/stats/fetch-rate-by-feed", s.fetchRateByFeed)
}

func (s *Server) globalStats(w http.ResponseWriter, r *http.Request) {
	stats, err := s.svc.GlobalStats()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, stats)
}

func (s *Server) articlesByCategory(w http.ResponseWriter, r *http.Request) {
	result, err := s.svc.ArticlesByCategory()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

func (s *Server) articlesByFeed(w http.ResponseWriter, r *http.Request) {
	result, err := s.svc.ArticlesByFeedTopN(10)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

func (s *Server) fetchRateByFeed(w http.ResponseWriter, r *http.Request) {
	result, err := s.svc.FetchRateByFeed()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}
