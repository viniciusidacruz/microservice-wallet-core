package web

import (
	"encoding/json"
	"net/http"

	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/apperror"
)

type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if payload == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		return
	}
}

func respondError(w http.ResponseWriter, err error) {
	respondJSON(w, apperror.HTTPStatus(err), ErrorResponse{
		Error: apperror.MessageFromError(err),
		Code:  string(apperror.CodeFromError(err)),
	})
}
