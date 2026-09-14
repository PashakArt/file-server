package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/PashakArt/file-server/internal/config"
	"github.com/PashakArt/file-server/internal/service"
	"github.com/PashakArt/file-server/internal/transport/http/middleware"
	"github.com/PashakArt/file-server/internal/transport/http/types"
	"github.com/PashakArt/file-server/pkg/validator"
)

type DocHandler struct {
	docService *service.DocService
	cfg        *config.Config
}

func NewDocsHandler(
	docService *service.DocService,
	cfg *config.Config,
) *DocHandler {
	return &DocHandler{
		docService: docService,
		cfg:        cfg,
	}
}

func (h *DocHandler) RegisterRoutes(mux *http.ServeMux, authMiddleware middleware.Middleware) {
	mux.Handle("POST /api/docs", authMiddleware(http.HandlerFunc(h.upload)))
	// mux.HandleFunc("POST /api/docs", h.upload)
	mux.HandleFunc("GET /api/docs", h.get)
	mux.HandleFunc("GET /api/docs/{id}", h.getById)
	mux.HandleFunc("DELETE /api/docs/{id}", h.deleteById)

}

func (h *DocHandler) upload(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		types.SendError(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, int64(h.cfg.MaxUploadBodySizeByte))
	err := r.ParseMultipartForm(int64(h.cfg.MaxUploadBodySizeByte))
	if err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			types.SendError(w, r, http.StatusBadRequest, "Request size exceeds maximum allowed limit")
			return
		}

		types.SendError(w, r, http.StatusBadRequest, "Invalid multipart form data")
		return
	}

	metaStr := r.FormValue("meta")
	if metaStr == "" {
		types.SendError(w, r, http.StatusBadRequest, "Meta field is required")
		return
	}

	var metaField types.DocMeta
	err = json.Unmarshal([]byte(metaStr), &metaField)
	if err != nil {
		types.SendError(w, r, http.StatusBadRequest, "Invalid field meta")
		return
	}

	var file io.Reader
	if metaField.File {
		f, _, err := r.FormFile("file")
		if err != nil {
			types.SendError(w, r, http.StatusBadRequest, "invalid file field")
			return
		}
		defer f.Close()

		file = f
	}
	err = validator.ValidateStruct(&metaField)
	if err != nil {
		types.SendError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	var jsonField json.RawMessage
	jsonStr := r.FormValue("json")
	if jsonStr != "" {
		if !json.Valid([]byte(jsonStr)) {
			types.SendError(w, r, http.StatusBadRequest, "Invalid json field")
			return
		}
		json.Unmarshal([]byte(jsonStr), &jsonField)
	}

	err = h.docService.Upload(r.Context(), userID, file, &metaField, &jsonField)
	if err != nil {
		types.SendError(w, r, http.StatusInternalServerError, "Internal server error")
		log.Println(err)
		return
	}

	types.SendData(w, r, map[string]string{
		"file": metaField.Name,
		"json": jsonStr,
	})
}

func (h *DocHandler) get(w http.ResponseWriter, r *http.Request) {}

func (h *DocHandler) getById(w http.ResponseWriter, r *http.Request) {}

func (h *DocHandler) deleteById(w http.ResponseWriter, r *http.Request) {}
