package get_account

import (
	"database/sql"
	"errors"
	"time"

	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/apperror"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/gateway"
)

type ClientSummaryDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type GetAccountOutputDTO struct {
	ID        string           `json:"id"`
	ClientID  string           `json:"client_id"`
	Balance   float64          `json:"balance"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	Client    ClientSummaryDTO `json:"client"`
}

type GetAccountUseCase struct {
	AccountGateway gateway.AccountGateway
}

func NewGetAccountUseCase(accountGateway gateway.AccountGateway) *GetAccountUseCase {
	return &GetAccountUseCase{
		AccountGateway: accountGateway,
	}
}

func (u *GetAccountUseCase) Execute(id string) (*GetAccountOutputDTO, error) {
	if id == "" {
		return nil, apperror.NewValidation("account id is required")
	}

	account, err := u.AccountGateway.FindByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFound("account not found")
		}

		return nil, apperror.NewInternal("failed to get account", err)
	}

	return &GetAccountOutputDTO{
		ID:        account.ID,
		ClientID:  account.Client.ID,
		Balance:   account.Balance,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
		Client: ClientSummaryDTO{
			ID:    account.Client.ID,
			Name:  account.Client.Name,
			Email: account.Client.Email,
		},
	}, nil
}
