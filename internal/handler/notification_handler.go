package handler

import (
	"net/http"
	"strconv"

	"rss-reader/internal/model"
	"rss-reader/pkg/httpx"
)

func (s *Server) registerNotificationRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/notifications", s.createNotification)
	mux.HandleFunc("GET /api/notifications", s.listNotifications)
	mux.HandleFunc("GET /api/notifications/{id}", s.getNotification)
	mux.HandleFunc("POST /api/notifications/{id}/read", s.markNotificationRead)
	mux.HandleFunc("DELETE /api/notifications/{id}", s.deleteNotification)
}

type createNotificationRequest struct {
	UserID  string `json:"user_id"`
	Type    string `json:"type"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (s *Server) createNotification(w http.ResponseWriter, r *http.Request) {
	var req createNotificationRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	n, err := s.svc.CreateNotification(model.Notification{
		UserID:  req.UserID,
		Type:    req.Type,
		Title:   req.Title,
		Content: req.Content,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, n)
}

func (s *Server) listNotifications(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.NotificationFilter{
		UserID: r.URL.Query().Get("user_id"),
		Type:   r.URL.Query().Get("type"),
	}
	if v := r.URL.Query().Get("is_read"); v != "" {
		b, _ := strconv.ParseBool(v)
		filter.IsRead = &b
	}
	items, total, err := s.svc.ListNotifications(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getNotification(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	n, err := s.svc.GetNotification(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, n)
}

func (s *Server) markNotificationRead(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	n, err := s.svc.MarkNotificationRead(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, n)
}

func (s *Server) deleteNotification(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteNotification(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
