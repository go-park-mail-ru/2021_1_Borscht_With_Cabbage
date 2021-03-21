package domain

import (
	errors "backend/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

type CustomContext struct {
	echo.Context
	User        *User
	Users       *[]User
	Restaurants *map[string]Restaurant // [id]RestaurantStruct
	Sessions    *map[string]string     // [session]user's phone number
}

type message struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

func (c *CustomContext) SendOK(data interface{}) error {
	return c.JSON(http.StatusOK, message{200, data})
}

func (c *CustomContext) SendERR(err error) error {
	// проверяем можно ли преобразовать в кастомную ошибку
	if customErr, ok := err.(*errors.CustomError); ok {
		return c.JSON(http.StatusOK, customErr.SendError)
	}

	customErr := errors.FailServer(err.Error())
	return c.JSON(http.StatusOK, customErr.SendError)
}
