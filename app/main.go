package main

import (
	"database/sql"
	"fmt"
	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/order"
	"github.com/borscht/backend/internal/order/delivery/http"
	"github.com/borscht/backend/internal/order/repository"
	"github.com/borscht/backend/internal/order/usecase"
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
	_ "github.com/lib/pq"
	"log"
)

type initRoute struct {
	e          *echo.Echo
	user       user.UserHandler
	restaurant restaurant.RestaurantHandler
	order      order.OrderHandler
	middleware custMiddleware.AuthMiddleware
}

func route(data initRoute) {
	user := data.e.Group("/user", data.middleware.Auth)

	data.e.POST("/signin", data.user.Login)
	data.e.POST("/signup", data.user.Create)
	user.GET("", data.user.GetUserData)
	user.PUT("", data.user.EditProfile)
	user.GET("/auth", data.user.CheckAuth)
	user.GET("/user/orders", data.order.GetUserOrder)
	data.e.GET("/logout", data.user.Logout)
	data.e.GET("/:id", data.restaurant.GetRestaurantPage)
	data.e.GET("/", data.restaurant.GetVendor)
	data.e.GET("/restaurants", data.restaurant.GetVendor)
}

func main() {
	e := echo.New()
	e.Static("/static", config.Static)
	e.Static("/default", config.DefaultStatic)

	e.Use(custMiddleware.CORS)

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s", config.DBUser, config.DBPass, config.DBName)
	db, err := sql.Open(config.PostgresDB, dsn)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(3)

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	userRepo := userRepo.NewUserRepo(db)
	sessionRepo := sessionRepo.NewSessionRepo(db)
	restaurantRepo := restaurantRepo.NewRestaurantRepo(db)
	orderRepo := repository.NewOrderRepo(db)

	userUcase := userUcase.NewUserUsecase(userRepo)
	sessionUcase := sessionUcase.NewSessionUsecase(sessionRepo)
	restaurantUsecase := restaurantUsecase.NewRestaurantUsecase(restaurantRepo)
	orderUsecase := usecase.NewOrderUsecase(orderRepo)

	userHandler := userDelivery.NewUserHandler(userUcase, sessionUcase)
	restaurantHandler := restaurantDelivery.NewRestaurantHandler(restaurantUsecase)
	orderHandler := http.NewOrderHandler(orderUsecase)

	initMiddleware := custMiddleware.InitMiddleware(userUcase, sessionUcase)

	route(initRoute{
		e:          e,
		user:       userHandler,
		restaurant: restaurantHandler,
		order:      orderHandler,
		middleware: *initMiddleware,
	})

	e.Logger.Fatal(e.Start(":5000"))
}
