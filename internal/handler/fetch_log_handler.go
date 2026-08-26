package handler

import (
	"net/http"

	"rss-reader/internal/model"
	"rss-reader/pkg/httpx"
)

func (s *Server) registerFetchLogRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/fetch-logs", s.createFetchLog)
	mux.HandleFunc("GET /api/fetch-logs", s.listFetchLogs)
	mux.HandleFunc("GET /api/fetch-logs/{id}", s.getFetchLog)
	mux.HandleFunc("DELETE /api/fetch-logs/{id}", s.deleteFetchLog)
}

type createFetchLogRequest struct {
	TaskID  string `json:"task_id"`
	FeedID  string `json:"feed_id"`
	Level   string `json:"level"`
	Message string `json:"message"`
}

func (s *Server) createFetchLog(w http.ResponseWriter, r *http.Request) {
	var req createFetchLogRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	l, err := s.svc.CreateFetchLog(model.FetchLog{
		TaskID:  req.TaskID,
		FeedID:  req.FeedID,
		Level:   req.Level,
		Message: req.Message,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, l)
}

func (s *Server) listFetchLogs(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.FetchLogFilter{
		FeedID: r.URL.Query().Get("feed_id"),
		TaskID: r.URL.Query().Get("task_id"),
		Level:  r.URL.Query().Get("level"),
	}
	items, total, err := s.svc.ListFetchLogs(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getFetchLog(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	l, err := s.svc.GetFetchLog(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, l)
}

func (s *Server) deleteFetchLog(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteFetchLog(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
