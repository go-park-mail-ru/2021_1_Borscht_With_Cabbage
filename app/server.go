package app

import (
	"backend/api/auth"
	"backend/api/page"
	"backend/api/profile"
	"backend/api/restaurant"
	"backend/api/domain"
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
		AllowOrigins:     []string{domain.Host + ":3000", "http://127.0.0.1:3000"},
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	}))

	Users := make([]domain.User, 0)
	Sessions := make(map[string]string, 0)
	Restaurants := make(map[string]domain.Restaurant, 0)

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &domain.CustomContext{Context: c, Users: &Users, Restaurants: &Restaurants, Sessions: &Sessions}
			domain.InitData(*cc)
			return next(cc)
		}
	})

	router(e)
	e.Logger.Fatal(e.Start(":5000"))
}
