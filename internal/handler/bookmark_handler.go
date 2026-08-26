package handler

import (
	"net/http"

	"rss-reader/internal/model"
	"rss-reader/pkg/httpx"
)

func (s *Server) registerBookmarkRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/bookmarks", s.createBookmark)
	mux.HandleFunc("GET /api/bookmarks", s.listBookmarks)
	mux.HandleFunc("GET /api/bookmarks/{id}", s.getBookmark)
	mux.HandleFunc("DELETE /api/bookmarks/{id}", s.deleteBookmark)
}

type createBookmarkRequest struct {
	UserID    string `json:"user_id"`
	ArticleID string `json:"article_id"`
}

func (s *Server) createBookmark(w http.ResponseWriter, r *http.Request) {
	var req createBookmarkRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	b, err := s.svc.CreateBookmark(model.Bookmark{
		UserID:    req.UserID,
		ArticleID: req.ArticleID,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, b)
}

func (s *Server) listBookmarks(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.BookmarkFilter{
		UserID:    r.URL.Query().Get("user_id"),
		ArticleID: r.URL.Query().Get("article_id"),
	}
	items, total, err := s.svc.ListBookmarks(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getBookmark(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	b, err := s.svc.GetBookmark(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, b)
}

func (s *Server) deleteBookmark(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteBookmark(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
