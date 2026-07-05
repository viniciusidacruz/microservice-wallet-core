package apperror

import (
	"errors"
	"net/http"
)

type Code string

const (
	CodeValidation Code = "VALIDATION_ERROR"
	CodeNotFound   Code = "NOT_FOUND"
	CodeConflict   Code = "CONFLICT"
	CodeInternal   Code = "INTERNAL_ERROR"
)

type AppError struct {
	Code    Code
	Message string
	Err     error
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewValidation(message string) *AppError {
	return &AppError{
		Code:    CodeValidation,
		Message: message,
	}
}

func NewNotFound(message string) *AppError {
	return &AppError{
		Code:    CodeNotFound,
		Message: message,
	}
}

func NewConflict(message string) *AppError {
	return &AppError{
		Code:    CodeConflict,
		Message: message,
	}
}

func NewInternal(message string, err error) *AppError {
	return &AppError{
		Code:    CodeInternal,
		Message: message,
		Err:     err,
	}
}

func HTTPStatus(err error) int {
	var appErr *AppError
	if errors.As(err, &appErr) {
		switch appErr.Code {
		case CodeValidation:
			return http.StatusBadRequest
		case CodeNotFound:
			return http.StatusNotFound
		case CodeConflict:
			return http.StatusConflict
		default:
			return http.StatusInternalServerError
		}
	}

	return http.StatusInternalServerError
}

func CodeFromError(err error) Code {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code
	}

	return CodeInternal
}

func MessageFromError(err error) string {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Message
	}

	return "internal server error"
}
