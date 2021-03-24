package models

import (
	"context"
	_errors "github.com/borscht/backend/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

type CustomContext struct {
	echo.Context
	User *User
}

type message struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

func (c *CustomContext) GetContext() context.Context {
	ctx := c.Request().Context()
	if c.User != nil {
		ctx = context.WithValue(ctx, "User", c.User)
	}

	return ctx
}

func (c *CustomContext) SendResponse(data interface{}) error {
	return c.JSON(http.StatusOK, message{200, data})
}

func (c *CustomContext) SendResponseWithError(err error) error {
	// проверяем можно ли преобразовать в кастомную ошибку
	if customErr, ok := err.(*_errors.CustomError); ok {
		return c.JSON(http.StatusOK, customErr.SendError)
	}

	customErr := _errors.FailServer(err.Error())
	return c.JSON(http.StatusOK, customErr.SendError)
}
