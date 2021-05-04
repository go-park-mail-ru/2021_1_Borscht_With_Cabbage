package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/chat"
	chatDelivery "github.com/borscht/backend/internal/chat/delivery/http"
	chatRepo "github.com/borscht/backend/internal/chat/repository"
	chatUsecase "github.com/borscht/backend/internal/chat/usecase"
	imageRepo "github.com/borscht/backend/internal/image/repository"
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
	"github.com/borscht/backend/utils/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	serviceChat "github.com/borscht/backend/internal/services/chat"
	protoChat "github.com/borscht/backend/services/proto/chat"

	"github.com/gomodule/redigo/redis"
)

type initRoute struct {
	e               *echo.Echo
	user            user.UserHandler
	restaurant      restaurant.RestaurantHandler
	restaurantAdmin restaurantAdmin.AdminRestaurantHandler
	dishAdmin       restaurantAdmin.AdminDishHandler
	sectionAdmin    restaurantAdmin.AdminSectionHandler
	order           order.OrderHandler
	chat            chat.ChatHandler
	authMiddleware  custMiddleware.AuthMiddleware
	userMiddleware  custMiddleware.UserAuthMiddleware
	adminMiddleware custMiddleware.AdminAuthMiddleware
	wsMiddleware    custMiddleware.WsAuthMiddleware
}

func route(data initRoute) {
	userGroup := data.e.Group("/user", data.userMiddleware.Auth)
	auth := data.e.Group("", data.authMiddleware.Auth)
	userGroup.GET("", data.user.GetUserData)
	userGroup.PUT("", data.user.UpdateData)
	userGroup.PUT("/avatar", data.user.UploadAvatar)
	userGroup.POST("/address", data.user.UpdateMainAddress)
	userGroup.GET("/address", data.user.GetMainAddress)
	auth.GET("/auth", data.user.CheckAuth)
	auth.GET("/connect/ws", data.chat.GetKey)
	auth.GET("/chats", data.chat.GetAllChats)
	auth.GET("/chat/:id", data.chat.GetAllMessages)
	data.e.GET("/ws/:key", data.chat.Connect, data.wsMiddleware.WsAuth)

	restaurantGroup := data.e.Group("/restaurant", data.adminMiddleware.Auth)
	restaurantGroup.POST("/dish", data.dishAdmin.AddDish)
	restaurantGroup.DELETE("/dish", data.dishAdmin.DeleteDish)
	restaurantGroup.PUT("/dish", data.dishAdmin.UpdateDishData)
	restaurantGroup.PUT("/dish/image", data.dishAdmin.UploadDishImage)
	restaurantGroup.GET("/dishes", data.dishAdmin.GetAllDishes)
	restaurantGroup.PUT("", data.restaurantAdmin.UpdateRestaurantData)
	restaurantGroup.PUT("/avatar", data.restaurantAdmin.UploadRestaurantImage)
	restaurantGroup.POST("/section", data.sectionAdmin.AddSection)
	restaurantGroup.DELETE("/section", data.sectionAdmin.DeleteSection)
	restaurantGroup.PUT("/section", data.sectionAdmin.UpdateSection)

	data.e.POST("/signin", data.user.Login)
	data.e.POST("/signup", data.user.Create)
	data.e.POST("/restaurant/signup", data.restaurantAdmin.CreateRestaurant)
	data.e.POST("/restaurant/signin", data.restaurantAdmin.Login)
	userGroup.GET("/orders", data.order.GetUserOrders)
	userGroup.POST("/order", data.order.Create)
	userGroup.PUT("/basket", data.order.AddToBasket)
	userGroup.GET("/basket", data.order.GetBasket)
	userGroup.POST("/basket", data.order.AddBasket)
	data.e.GET("/logout", data.user.Logout)
	data.e.GET("/:id", data.restaurant.GetRestaurantPage)
	data.e.GET("/", data.restaurant.GetVendor)
	data.e.GET("/restaurants", data.restaurant.GetVendor)
}

