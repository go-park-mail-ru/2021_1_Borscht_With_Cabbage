package models

import (
"github.com/labstack/echo/v4"
"net/http"
)

type ServerError struct {
	err  error
	body infoError
}

type infoError struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

func NewError(err error, code string) ServerError {
	info := infoError{
		Code:        code,
		Description: err.Error(),
	}

	return ServerError{
		err:  err,
		body: info,
	}
}

func (err ServerError) Response(c echo.Context) error {
	return c.JSON(http.StatusOK, err.body)
}
