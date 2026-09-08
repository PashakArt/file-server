package types

import (
	"encoding/json"
	"net/http"
)

type APIError struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

type APIResponse struct {
	Error    *APIError   `json:"error,omitempty"`
	Response interface{} `json:"response,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}

func prepareResponse(w http.ResponseWriter, r *http.Request, code int, resp APIResponse) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(resp)
}

func SendError(w http.ResponseWriter, r *http.Request, code int, errText string) {
	prepareResponse(w, r, code, APIResponse{
		Error: &APIError{
			Code: code,
			Text: errText,
		},
	})
}

func SendData(w http.ResponseWriter, r *http.Request, payload interface{}) {
	prepareResponse(w, r, http.StatusOK, APIResponse{
		Response: payload,
	})
}
