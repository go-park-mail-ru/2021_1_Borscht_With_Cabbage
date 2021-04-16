package errors

import (
	"net/http"
)

const ServerErrorCode = 420

type CustomError struct {
	SendError
	Description string
}

type SendError struct {
	Code    int    `json:"code"`
	Message string `json:"data"`
}

func (err *CustomError) Error() string {
	return err.Description
}

func (err *CustomError) SetDescription(description string) *CustomError {
	err.Description = description
	return err
}

func NewErrorWithMessage(message string) *CustomError {
	return &CustomError{
		SendError: SendError{
			Code:    ServerErrorCode,
			Message: message,
		},
		Description: message,
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
