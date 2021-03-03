package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// тут будет загрузка главной страницы с ресторанами
func mainPage(c echo.Context) error {
	return c.String(http.StatusOK, "It will be the main page")
}

func main() {
	var users []user // слайс пользователей
	var restourants []restourant // рестораны
	e := echo.New()
	e.GET("/", mainPage)
	// TODO: 
	//e.GET("/:id", restaurantPage) // урл на получение странички ресторана номер id
	//e.POST("/signup", createUser) // урл на регистрацию пользователя
	//e.POST("/signin", logUser) // урл на авторизацию
	//e.POST("/edituser", updateUser) // обновить пользователя после редактирования профиля

	e.Logger.Fatal(e.Start(":5000"))
}

