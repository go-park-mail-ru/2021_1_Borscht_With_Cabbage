package errors

import (
	"net/http"
)

type CustomError struct {
	SendError
	Description string
}

type SendError struct {
	Code int `json:"code"`
	Message string `json:"message"`
}

func (err *CustomError) Error() string {
	return err.Description
}

func BadRequest(err error) *CustomError {
	cErr := CustomError{}

	cErr.Code = http.StatusBadRequest
	cErr.Message = "Bad request"
	cErr.Description = err.Error()

	return &cErr
}

func FailServer(err error) *CustomError {
	cErr := CustomError{}

	cErr.Code = http.StatusInternalServerError
	cErr.Message = "Server error"
	cErr.Description = err.Error()

	return &cErr
}

func Authorization(err error) *CustomError {
	cErr := CustomError{}

	cErr.Code = http.StatusUnauthorized
	cErr.Message = "Not authorized"
	cErr.Description = err.Error()

	return &cErr
}
