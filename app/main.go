package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/order"
	"github.com/borscht/backend/internal/order/delivery/http"
	"github.com/borscht/backend/internal/order/repository"
	"github.com/borscht/backend/internal/order/usecase"
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
	order           order.OrderHandler
	authMiddleware  custMiddleware.AuthMiddleware
	userMiddleware  custMiddleware.UserAuthMiddleware
	adminMiddleware custMiddleware.AdminAuthMiddleware
}

func route(data initRoute) {
	user := data.e.Group("/user", data.userMiddleware.Auth)
	// restaurant := data.e.Group("/restaurant", data.adminMiddleware.Auth)
	auth := data.e.Group("", data.authMiddleware.Auth)

	data.e.POST("/signin", data.user.Login)
	data.e.POST("/signup", data.user.Create)
	data.e.POST("/restaurant/signin", data.restaurantAdmin.Login)
	data.e.POST("/restaurant/signup", data.restaurantAdmin.Create)
	user.GET("", data.user.GetUserData)
	user.PUT("", data.user.EditProfile)
	user.GET("/orders", data.order.AddToBasket) // TODO TODO TODO НЕ ЭТО
	user.POST("/order", data.order.Create)
	user.PUT("/basket", data.order.AddToBasket)
	auth.GET("/auth", data.user.CheckAuth)
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
	orderRepo := repository.NewOrderRepo(db)

	userUcase := userUcase.NewUserUsecase(userRepo)
	sessionUcase := sessionUcase.NewSessionUsecase(sessionRepo)
	restaurantAdminUsecase := restaurantAdminUsecase.NewAdminUsecase(restaurantAdminRepo)
	restaurantUsecase := restaurantUsecase.NewRestaurantUsecase(restaurantRepo)
	orderUsecase := usecase.NewOrderUsecase(orderRepo)

	userHandler := userDelivery.NewUserHandler(userUcase, restaurantAdminUsecase, sessionUcase)
	restaurantAdminHandler := restaurantAdminDelivery.NewAdminHandler(restaurantAdminUsecase, sessionUcase)
	restaurantHandler := restaurantDelivery.NewRestaurantHandler(restaurantUsecase)
	orderHandler := http.NewOrderHandler(orderUsecase)

	initUserMiddleware := custMiddleware.InitUserMiddleware(userUcase, sessionUcase)
	initAdminMiddleware := custMiddleware.InitAdminMiddleware(restaurantAdminUsecase, sessionUcase)
	initAuthMiddleware := custMiddleware.InitAuthMiddleware(userUcase, restaurantAdminUsecase, sessionUcase)

	route(initRoute{
		e:               e,
		user:            userHandler,
		restaurantAdmin: restaurantAdminHandler,
		restaurant:      restaurantHandler,
		order:           orderHandler,
		userMiddleware:  *initUserMiddleware,
		adminMiddleware: *initAdminMiddleware,
		authMiddleware:  *initAuthMiddleware,
	})

	e.Logger.Fatal(e.Start(":5000"))
}
