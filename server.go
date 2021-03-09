package main

import (
	"backend/api"
	"backend/api/auth"
	"backend/api/image"
	"backend/api/page"
	"backend/api/profile"
	"backend/api/restaurant"
	"github.com/labstack/echo/v4"
)

func router(e *echo.Echo) {
	e.POST("/signin", auth.LoginUser)
	e.POST("/signup", auth.CreateUser)
	e.POST("/avatar", image.UploadAvatar)
	e.GET("/avatar", image.DownloadAvatar)
	e.GET("/:id", restaurant.GetRestaurantPage)
	e.GET("/user", profile.GetUserData)
	e.PUT("/user", profile.EditProfile)
	e.GET("/vendor", page.GetVendor)
	e.GET("/auth", auth.CheckAuth)
}

func main() {
	e := echo.New()

	Users := make([]api.User, 0)
	Sessions := make(map[string]string, 0)
	Restaurants := make(map[string]api.Restaurant, 0)

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &api.CustomContext{Context: c, Users: &Users, Restaurants: &Restaurants, Sessions: &Sessions}
			api.InitData(*cc)
			return next(cc)
		}
	})

	router(e)
	e.Logger.Fatal(e.Start(":5000"))
}
