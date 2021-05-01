package main

import (
	"database/sql"
	"fmt"
	"github.com/borscht/backend/config"
	"github.com/borscht/backend/services/basket/internal"
	basketServiceRepo "github.com/borscht/backend/services/basket/repository"
	protoBasket "github.com/borscht/backend/services/proto/basket"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
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

	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatalln("cant listen port", err)
	}

	server := grpc.NewServer()

	basketRepository := basketServiceRepo.NewBasketRepository(db)
	basketService := internal.NewService(basketRepository)
	protoBasket.RegisterBasketServer(server, basketService)

	log.Print("starting server at :8082")
	err = server.Serve(lis)
	if err != nil {
		log.Fatalln("Serve auth error: ", err)
	}
}
