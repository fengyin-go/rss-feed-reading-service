package handler

import (
	"io"
	"net/http"

	"rss-reader/pkg/httpx"
)

func (s *Server) registerImportRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/import/feeds", s.importFeeds)
	mux.HandleFunc("POST /api/import/feeds/{id}/articles", s.importArticles)
	mux.HandleFunc("GET /api/export/feeds", s.exportFeeds)
	mux.HandleFunc("GET /api/export/all", s.exportAll)
}

func (s *Server) importFeeds(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(io.LimitReader(r.Body, 10<<20))
	if err != nil {
		httpx.BadRequest(w, "读取请求体失败: "+err.Error())
		return
	}
	result, err := s.svc.ImportFeeds(data)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

func (s *Server) importArticles(w http.ResponseWriter, r *http.Request) {
	feedID := r.PathValue("id")
	data, err := io.ReadAll(io.LimitReader(r.Body, 10<<20))
	if err != nil {
		httpx.BadRequest(w, "读取请求体失败: "+err.Error())
		return
	}
	result, err := s.svc.ImportArticles(feedID, data)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

func (s *Server) exportFeeds(w http.ResponseWriter, r *http.Request) {
	items, err := s.svc.ExportFeeds()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, items)
}

func (s *Server) exportAll(w http.ResponseWriter, r *http.Request) {
	data, err := s.svc.ExportAll()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, data)
}
