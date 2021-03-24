package main

import (
	"backend/api/domain"
	_restaurantDelivery "backend/api/restaurant/delivery/http"
	_restaurantRepo "backend/api/restaurant/repository"
	_restaurantUsecase "backend/api/restaurant/usecase"
	_sessionRepo "backend/api/session/repository"
	_sessionUcase "backend/api/session/usecase"
	_userDelivery "backend/api/user/delivery/http"
	_userRepo "backend/api/user/repository"
	_userUcase "backend/api/user/usecase"
	"backend/config"
	"database/sql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func main() {
	e := echo.New()
	e.Static("/static/avatar", "static/avatar")

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{domain.Host + ":3000", "http://127.0.0.1:3000"},
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	}))

	dsn := ""
	db, err := sql.Open(config.PostgresDB, dsn)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(10) // TODO mn 10 соединений до бд
	db.SetMaxIdleConns(10) // TODO mn

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

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
	restaurantRepo := _restaurantRepo.NewRestaurantRepo()
	imageRepo := _userRepo.NewImageRepo()
	userUcase := _userUcase.NewUserUsecase(userRepo)
	sessionUcase := _sessionUcase.NewSessionUsecase(sessionRepo)
	restaurantUsecase := _restaurantUsecase.NewRestaurantUsecase(restaurantRepo)
	imageUcase := _userUcase.NewImageUsecase(imageRepo)

	_userDelivery.NewUserHandler(e, userUcase, sessionUcase, imageUcase)
	_restaurantDelivery.NewRestaurantHandler(e, restaurantUsecase)

	e.Logger.Fatal(e.Start(":5000"))
}
