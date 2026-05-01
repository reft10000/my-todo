package domain

import "net/http"

type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string { return e.Message }

var (
	ErrInvalidTitle = &AppError{Code: http.StatusBadRequest, Message: "invalid title"}
	ErrBadRequest   = &AppError{Code: http.StatusBadRequest, Message: "bad request"}
)
