package handler

import (
	"net/http"

	"rss-reader/internal/model"
	"rss-reader/pkg/httpx"
)

func (s *Server) registerTagRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/tags", s.createTag)
	mux.HandleFunc("GET /api/tags", s.listTags)
	mux.HandleFunc("GET /api/tags/{id}", s.getTag)
	mux.HandleFunc("PUT /api/tags/{id}", s.updateTag)
	mux.HandleFunc("DELETE /api/tags/{id}", s.deleteTag)
	mux.HandleFunc("GET /api/tags/{id}/articles", s.articlesByTag)
}

type createTagRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

func (s *Server) createTag(w http.ResponseWriter, r *http.Request) {
	var req createTagRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.CreateTag(model.Tag{
		Name:  req.Name,
		Color: req.Color,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, t)
}

func (s *Server) listTags(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.TagFilter{
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListTags(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getTag(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	t, err := s.svc.GetTag(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

type updateTagRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

func (s *Server) updateTag(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateTagRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.UpdateTag(id, model.Tag{
		Name:  req.Name,
		Color: req.Color,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

func (s *Server) deleteTag(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteTag(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) articlesByTag(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	items, err := s.svc.ArticlesByTag(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, items)
}
