package main

import (
	"backend/api/page"
	"github.com/labstack/echo/v4"
)

func router(e *echo.Echo)  {
	e.GET("/", page.MainPage)
	e.GET("/vendor", page.GetVendor) // урл на получение json с растаранами
	//e.GET("/:id", restaurantPage) // урл на получение странички ресторана номер id
	//e.POST("/signup", createUser) // урл на регистрацию пользователя
	//e.POST("/signin", logUser) // урл на авторизацию
	//e.POST("/edituser", updateUser) // обновить пользователя после редактирования профиля
}

func main() {
	e := echo.New()
	router(e)

	e.Logger.Fatal(e.Start(":5000"))
}

