package main

import (
	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/restaurant"
	restaurantDelivery "github.com/borscht/backend/internal/restaurant/delivery/http"
	restaurantRepo "github.com/borscht/backend/internal/restaurant/repository"
	restaurantUsecase "github.com/borscht/backend/internal/restaurant/usecase"
	sessionRepo "github.com/borscht/backend/internal/session/repository"
	sessionUcase "github.com/borscht/backend/internal/session/usecase"
	"github.com/borscht/backend/internal/user"
	userDelivery "github.com/borscht/backend/internal/user/delivery/http"
	userRepo "github.com/borscht/backend/internal/user/repository"
	userUcase "github.com/borscht/backend/internal/user/usecase"
	custMiddleware "github.com/borscht/backend/middleware"
	"github.com/labstack/echo/v4"
)

type initRoute struct {
	e          *echo.Echo
	user       user.UserHandler
	restaurant restaurant.RestaurantHandler
	middleware custMiddleware.AuthMiddleware
}

func route(data initRoute) {
	user := data.e.Group("/user", data.middleware.Auth)

	data.e.POST("/signin", data.user.Login)
	data.e.POST("/signup", data.user.Create)
	user.GET("", data.user.GetUserData)
	user.PUT("", data.user.EditProfile)
	data.e.GET("/auth", data.user.CheckAuth)
	data.e.GET("/logout", data.user.Logout)
	data.e.GET("/:id", data.restaurant.GetRestaurantPage)
	data.e.GET("/", data.restaurant.GetVendor)
	data.e.GET("/restaurants", data.restaurant.GetVendor)
}

func main() {
	e := echo.New()
	e.Static("/static/avatar", config.Static)

	// TODO убрать мидлвар в отедльный файл
	e.Use(custMiddleware.CORS)

	userRepo := userRepo.NewUserRepo()
	sessionRepo := sessionRepo.NewSessionRepo()
	restaurantRepo := restaurantRepo.NewRestaurantRepo()
	userUcase := userUcase.NewUserUsecase(userRepo)
	sessionUcase := sessionUcase.NewSessionUsecase(sessionRepo)
	restaurantUsecase := restaurantUsecase.NewRestaurantUsecase(restaurantRepo)

	userHandler := userDelivery.NewUserHandler(userUcase, sessionUcase)
	restaurantHandler := restaurantDelivery.NewRestaurantHandler(restaurantUsecase)

	initMiddleware := custMiddleware.InitMiddleware(userUcase, sessionUcase)

	route(initRoute{
		e:          e,
		user:       userHandler,
		restaurant: restaurantHandler,
		middleware: *initMiddleware,
	})

	e.Logger.Fatal(e.Start(":5000"))
}
