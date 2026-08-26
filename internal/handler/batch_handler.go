package handler

import (
	"net/http"

	"rss-reader/internal/model"
	"rss-reader/pkg/httpx"
)

func (s *Server) registerBatchRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/batch/articles", s.batchCreateArticles)
	mux.HandleFunc("POST /api/batch/articles/read", s.batchMarkArticlesRead)
	mux.HandleFunc("POST /api/batch/articles/star", s.batchStarArticles)
	mux.HandleFunc("POST /api/batch/articles/unstar", s.batchUnstarArticles)
	mux.HandleFunc("POST /api/batch/articles/delete", s.batchDeleteArticles)
	mux.HandleFunc("POST /api/batch/feeds", s.batchCreateFeeds)
}

type batchCreateArticlesRequest struct {
	FeedID   string          `json:"feed_id"`
	Articles []model.Article `json:"articles"`
}

func (s *Server) batchCreateArticles(w http.ResponseWriter, r *http.Request) {
	var req batchCreateArticlesRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	result, err := s.svc.BatchCreateArticles(req.FeedID, req.Articles)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

type batchMarkArticlesReadRequest struct {
	ArticleIDs []string `json:"article_ids"`
}

func (s *Server) batchMarkArticlesRead(w http.ResponseWriter, r *http.Request) {
	var req batchMarkArticlesReadRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	result, err := s.svc.BatchMarkArticlesRead(req.ArticleIDs)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

type batchStarArticlesRequest struct {
	ArticleIDs []string `json:"article_ids"`
}

func (s *Server) batchStarArticles(w http.ResponseWriter, r *http.Request) {
	var req batchStarArticlesRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	result, err := s.svc.BatchStarArticles(req.ArticleIDs, true)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

func (s *Server) batchUnstarArticles(w http.ResponseWriter, r *http.Request) {
	var req batchStarArticlesRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	result, err := s.svc.BatchStarArticles(req.ArticleIDs, false)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

type batchDeleteArticlesRequest struct {
	ArticleIDs []string `json:"article_ids"`
}

func (s *Server) batchDeleteArticles(w http.ResponseWriter, r *http.Request) {
	var req batchDeleteArticlesRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	result, err := s.svc.BatchDeleteArticles(req.ArticleIDs)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

type batchCreateFeedsRequest struct {
	Feeds []model.Feed `json:"feeds"`
}

func (s *Server) batchCreateFeeds(w http.ResponseWriter, r *http.Request) {
	var req batchCreateFeedsRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	result, err := s.svc.BatchCreateFeeds(req.Feeds)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}
