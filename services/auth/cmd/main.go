package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	"github.com/borscht/backend/services/auth/config"
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
	if config.ReadConfig() != nil {
		return
	}

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

	lis, err := net.Listen("tcp", ":"+config.Port)
	if err != nil {
		log.Fatalln("cant listen port", err)
	}

	server := grpc.NewServer()

	userAuthRepository := repo.NewUserAuthRepo(db)
	restaurantAuthRepo := repo.NewRestaurantAuthRepo(db)
	sessionRepository := repo.NewSessionRepo(redisConn)

	authService := internal.NewService(userAuthRepository, restaurantAuthRepo, sessionRepository)
	protoAuth.RegisterAuthServer(server, authService)

	log.Print("starting server at :" + config.Port)
	err = server.Serve(lis)
	if err != nil {
		log.Fatalln("Serve auth error: ", err)
	}
}
