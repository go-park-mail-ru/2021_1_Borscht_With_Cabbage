package models

import (
	"github.com/labstack/echo/v4"
	"go/types"
	"net/http"
)

type ServerError struct {
	error
	Code int `json:"code"`
	Message string `json:"message"`
	Description string
}

type Error *ServerError

func BadRequest(err error) Error {
	serverError := ServerError{
		Code: http.StatusBadRequest,
		Message: "Bad request",
		Description: err.Error(),
	}

	return &serverError
}

func FailServer(err error) Error {
	serverError := ServerError{
		Code: http.StatusInternalServerError,
		Message: "Server error",
		Description: err.Error(),
	}

	return &serverError
}

func addData(data interface{}) interface{} {
	t := {

	}
}

func NotError(c echo.Context, data interface{}) error {
	t := ServerError{

	}
	return c.JSON(http.StatusOK, data, t)

}

func AuthorizationErr(err error) Error {
	serverError := ServerError{
		Code: http.StatusUnauthorized,
		Message: "Not authorized",
		Description: err.Error(),
	}

	return &serverError
}

func RegistrationErr() ServerError {
	return ServerError{
		Code: http.StatusUnauthorized,
		Message: "not auth",
	}
}

func (err ServerError) Response(c echo.Context) error {
	return c.JSON(http.StatusOK, err)
}
