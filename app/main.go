package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/borscht/backend/utils/notifications"
	"github.com/borscht/backend/utils/websocketPool"

	"github.com/borscht/backend/internal/services/auth"
	"github.com/borscht/backend/internal/services/basket"
	protoAuth "github.com/borscht/backend/services/proto/auth"
	protoBasket "github.com/borscht/backend/services/proto/basket"

	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/chat"
	chatDelivery "github.com/borscht/backend/internal/chat/delivery/http"
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
	apiGroup := data.e.Group("/api")
	userGroup := apiGroup.Group("/user", data.userMiddleware.Auth)
	auth := apiGroup.Group("", data.authMiddleware.Auth)
	userGroup.GET("", data.user.GetUserData)
	userGroup.PUT("", data.user.UpdateData)
	userGroup.PUT("/avatar", data.user.UploadAvatar)
	userGroup.POST("/address", data.user.UpdateMainAddress)
	userGroup.GET("/address", data.user.GetMainAddress)
	auth.GET("/auth", data.user.CheckAuth)
	auth.GET("/connect/ws", data.chat.GetKey)
	auth.GET("/chats", data.chat.GetAllChats)
	auth.GET("/chat/:id", data.chat.GetAllMessages)
	apiGroup.GET("/ws/:key", data.chat.Connect, data.wsMiddleware.WsAuth)

	restaurantGroup := apiGroup.Group("/restaurant", data.adminMiddleware.Auth)
	restaurantGroup.POST("/categories", data.restaurantAdmin.AddCategories)
	restaurantGroup.GET("/orders", data.order.GetRestaurantOrders)
	restaurantGroup.PUT("/order/status", data.order.SetNewStatus)
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

	apiGroup.POST("/signin", data.user.Login)
	apiGroup.POST("/signup", data.user.Create)
	apiGroup.POST("/restaurant/signup", data.restaurantAdmin.CreateRestaurant)
	apiGroup.POST("/restaurant/signin", data.restaurantAdmin.Login)
	userGroup.GET("/orders", data.order.GetUserOrders)
	userGroup.POST("/order/:basketID", data.order.Create)
	userGroup.POST("/order/review", data.order.CreateReview)
	userGroup.PUT("/basket/:restID", data.order.AddToBasket)
	userGroup.GET("/basket/:restID", data.order.GetBasket)
	userGroup.GET("/baskets", data.order.GetBaskets)
	userGroup.POST("/basket", data.order.AddBaskets)
	userGroup.DELETE("/basket/all", data.order.DeleteBaskets)
	userGroup.POST("/basket/delete/:basketID", data.order.DeleteBasket)
	apiGroup.GET("/logout", data.user.Logout)
	apiGroup.GET("/:id", data.restaurant.GetRestaurantPage)
	apiGroup.GET("/", data.restaurant.GetVendor)
	apiGroup.GET("/restaurants", data.restaurant.GetVendor)
	apiGroup.GET("/restaurant/:id/reviews", data.restaurant.GetReviews)
	apiGroup.GET("/restaurant/:id/recommendations", data.restaurant.GetRecommendations)
}

func initServer(e *echo.Echo) {
	e.Static("/static", config.ConfigStatic.Folder+"/static")
	e.Static("/default", config.ConfigStatic.Default)

	logger.InitLogger()
	e.Use(custMiddleware.LogMiddleware)
	e.Use(custMiddleware.CORS)
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "header:X-XSRF-TOKEN",
		CookiePath:  "/",
	}))

	e.Use(middleware.Secure())

	e.HTTPErrorHandler = custMiddleware.ErrorHandler
}

func main() {
	if config.ReadConfig() != nil {
		return
	}

	e := echo.New()
	initServer(e)
	grpcConnAuth, errr := grpc.Dial(
		config.AuthServiceAddress,
		grpc.WithInsecure(),
	)
	if errr != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer grpcConnAuth.Close()

	authClient := protoAuth.NewAuthClient(grpcConnAuth)
	authService := auth.NewService(authClient)

	grpcConnBasket, errr := grpc.Dial(
		config.BasketServiceAddress,
		grpc.WithInsecure(),
	)
	if errr != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer grpcConnBasket.Close()

	basketClient := protoBasket.NewBasketClient(grpcConnBasket)
	basketService := basket.NewService(basketClient)

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
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.ConfigDb.Host, config.ConfigDb.Port, config.ConfigDb.User,
		config.ConfigDb.Password, config.ConfigDb.NameDb)

	db, err := sql.Open(config.ConfigDb.NameSql, dsn)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(3)

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	websocketConnectionsUsers := websocketPool.NewConnectionPool()
	websocketConnectionsRestaurants := websocketPool.NewConnectionPool()
	orderNotificator := notifications.NewOrderNotificator(&websocketConnectionsUsers, &websocketConnectionsRestaurants)

	userRepo := userRepo.NewUserRepo(db)
	adminRestaurantRepo := restaurantAdminRepo.NewRestaurantRepo(db)
	adminDishRepo := restaurantAdminRepo.NewDishRepo(db)
	adminSectionRepo := restaurantAdminRepo.NewSectionRepo(db)
	restaurantRepo := restaurantRepo.NewRestaurantRepo(db)
	imageRepo := imageRepo.NewImageRepo()

	userUcase := userUcase.NewUserUsecase(userRepo, imageRepo)
	orderRepo := repository.NewOrderRepo(db)
	adminRestaurantUsecase := restaurantAdminUsecase.NewRestaurantUsecase(adminRestaurantRepo, imageRepo)
	adminDishUsecase := restaurantAdminUsecase.NewDishUsecase(adminDishRepo, adminSectionRepo, imageRepo)
	adminSectionUsecase := restaurantAdminUsecase.NewSectionUsecase(adminSectionRepo)
	restaurantUsecase := restaurantUsecase.NewRestaurantUsecase(restaurantRepo, adminRestaurantRepo)
	orderUsecase := usecase.NewOrderUsecase(orderRepo, adminRestaurantRepo)
	chatUsecase := chatUsecase.NewChatUsecase(chatService, authService, &websocketConnectionsUsers, &websocketConnectionsRestaurants)

	userHandler := userDelivery.NewUserHandler(userUcase, adminRestaurantUsecase, authService)
	adminRestaurantHandler := restaurantAdminDelivery.NewRestaurantHandler(adminRestaurantUsecase, authService)
	adminDishHandler := restaurantAdminDelivery.NewDishHandler(adminDishUsecase)
	adminSectionHandler := restaurantAdminDelivery.NewSectionHandler(adminSectionUsecase)
	restaurantHandler := restaurantDelivery.NewRestaurantHandler(restaurantUsecase)
	orderHandler := http.NewOrderHandler(orderUsecase, basketService, orderNotificator)
	chatHandler := chatDelivery.NewChatHandler(chatUsecase, authService)

	initUserMiddleware := custMiddleware.InitUserMiddleware(authService)
	initAdminMiddleware := custMiddleware.InitAdminMiddleware(authService, adminRestaurantUsecase)
	initAuthMiddleware := custMiddleware.InitAuthMiddleware(authService, userUcase, adminRestaurantUsecase)
	initWsMiddleware := custMiddleware.InitWsMiddleware(authService)

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

	e.Logger.Fatal(e.Start(":" + config.Config.Server.Port))
}
