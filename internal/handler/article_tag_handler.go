package handler

import (
	"net/http"

	"rss-reader/internal/model"
	"rss-reader/pkg/httpx"
)

func (s *Server) registerArticleTagRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/article-tags", s.createArticleTag)
	mux.HandleFunc("GET /api/article-tags", s.listArticleTags)
	mux.HandleFunc("GET /api/article-tags/{id}", s.getArticleTag)
	mux.HandleFunc("DELETE /api/article-tags/{id}", s.deleteArticleTag)
	mux.HandleFunc("GET /api/articles/{id}/tags", s.tagsByArticle)
}

type createArticleTagRequest struct {
	ArticleID string `json:"article_id"`
	TagID     string `json:"tag_id"`
}

func (s *Server) createArticleTag(w http.ResponseWriter, r *http.Request) {
	var req createArticleTagRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	at, err := s.svc.CreateArticleTag(model.ArticleTag{
		ArticleID: req.ArticleID,
		TagID:     req.TagID,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, at)
}

func (s *Server) listArticleTags(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ArticleTagFilter{
		ArticleID: r.URL.Query().Get("article_id"),
		TagID:     r.URL.Query().Get("tag_id"),
	}
	items, total, err := s.svc.ListArticleTags(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getArticleTag(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	at, err := s.svc.GetArticleTag(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, at)
}

func (s *Server) deleteArticleTag(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteArticleTag(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) tagsByArticle(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	items, err := s.svc.TagsByArticle(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, items)
}
