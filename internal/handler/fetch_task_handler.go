package handler

import (
	"net/http"

	"rss-reader/internal/model"
	"rss-reader/pkg/httpx"
)

func (s *Server) registerFetchTaskRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/fetch-tasks", s.createFetchTask)
	mux.HandleFunc("GET /api/fetch-tasks", s.listFetchTasks)
	mux.HandleFunc("GET /api/fetch-tasks/{id}", s.getFetchTask)
	mux.HandleFunc("POST /api/fetch-tasks/{id}/start", s.startFetchTask)
	mux.HandleFunc("POST /api/fetch-tasks/{id}/complete", s.completeFetchTask)
	mux.HandleFunc("DELETE /api/fetch-tasks/{id}", s.deleteFetchTask)
}

type createFetchTaskRequest struct {
	FeedID string `json:"feed_id"`
}

func (s *Server) createFetchTask(w http.ResponseWriter, r *http.Request) {
	var req createFetchTaskRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	tk, err := s.svc.CreateFetchTask(model.FetchTask{FeedID: req.FeedID})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, tk)
}

func (s *Server) listFetchTasks(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.FetchTaskFilter{
		FeedID: r.URL.Query().Get("feed_id"),
		Status: r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListFetchTasks(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getFetchTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	tk, err := s.svc.GetFetchTask(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, tk)
}

func (s *Server) startFetchTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	tk, err := s.svc.StartFetchTask(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, tk)
}

type completeFetchTaskRequest struct {
	Success      bool   `json:"success"`
	FetchedCount int    `json:"fetched_count"`
	Error        string `json:"error"`
}

func (s *Server) completeFetchTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req completeFetchTaskRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	tk, err := s.svc.CompleteFetchTask(id, req.Success, req.FetchedCount, req.Error)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, tk)
}

func (s *Server) deleteFetchTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteFetchTask(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
