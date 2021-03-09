package models

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type ServerError struct {
	Description  error
	Code int `json:"code"`
	Message string `json:"message"`
}

func AuthorizationErr() ServerError {
	return ServerError{
		Code: http.StatusUnauthorized,
		Message: "not authorized",
	}
}

func RegistrationErr() ServerError {
	return ServerError{
		Code: http.StatusUnauthorized,
		Message: "not auth",
	}
}

func (err* ServerError) AddError(desc error) {
	err.Description = desc
}

func (err ServerError) Response(c echo.Context) error {
	return c.JSON(http.StatusOK, err)
}
