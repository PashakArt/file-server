package http

import (
	"net/http"

	"github.com/PashakArt/file-server/internal/transport/http/handler"
)

func NewRouter(
	ah *handler.AuthHandler,
	dh *handler.DocsHandler,
) http.Handler {
	mux := http.NewServeMux()

	ah.RegisterRoutes(mux)
	dh.RegisterRoutes(mux)

	return mux
}
