package main

import (
	"backend/api/domain"
	"backend/api/page"
	"backend/api/restaurant"
	_sessionRepo "backend/api/session/repository"
	_sessionUcase "backend/api/session/usecase"
	_userDelivery "backend/api/user/delivery/http"
	_userRepo "backend/api/user/repository"
	_userUcase "backend/api/user/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func router(e *echo.Echo) {
	e.GET("/:id", restaurant.GetRestaurantPage)
	e.GET("/", page.GetVendor)
	e.GET("/restaurants", page.GetVendor)
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

	userRepo := _userRepo.NewUserRepo()
	sessionRepo := _sessionRepo.NewSessionRepo()
	userUcase := _userUcase.NewUserUsecase(userRepo)
	sessionUcase := _sessionUcase.NewSessionUsecase(sessionRepo)

	_userDelivery.NewUserHandler(e, userUcase, sessionUcase)

	router(e)
	e.Logger.Fatal(e.Start(":5000"))
}
