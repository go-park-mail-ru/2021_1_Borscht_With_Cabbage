package main

import (
	"backend/api"
	"backend/api/auth"
	"github.com/labstack/echo/v4"
	"net/http"
)

// тут будет загрузка главной страницы с ресторанами
func mainPage(c echo.Context) error {
	return c.String(http.StatusOK, "It will be the main page")
}

func router(e *echo.Echo) {
	e.GET("/", mainPage)
	e.POST("/signup", auth.CreateUser) // урл на регистрацию пользователя
	e.POST("/signin", auth.LogUser)    // урл на авторизацию
}

var Users []api.User
var restaurants []api.Restaurant

func main() {
	e := echo.New()
	router(e)

	e.Logger.Fatal(e.Start(":5000"))
}

