package list_accounts

import (
	"time"

	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/apperror"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/gateway"
)

type ClientSummaryDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type AccountItemDTO struct {
	ID        string           `json:"id"`
	ClientID  string           `json:"client_id"`
	Balance   float64          `json:"balance"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	Client    ClientSummaryDTO `json:"client"`
}

type ListAccountsOutputDTO struct {
	Accounts []AccountItemDTO `json:"accounts"`
}

type ListAccountsUseCase struct {
	AccountGateway gateway.AccountGateway
}

func NewListAccountsUseCase(accountGateway gateway.AccountGateway) *ListAccountsUseCase {
	return &ListAccountsUseCase{
		AccountGateway: accountGateway,
	}
}

func (u *ListAccountsUseCase) Execute() (*ListAccountsOutputDTO, error) {
	accounts, err := u.AccountGateway.List()
	if err != nil {
		return nil, apperror.NewInternal("failed to list accounts", err)
	}

	items := make([]AccountItemDTO, 0, len(accounts))
	for _, account := range accounts {
		items = append(items, AccountItemDTO{
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
		})
	}

	return &ListAccountsOutputDTO{
		Accounts: items,
	}, nil
}
