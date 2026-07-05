package create_client

import (
	"database/sql"
	"testing"

	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/apperror"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ClientGatewayMock struct {
	mock.Mock
}

func (m *ClientGatewayMock) Save(client *entity.Client) error {
	args := m.Called(client)
	return args.Error(0)
}

func (m *ClientGatewayMock) Get(id string) (*entity.Client, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.Client), args.Error(1)
}

func (m *ClientGatewayMock) GetByEmail(email string) (*entity.Client, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.Client), args.Error(1)
}

func (m *ClientGatewayMock) List() ([]*entity.Client, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*entity.Client), args.Error(1)
}

func (m *ClientGatewayMock) Delete(id string) error {
	return m.Called(id).Error(0)
}

func (m *ClientGatewayMock) DeleteAll() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func TestCreateClientUseCase(t *testing.T) {
	m := &ClientGatewayMock{}

	m.On("GetByEmail", "test").Return((*entity.Client)(nil), sql.ErrNoRows)
	m.On("Save", mock.Anything).Return(nil)

	uc := NewCreateClientUseCase(m)
	output, err := uc.Execute(CreateClientInputDTO{
		Name:  "John",
		Email: "test",
	})

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.ID)
	assert.Equal(t, "John", output.Name)
	assert.Equal(t, "test", output.Email)
	m.AssertExpectations(t)
	m.AssertNumberOfCalls(t, "Save", 1)
}

func TestCreateClientUseCase_ValidationError(t *testing.T) {
	uc := NewCreateClientUseCase(&ClientGatewayMock{})

	output, err := uc.Execute(CreateClientInputDTO{})

	assert.Nil(t, output)
	assert.Error(t, err)
	assert.Equal(t, apperror.CodeValidation, apperror.CodeFromError(err))
	assert.Equal(t, "Name is required", err.Error())
}

func TestCreateClientUseCase_SaveError(t *testing.T) {
	m := &ClientGatewayMock{}
	m.On("GetByEmail", "test@test.com").Return((*entity.Client)(nil), sql.ErrNoRows)
	m.On("Save", mock.Anything).Return(sql.ErrConnDone)

	uc := NewCreateClientUseCase(m)
	output, err := uc.Execute(CreateClientInputDTO{
		Name:  "John",
		Email: "test@test.com",
	})

	assert.Nil(t, output)
	assert.Error(t, err)
	assert.Equal(t, apperror.CodeInternal, apperror.CodeFromError(err))
	assert.Equal(t, "failed to save client", err.Error())
}

func TestCreateClientUseCase_DuplicateEmail(t *testing.T) {
	existingClient, _ := entity.NewClient("John", "test@test.com")
	m := &ClientGatewayMock{}
	m.On("GetByEmail", "test@test.com").Return(existingClient, nil)

	uc := NewCreateClientUseCase(m)
	output, err := uc.Execute(CreateClientInputDTO{
		Name:  "Jane",
		Email: "test@test.com",
	})

	assert.Nil(t, output)
	assert.Error(t, err)
	assert.Equal(t, apperror.CodeConflict, apperror.CodeFromError(err))
	assert.Equal(t, "client with this email already exists", err.Error())
}
