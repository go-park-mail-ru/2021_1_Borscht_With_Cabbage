package errors

import (
	"net/http"
)

type CustomError struct {
	SendError
	Description string
}

type SendError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err *CustomError) Error() string {
	return err.Description
}

func NewCustomError(code int, mess string) *CustomError {
	return &CustomError{
		SendError: SendError{
			Code:    code,
			Message: mess,
		},
		Description: mess,
	}
}

func BadRequestError(desc string) *CustomError {
	return &CustomError{
		SendError: SendError{
			Code:    http.StatusBadRequest,
			Message: "bad request",
		},
		Description: desc,
	}
}

func FailServerError(desc string) *CustomError {
	return &CustomError{
		SendError: SendError{
			Code:    http.StatusInternalServerError,
			Message: "server error",
		},
		Description: desc,
	}
}

func AuthorizationError(desc string) *CustomError {
	return &CustomError{
		SendError: SendError{
			Code:    http.StatusUnauthorized,
			Message: "not authorized",
		},
		Description: desc,
	}
}
