package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/restaurant"
	restaurantDelivery "github.com/borscht/backend/internal/restaurant/delivery/http"
	restaurantRepo "github.com/borscht/backend/internal/restaurant/repository"
	restaurantUsecase "github.com/borscht/backend/internal/restaurant/usecase"
	"github.com/borscht/backend/internal/restaurantAdmin"
	restaurantAdminDelivery "github.com/borscht/backend/internal/restaurantAdmin/delivery/http"
	restaurantAdminRepo "github.com/borscht/backend/internal/restaurantAdmin/repository"
	restaurantAdminUsecase "github.com/borscht/backend/internal/restaurantAdmin/usecase"
	sessionRepo "github.com/borscht/backend/internal/session/repository"
	sessionUcase "github.com/borscht/backend/internal/session/usecase"
	"github.com/borscht/backend/internal/user"
	userDelivery "github.com/borscht/backend/internal/user/delivery/http"
	userRepo "github.com/borscht/backend/internal/user/repository"
	userUcase "github.com/borscht/backend/internal/user/usecase"
	custMiddleware "github.com/borscht/backend/middleware"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"

	"github.com/gomodule/redigo/redis"
)

type initRoute struct {
	e               *echo.Echo
	user            user.UserHandler
	restaurant      restaurant.RestaurantHandler
	restaurantAdmin restaurantAdmin.AdminHandler
	authMiddleware  custMiddleware.AuthMiddleware
	adminMiddleware custMiddleware.AdminAuthMiddleware
}

func route(data initRoute) {
	user := data.e.Group("/user", data.authMiddleware.Auth)

	data.e.POST("/signin", data.user.Login)
	data.e.POST("/signup", data.user.Create)
	data.e.POST("/restaurant/signin", data.restaurantAdmin.Login)
	data.e.POST("/restaurant/signup", data.restaurantAdmin.Create)
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

	e.Use(custMiddleware.CORS)

	// подключение postgres
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

	// подключение redis
	redisConn, err := redis.Dial("tcp", config.RedisHost)
	if err != nil {
		log.Fatal(err)
	}
	defer redisConn.Close()

	userRepo := userRepo.NewUserRepo(db)
	sessionRepo := sessionRepo.NewSessionRepo(redisConn)
	restaurantAdminRepo := restaurantAdminRepo.NewAdminRepo(db)
	restaurantRepo := restaurantRepo.NewRestaurantRepo(db)
	userUcase := userUcase.NewUserUsecase(userRepo)
	sessionUcase := sessionUcase.NewSessionUsecase(sessionRepo)
	restaurantAdminUsecase := restaurantAdminUsecase.NewAdminUsecase(restaurantAdminRepo)
	restaurantUsecase := restaurantUsecase.NewRestaurantUsecase(restaurantRepo)

	userHandler := userDelivery.NewUserHandler(userUcase, restaurantAdminUsecase, sessionUcase)
	restaurantAdminHandler := restaurantAdminDelivery.NewAdminHandler(restaurantAdminUsecase, sessionUcase)
	restaurantHandler := restaurantDelivery.NewRestaurantHandler(restaurantUsecase)

	initAuthMiddleware := custMiddleware.InitMiddleware(userUcase, sessionUcase)
	initAdminMiddleware := custMiddleware.InitAdminMiddleware(restaurantAdminUsecase, sessionUcase)

	route(initRoute{
		e:               e,
		user:            userHandler,
		restaurantAdmin: restaurantAdminHandler,
		restaurant:      restaurantHandler,
		authMiddleware:  *initAuthMiddleware,
		adminMiddleware: *initAdminMiddleware,
	})

	e.Logger.Fatal(e.Start(":5000"))
}
