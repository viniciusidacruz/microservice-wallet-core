package apperror

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppError(t *testing.T) {
	validationErr := NewValidation("client_id is required")
	assert.Equal(t, "client_id is required", validationErr.Error())
	assert.Equal(t, http.StatusBadRequest, HTTPStatus(validationErr))
	assert.Equal(t, CodeValidation, CodeFromError(validationErr))
	assert.Equal(t, "client_id is required", MessageFromError(validationErr))

	notFoundErr := NewNotFound("client not found")
	assert.Equal(t, http.StatusNotFound, HTTPStatus(notFoundErr))
	assert.Equal(t, CodeNotFound, CodeFromError(notFoundErr))

	conflictErr := NewConflict("client with this email already exists")
	assert.Equal(t, http.StatusConflict, HTTPStatus(conflictErr))
	assert.Equal(t, CodeConflict, CodeFromError(conflictErr))

	internalErr := NewInternal("failed to save account", errors.New("db error"))
	assert.Equal(t, http.StatusInternalServerError, HTTPStatus(internalErr))
	assert.Equal(t, CodeInternal, CodeFromError(internalErr))
	assert.Equal(t, "failed to save account", MessageFromError(internalErr))

	assert.Equal(t, http.StatusInternalServerError, HTTPStatus(errors.New("unknown error")))
	assert.Equal(t, CodeInternal, CodeFromError(errors.New("unknown error")))
	assert.Equal(t, "internal server error", MessageFromError(errors.New("unknown error")))
}
