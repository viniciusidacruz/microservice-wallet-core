package web

import (
	"encoding/json"
	"net/http"

	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/apperror"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/usecases/create_account"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/usecases/delete_account"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/usecases/delete_all_accounts"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/usecases/get_account"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/usecases/list_accounts"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/usecases/set_account_balance"
	"github.com/go-chi/chi/v5"
)

type WebAccountHandler struct {
	CreateAccountUseCase      create_account.CreateAccountUseCase
	SetAccountBalanceUseCase  set_account_balance.SetAccountBalanceUseCase
	GetAccountUseCase         get_account.GetAccountUseCase
	ListAccountsUseCase       list_accounts.ListAccountsUseCase
	DeleteAccountUseCase      delete_account.DeleteAccountUseCase
	DeleteAllAccountsUseCase delete_all_accounts.DeleteAllAccountsUseCase
}

func NewWebAccountHandler(
	createAccountUseCase create_account.CreateAccountUseCase,
	setAccountBalanceUseCase set_account_balance.SetAccountBalanceUseCase,
	getAccountUseCase get_account.GetAccountUseCase,
	listAccountsUseCase list_accounts.ListAccountsUseCase,
	deleteAccountUseCase delete_account.DeleteAccountUseCase,
	deleteAllAccountsUseCase delete_all_accounts.DeleteAllAccountsUseCase,
) *WebAccountHandler {
	return &WebAccountHandler{
		CreateAccountUseCase:      createAccountUseCase,
		SetAccountBalanceUseCase:  setAccountBalanceUseCase,
		GetAccountUseCase:         getAccountUseCase,
		ListAccountsUseCase:       listAccountsUseCase,
		DeleteAccountUseCase:      deleteAccountUseCase,
		DeleteAllAccountsUseCase:  deleteAllAccountsUseCase,
	}
}

func (h *WebAccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var input create_account.CreateAccountInputDTO
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		respondError(w, apperror.NewValidation("invalid request body"))
		return
	}

	output, err := h.CreateAccountUseCase.Execute(input)
	if err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, output)
}

func (h *WebAccountHandler) SetAccountBalance(w http.ResponseWriter, r *http.Request) {
	var input set_account_balance.SetAccountBalanceInputDTO
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		respondError(w, apperror.NewValidation("invalid request body"))
		return
	}

	output, err := h.SetAccountBalanceUseCase.Execute(input)
	if err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, output)
}

func (h *WebAccountHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	output, err := h.GetAccountUseCase.Execute(id)
	if err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, output)
}

func (h *WebAccountHandler) ListAccounts(w http.ResponseWriter, r *http.Request) {
	output, err := h.ListAccountsUseCase.Execute()
	if err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, output)
}

func (h *WebAccountHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	output, err := h.DeleteAccountUseCase.Execute(id)
	if err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, output)
}

func (h *WebAccountHandler) DeleteAllAccounts(w http.ResponseWriter, r *http.Request) {
	output, err := h.DeleteAllAccountsUseCase.Execute()
	if err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, output)
}
