package models

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type InfoError struct {
	Code        int `json:"code"`
	Message     string `json:"message"`
	Description string
}

func InternalError(err error) InfoError {
	return InfoError{
		Code: http.StatusInternalServerError,
		Message: "internal server error",
		Description: err.Error(),
	}
}

func UnAuthError(err error) InfoError {
	return InfoError{
		Code:        http.StatusUnauthorized,
		Message: "user not authorized",
		Description: err.Error(),
	}
}

func (err InfoError) Response(c echo.Context) error {
	return c.JSON(http.StatusOK, err)
}
