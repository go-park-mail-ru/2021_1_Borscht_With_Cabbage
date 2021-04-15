package models

import (
	"context"
	"net/http"

	errors "github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
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

type Response interface {
	Sanitize()
}

func SendRedirectLogin(c echo.Context) error {
	ctx := GetContext(c)

	logger.ResponseLevel().InfoLog(ctx, logger.Fields{
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
	ctx = context.WithValue(ctx, "Restaurant", c.Get("Restaurant"))
	return context.WithValue(ctx, "request_id", c.Get("request_id"))
}

func SendMoreResponse(c echo.Context, data ...Response) error {
	ctx := GetContext(c)

	for i := range data {
		if data[i] != nil {
			data[i].Sanitize()
		}
	}

	serverMessage := message{http.StatusOK, data}

	logger.ResponseLevel().InfoLog(ctx, logger.Fields{
		"code":     http.StatusOK,
		"response": data,
	})

	return c.JSON(http.StatusOK, serverMessage)
}

func SendResponse(c echo.Context, data Response) error {
	ctx := GetContext(c)

	if data != nil {
		data.Sanitize()
	}

	serverMessage := message{http.StatusOK, data}

	logger.ResponseLevel().InfoLog(ctx, logger.Fields{
		"code":     http.StatusOK,
		"response": data,
	})

	return c.JSON(http.StatusOK, serverMessage)
}

func SendResponseWithError(c echo.Context, err error) error {
	// проверяем можно ли преобразовать в кастомную ошибку
	if customErr, ok := err.(*errors.CustomError); ok {
		return c.JSON(http.StatusOK, customErr.SendError)
	}

	ctx := GetContext(c)
	customErr := errors.FailServerError(err.Error())

	logger.ResponseLevel().InfoLog(ctx, logger.Fields{
		"code":     customErr.Code,
		"response": customErr.SendError,
	})

	return c.JSON(http.StatusOK, customErr.SendError)
}