func initServer(e *echo.Echo) {
	e.Static("/static", config.Static)
	e.Static("/", "../public")
	e.Static("/default", config.DefaultStatic)

	logger.InitLogger()
	e.Use(custMiddleware.LogMiddleware)
	e.Use(custMiddleware.CORS)
	// e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
	// 	TokenLookup: "header:X-XSRF-TOKEN",
	// }))

	e.Use(middleware.Secure())

	e.HTTPErrorHandler = custMiddleware.ErrorHandler
}

func main() {
	e := echo.New()
	initServer(e)
	grpcConnChat, errr := grpc.Dial(
		config.ChatServiceAddress,
		grpc.WithInsecure(),
	)
	if errr != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer grpcConnChat.Close()

	chatClient := protoChat.NewChatClient(grpcConnChat)
	chatService := serviceChat.NewService(chatClient)

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
	adminRestaurantRepo := restaurantAdminRepo.NewRestaurantRepo(db)
	adminDishRepo := restaurantAdminRepo.NewDishRepo(db)
	adminSectionRepo := restaurantAdminRepo.NewSectionRepo(db)
	restaurantRepo := restaurantRepo.NewRestaurantRepo(db)
	imageRepo := imageRepo.NewImageRepo()
	chatRepo := chatRepo.NewChattRepo(db)

	userUcase := userUcase.NewUserUsecase(userRepo, imageRepo)
	orderRepo := repository.NewOrderRepo(db)
	sessionUcase := sessionUcase.NewSessionUsecase(sessionRepo)
	adminRestaurantUsecase := restaurantAdminUsecase.NewRestaurantUsecase(adminRestaurantRepo, imageRepo)
	adminDishUsecase := restaurantAdminUsecase.NewDishUsecase(adminDishRepo, adminSectionRepo, imageRepo)
	adminSectionUsecase := restaurantAdminUsecase.NewSectionUsecase(adminSectionRepo)
	restaurantUsecase := restaurantUsecase.NewRestaurantUsecase(restaurantRepo, adminRestaurantRepo)
	orderUsecase := usecase.NewOrderUsecase(orderRepo, adminRestaurantRepo)
	chatUsecase := chatUsecase.NewChatUsecase(chatRepo, chatService)

	userHandler := userDelivery.NewUserHandler(userUcase, adminRestaurantUsecase, sessionUcase)
	adminRestaurantHandler := restaurantAdminDelivery.NewRestaurantHandler(adminRestaurantUsecase, sessionUcase)
	adminDishHandler := restaurantAdminDelivery.NewDishHandler(adminDishUsecase)
	adminSectionHandler := restaurantAdminDelivery.NewSectionHandler(adminSectionUsecase)
	restaurantHandler := restaurantDelivery.NewRestaurantHandler(restaurantUsecase)
	orderHandler := http.NewOrderHandler(orderUsecase)
	chatHandler := chatDelivery.NewChatHandler(chatUsecase, sessionUcase)

	initUserMiddleware := custMiddleware.InitUserMiddleware(userUcase, sessionUcase)
	initAdminMiddleware := custMiddleware.InitAdminMiddleware(adminRestaurantUsecase, sessionUcase)
	initAuthMiddleware := custMiddleware.InitAuthMiddleware(userUcase, adminRestaurantUsecase, sessionUcase)
	initWsMiddleware := custMiddleware.InitWsMiddleware(userUcase, adminRestaurantUsecase, sessionUcase)

	route(initRoute{
		e:               e,
		user:            userHandler,
		restaurantAdmin: adminRestaurantHandler,
		dishAdmin:       adminDishHandler,
		sectionAdmin:    adminSectionHandler,
		restaurant:      restaurantHandler,
		order:           orderHandler,
		chat:            chatHandler,
		userMiddleware:  *initUserMiddleware,
		adminMiddleware: *initAdminMiddleware,
		authMiddleware:  *initAuthMiddleware,
		wsMiddleware:    *initWsMiddleware,
	})

	e.Logger.Fatal(e.Start(":5000"))
}
