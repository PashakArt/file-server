package handler

import (
	"net/http"
)

type DocsHandler struct{}

func NewDocsHandler() *DocsHandler {
	return &DocsHandler{}
}

func (h *DocsHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/docs", h.upload)
	mux.HandleFunc("GET /api/docs", h.get)
	mux.HandleFunc("GET /api/docs/{id}", h.getById)
	mux.HandleFunc("DELETE /api/docs/{id}", h.deleteById)

}

func (h *DocsHandler) upload(w http.ResponseWriter, r *http.Request) {}

func (h *DocsHandler) get(w http.ResponseWriter, r *http.Request) {}

func (h *DocsHandler) getById(w http.ResponseWriter, r *http.Request) {}

func (h *DocsHandler) deleteById(w http.ResponseWriter, r *http.Request) {}
