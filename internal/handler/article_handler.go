package handler

import (
	"net/http"
	"strconv"

	"rss-reader/internal/model"
	"rss-reader/pkg/httpx"
)

func (s *Server) registerArticleRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/articles", s.createArticle)
	mux.HandleFunc("GET /api/articles", s.listArticles)
	mux.HandleFunc("GET /api/articles/{id}", s.getArticle)
	mux.HandleFunc("PUT /api/articles/{id}", s.updateArticle)
	mux.HandleFunc("DELETE /api/articles/{id}", s.deleteArticle)
	mux.HandleFunc("POST /api/articles/{id}/read", s.markArticleRead)
	mux.HandleFunc("POST /api/articles/{id}/star", s.starArticle)
	mux.HandleFunc("POST /api/articles/{id}/unstar", s.unstarArticle)
	mux.HandleFunc("GET /api/feeds/{id}/export", s.exportFeedArticles)
}

type createArticleRequest struct {
	FeedID      string `json:"feed_id"`
	GUID        string `json:"guid"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	Summary     string `json:"summary"`
	Content     string `json:"content"`
	Author      string `json:"author"`
	PublishedAt string `json:"published_at"`
}

func (s *Server) createArticle(w http.ResponseWriter, r *http.Request) {
	var req createArticleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	a, err := s.svc.CreateArticle(model.Article{
		FeedID:  req.FeedID,
		GUID:    req.GUID,
		Title:   req.Title,
		URL:     req.URL,
		Summary: req.Summary,
		Content: req.Content,
		Author:  req.Author,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, a)
}

func (s *Server) listArticles(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ArticleFilter{
		FeedID:  r.URL.Query().Get("feed_id"),
		Author:  r.URL.Query().Get("author"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	if v := r.URL.Query().Get("is_read"); v != "" {
		b, _ := strconv.ParseBool(v)
		filter.IsRead = &b
	}
	if v := r.URL.Query().Get("is_starred"); v != "" {
		b, _ := strconv.ParseBool(v)
		filter.IsStarred = &b
	}
	items, total, err := s.svc.ListArticles(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getArticle(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	a, err := s.svc.GetArticle(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, a)
}

type updateArticleRequest struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Summary     string `json:"summary"`
	Content     string `json:"content"`
	Author      string `json:"author"`
	PublishedAt string `json:"published_at"`
}

func (s *Server) updateArticle(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateArticleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	a, err := s.svc.UpdateArticle(id, model.Article{
		Title:   req.Title,
		URL:     req.URL,
		Summary: req.Summary,
		Content: req.Content,
		Author:  req.Author,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, a)
}

func (s *Server) deleteArticle(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteArticle(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) markArticleRead(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	a, err := s.svc.MarkArticleRead(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, a)
}

func (s *Server) starArticle(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	a, err := s.svc.MarkArticleStarred(id, true)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, a)
}

func (s *Server) unstarArticle(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	a, err := s.svc.MarkArticleStarred(id, false)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, a)
}

func (s *Server) exportFeedArticles(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	items, err := s.svc.ExportArticlesByFeed(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{"feed_id": id, "articles": items})
}
