package models

import (
	"context"
	"net/http"

	"github.com/borscht/backend/utils"
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
	ctx := GetContext(c)

	utils.InfoLog(ctx, utils.Fields{
		"code":     http.StatusFound,
		"response": "redirect '/login'",
	})

	return c.JSON(http.StatusOK, redirect{
		Code:     http.StatusFound,
		Redirect: "/login",
	})
}

func GetContext(c echo.Context) context.Context {
	ctx := c.Request().Context()

	ctx = context.WithValue(ctx, "User", c.Get("User"))
	return context.WithValue(ctx, "request_id", c.Get("request_id"))
}

func SendResponse(c echo.Context, data interface{}) error {

	serverMessage := message{http.StatusOK, data}
	ctx := GetContext(c)

	utils.InfoLog(ctx, utils.Fields{
		"code":     http.StatusOK,
		"response": serverMessage,
	})

	return c.JSON(http.StatusOK, serverMessage)
}

func SendResponseWithError(c echo.Context, err error) error {
	// проверяем можно ли преобразовать в кастомную ошибку
	if customErr, ok := err.(*errors.CustomError); ok {
		return c.JSON(http.StatusOK, customErr.SendError)
	}

	ctx := GetContext(c)
	customErr := errors.FailServer(ctx, err.Error())

	utils.InfoLog(ctx, utils.Fields{
		"code":     customErr.Code,
		"response": customErr.SendError,
	})

	return c.JSON(http.StatusOK, customErr.SendError)
}
