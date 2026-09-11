package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/PashakArt/file-server/internal/domain"
	"github.com/PashakArt/file-server/internal/service"
	"github.com/PashakArt/file-server/internal/transport/http/types"
	"github.com/PashakArt/file-server/internal/validator"
)

type AuthHandler struct {
	authService *service.AuthService
	adminToken  string
}

func NewAuthHandler(
	authService *service.AuthService,
	adminToken string,
) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		adminToken:  adminToken,
	}
}

func (h *AuthHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/register", h.register)
	mux.HandleFunc("POST /api/auth", h.login)
	mux.HandleFunc("DELETE /api/auth/{token}", h.logout)

}

func (h *AuthHandler) register(w http.ResponseWriter, r *http.Request) {
	var body types.RegisterBody

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		types.SendError(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	err = validator.ValidateRegister(h.adminToken, body.Token, body.Login, body.Password)
	if err != nil {
		types.SendError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	err = h.authService.Register(r.Context(), body.Login, body.Password)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			types.SendError(w, r, http.StatusConflict, err.Error())
			return
		}
		log.Println(err)
		types.SendError(w, r, http.StatusInternalServerError, "internal server error")
		return
	}

	types.SendResponse(w, r, map[string]string{
		"login": body.Login,
	})
}

func (h *AuthHandler) login(w http.ResponseWriter, r *http.Request) {
	var body types.LoginBody

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		types.SendError(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	err = validator.ValidateLogin(body.Login, body.Password)
	if err != nil {
		types.SendError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.authService.Login(r.Context(), body.Login, body.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			types.SendError(w, r, http.StatusUnauthorized, "invalid login or password")
			return
		}
		log.Println(err)
		types.SendError(w, r, http.StatusInternalServerError, "internal server error")
		return
	}

	types.SendResponse(w, r, map[string]string{
		"token": token,
	})
}

func (h *AuthHandler) logout(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")
	if token == "" {
		types.SendError(w, r, http.StatusBadRequest, "token parameter is required")
		return
	}

	err := h.authService.Logout(r.Context(), token)
	if err != nil {
		log.Println(err)
		types.SendError(w, r, http.StatusInternalServerError, "failed to logout")
		return
	}
	types.SendResponse(w, r,
		map[string]bool{
			"token": true,
		})
}
