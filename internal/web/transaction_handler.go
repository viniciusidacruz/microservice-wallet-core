package web

import (
	"encoding/json"
	"net/http"

	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/apperror"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/usecases/create_transaction"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/usecases/delete_all_transactions"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/usecases/delete_transaction"
	"github.com/go-chi/chi/v5"
)

type WebTransactionHandler struct {
	CreateTransactionUseCase     create_transaction.CreateTransactionUseCase
	DeleteTransactionUseCase     delete_transaction.DeleteTransactionUseCase
	DeleteAllTransactionsUseCase delete_all_transactions.DeleteAllTransactionsUseCase
}

func NewWebTransactionHandler(
	createTransactionUseCase create_transaction.CreateTransactionUseCase,
	deleteTransactionUseCase delete_transaction.DeleteTransactionUseCase,
	deleteAllTransactionsUseCase delete_all_transactions.DeleteAllTransactionsUseCase,
) *WebTransactionHandler {
	return &WebTransactionHandler{
		CreateTransactionUseCase:     createTransactionUseCase,
		DeleteTransactionUseCase:     deleteTransactionUseCase,
		DeleteAllTransactionsUseCase: deleteAllTransactionsUseCase,
	}
}

func (h *WebTransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var input create_transaction.CreateTransactionInputDTO
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		respondError(w, apperror.NewValidation("invalid request body"))
		return
	}

	output, err := h.CreateTransactionUseCase.Execute(input)
	if err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, output)
}

func (h *WebTransactionHandler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	output, err := h.DeleteTransactionUseCase.Execute(id)
	if err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, output)
}

func (h *WebTransactionHandler) DeleteAllTransactions(w http.ResponseWriter, r *http.Request) {
	output, err := h.DeleteAllTransactionsUseCase.Execute()
	if err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, output)
}
