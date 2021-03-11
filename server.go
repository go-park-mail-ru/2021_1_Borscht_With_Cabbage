package main

import (
	"backend/api"
	"backend/api/auth"
	"backend/api/page"
	"backend/api/profile"
	"backend/api/restaurant"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func router(e *echo.Echo) {
	e.POST("/signin", auth.LoginUser)
	e.POST("/signup", auth.CreateUser)
	e.GET("/:id", restaurant.GetRestaurantPage)
	e.GET("/user", profile.GetUserData)
	e.PUT("/user", profile.EditProfile)
	e.GET("/", page.GetVendor)
	e.GET("/restaurants", page.GetVendor)
	e.GET("/auth", auth.CheckAuth)
	e.GET("/logout", auth.LogoutUser)
}

func main() {
	e := echo.New()
	e.Static("/static/avatar", "static/avatar")

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{api.Host + ":3000"},
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	}))

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
