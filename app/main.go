package main

import (
	"database/sql"
	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/restaurant"
	_restaurantDelivery "github.com/borscht/backend/internal/restaurant/delivery/http"
	_restaurantRepo "github.com/borscht/backend/internal/restaurant/repository"
	_restaurantUsecase "github.com/borscht/backend/internal/restaurant/usecase"
	_sessionRepo "github.com/borscht/backend/internal/session/repository"
	_sessionUcase "github.com/borscht/backend/internal/session/usecase"
	"github.com/borscht/backend/internal/user"
	_userDelivery "github.com/borscht/backend/internal/user/delivery/http"
	_userRepo "github.com/borscht/backend/internal/user/repository"
	_userUcase "github.com/borscht/backend/internal/user/usecase"
	_middleware "github.com/borscht/backend/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func route(e *echo.Echo, user user.UserHandler, restaurant restaurant.RestaurantHandler) {
	e.POST("/signin", user.Login)
	e.POST("/signup", user.Create)
	e.GET("/user", user.GetUserData)
	e.PUT("/user", user.EditProfile)
	e.GET("/auth", user.CheckAuth)
	e.GET("/logout", user.Logout)
	e.GET("/:id", restaurant.GetRestaurantPage)
	e.GET("/", restaurant.GetVendor)
	e.GET("/restaurants", restaurant.GetVendor)
}

func main() {
	e := echo.New()
	e.Static("/static/avatar", config.Static)

	// TODO убрать мидлвар в отедльный файл
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{config.Client, "http://127.0.0.1:3000"},
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	}))

	//dsn := "jdbc:postgresql://localhost:5432/postgres?user=labzunova&password=1111"
	dsn := "user=labzunova password=1111 dbname=postgres"
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

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &models.CustomContext{Context: c}
			return next(cc)
		}
	})

	userRepo := _userRepo.NewUserRepo(db)
	sessionRepo := _sessionRepo.NewSessionRepo(db)
	restaurantRepo := _restaurantRepo.NewRestaurantRepo(db)
	userUcase := _userUcase.NewUserUsecase(userRepo)
	sessionUcase := _sessionUcase.NewSessionUsecase(sessionRepo)
	restaurantUsecase := _restaurantUsecase.NewRestaurantUsecase(restaurantRepo)

	userHandler := _userDelivery.NewUserHandler(userUcase, sessionUcase)
	restaurantHandler := _restaurantDelivery.NewRestaurantHandler(restaurantUsecase)

	initMiddleware := _middleware.InitMiddleware(userUcase, sessionUcase)
	e.Use(initMiddleware.Auth)

	route(e, userHandler, restaurantHandler)

	e.Logger.Fatal(e.Start(":5000"))
}
