package main

import (
	"backend/api/auth"
	"backend/api/restaurant"
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
	e.GET("/:id", restaurant.RestaurantPage)    // урл на страницу ресторана
}

func main() {
	e := echo.New()
	router(e)

	e.Logger.Fatal(e.Start(":5000"))
}
