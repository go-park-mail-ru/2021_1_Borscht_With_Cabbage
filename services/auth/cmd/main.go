package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	"github.com/borscht/backend/config"
	"github.com/borscht/backend/services/auth/internal"
	repo "github.com/borscht/backend/services/auth/repository"
	protoAuth "github.com/borscht/backend/services/proto/auth"
	"github.com/borscht/backend/utils/logger"
	"github.com/gomodule/redigo/redis"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	logger.InitLogger()
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

	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalln("cant listen port", err)
	}

	server := grpc.NewServer()

	userAuthRepository := repo.NewUserAuthRepo(db)
	restaurantAuthRepo := repo.NewRestaurantAuthRepo(db)
	sessionRepository := repo.NewSessionRepo(redisConn)

	authService := internal.NewService(userAuthRepository, restaurantAuthRepo, sessionRepository)
	protoAuth.RegisterAuthServer(server, authService)

	log.Print("starting server at :8081")
	err = server.Serve(lis)
	if err != nil {
		log.Fatalln("Serve auth error: ", err)
	}
}
