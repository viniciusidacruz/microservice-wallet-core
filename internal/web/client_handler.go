package web

import (
	"encoding/json"
	"net/http"

	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/apperror"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/usecases/create_client"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/usecases/delete_all_clients"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/usecases/delete_client"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/usecases/get_client"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/usecases/list_clients"
	"github.com/go-chi/chi/v5"
)

type WebClientHandler struct {
	CreateClientUseCase     create_client.CreateClientUseCase
	GetClientUseCase        get_client.GetClientUseCase
	ListClientsUseCase      list_clients.ListClientsUseCase
	DeleteClientUseCase     delete_client.DeleteClientUseCase
	DeleteAllClientsUseCase delete_all_clients.DeleteAllClientsUseCase
}

func NewWebClientHandler(
	createClientUseCase create_client.CreateClientUseCase,
	getClientUseCase get_client.GetClientUseCase,
	listClientsUseCase list_clients.ListClientsUseCase,
	deleteClientUseCase delete_client.DeleteClientUseCase,
	deleteAllClientsUseCase delete_all_clients.DeleteAllClientsUseCase,
) *WebClientHandler {
	return &WebClientHandler{
		CreateClientUseCase:     createClientUseCase,
		GetClientUseCase:        getClientUseCase,
		ListClientsUseCase:      listClientsUseCase,
		DeleteClientUseCase:     deleteClientUseCase,
		DeleteAllClientsUseCase: deleteAllClientsUseCase,
	}
}

func (h *WebClientHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	var input create_client.CreateClientInputDTO
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		respondError(w, apperror.NewValidation("invalid request body"))
		return
	}

	output, err := h.CreateClientUseCase.Execute(input)
	if err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, output)
}

func (h *WebClientHandler) GetClient(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	output, err := h.GetClientUseCase.Execute(id)
	if err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, output)
}

func (h *WebClientHandler) ListClients(w http.ResponseWriter, r *http.Request) {
	output, err := h.ListClientsUseCase.Execute()
	if err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, output)
}

func (h *WebClientHandler) DeleteClient(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	output, err := h.DeleteClientUseCase.Execute(id)
	if err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, output)
}

func (h *WebClientHandler) DeleteAllClients(w http.ResponseWriter, r *http.Request) {
	output, err := h.DeleteAllClientsUseCase.Execute()
	if err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, output)
}
