package handler

import (
	"net/http"

	"rss-reader/internal/model"
	"rss-reader/pkg/httpx"
)

func (s *Server) registerFeedRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/feeds", s.createFeed)
	mux.HandleFunc("GET /api/feeds", s.listFeeds)
	mux.HandleFunc("GET /api/feeds/{id}", s.getFeed)
	mux.HandleFunc("PUT /api/feeds/{id}", s.updateFeed)
	mux.HandleFunc("DELETE /api/feeds/{id}", s.deleteFeed)
	mux.HandleFunc("POST /api/feeds/{id}/pause", s.pauseFeed)
	mux.HandleFunc("POST /api/feeds/{id}/resume", s.resumeFeed)
}

type createFeedRequest struct {
	Title         string `json:"title"`
	URL           string `json:"url"`
	Description   string `json:"description"`
	Category      string `json:"category"`
	FetchInterval int    `json:"fetch_interval"`
}

func (s *Server) createFeed(w http.ResponseWriter, r *http.Request) {
	var req createFeedRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	f, err := s.svc.CreateFeed(model.Feed{
		Title:         req.Title,
		URL:           req.URL,
		Description:   req.Description,
		Category:      req.Category,
		FetchInterval: req.FetchInterval,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, f)
}

func (s *Server) listFeeds(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.FeedFilter{
		Category: r.URL.Query().Get("category"),
		Status:   r.URL.Query().Get("status"),
		Keyword:  r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListFeeds(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getFeed(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	f, err := s.svc.GetFeed(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, f)
}

type updateFeedRequest struct {
	Title         string `json:"title"`
	URL           string `json:"url"`
	Description   string `json:"description"`
	Category      string `json:"category"`
	FetchInterval int    `json:"fetch_interval"`
}

func (s *Server) updateFeed(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateFeedRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	f, err := s.svc.UpdateFeed(id, model.Feed{
		Title:         req.Title,
		URL:           req.URL,
		Description:   req.Description,
		Category:      req.Category,
		FetchInterval: req.FetchInterval,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, f)
}

func (s *Server) deleteFeed(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteFeed(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) pauseFeed(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	f, err := s.svc.PauseFeed(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, f)
}

func (s *Server) resumeFeed(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	f, err := s.svc.ResumeFeed(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, f)
}
