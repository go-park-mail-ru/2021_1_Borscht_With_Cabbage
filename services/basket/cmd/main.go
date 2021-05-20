package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	"github.com/borscht/backend/services/basket/config"
	"github.com/borscht/backend/services/basket/internal"
	basketServiceRepo "github.com/borscht/backend/services/basket/repository"
	protoBasket "github.com/borscht/backend/services/proto/basket"
	"github.com/borscht/backend/utils/logger"
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

	lis, err := net.Listen("tcp", ":"+config.Port)
	if err != nil {
		log.Fatalln("cant listen port", err)
	}

	server := grpc.NewServer()

	basketRepository := basketServiceRepo.NewBasketRepository(db)
	basketService := internal.NewService(basketRepository)
	protoBasket.RegisterBasketServer(server, basketService)

	log.Print("starting server at :" + config.Port)
	err = server.Serve(lis)
	if err != nil {
		log.Fatalln("Serve auth error: ", err)
	}
}
