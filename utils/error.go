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

func BadRequest(desc string) *CustomError {
	return &CustomError{
		SendError: SendError{
			Code:    http.StatusBadRequest,
			Message: "bad request",
		},
		Description: desc,
	}
}

func FailServer(desc string) *CustomError {
	return &CustomError{
		SendError: SendError{
			Code:    http.StatusInternalServerError,
			Message: "server error",
		},
		Description: desc,
	}
}

func Authorization(desc string) *CustomError {
	return &CustomError{
		SendError: SendError{
			Code:    http.StatusUnauthorized,
			Message: "not authorized",
		},
		Description: desc,
	}
}
