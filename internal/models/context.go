package models

import (
	"net/http"

	errors "github.com/borscht/backend/utils"
	"github.com/labstack/echo/v4"
)

type message struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

type redirect struct {
	Code     int    `json:"code"`
	Redirect string `json:"redirect"`
}

func SendRedirectLogin(c echo.Context) error {
	return c.JSON(http.StatusOK, redirect{
		Code:     http.StatusFound,
		Redirect: "/login",
	})
}

func SendResponse(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, message{http.StatusOK, data})
}

func SendResponseWithError(c echo.Context, err error) error {
	// проверяем можно ли преобразовать в кастомную ошибку
	if customErr, ok := err.(*errors.CustomError); ok {
		return c.JSON(http.StatusOK, customErr.SendError)
	}

	customErr := errors.FailServer(err.Error())

	return c.JSON(http.StatusOK, customErr.SendError)
}
