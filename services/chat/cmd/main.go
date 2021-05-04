package main

import (
	"database/sql"
	"fmt"

	"github.com/borscht/backend/config"

	"github.com/borscht/backend/services/chat/internal"
	chatServiceRepo "github.com/borscht/backend/services/chat/repository"
	protoChat "github.com/borscht/backend/services/proto/chat"

	"log"
	"net"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
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

	lis, err := net.Listen("tcp", ":8083")
	if err != nil {
		log.Fatalln("cant listen port", err)
	}

	server := grpc.NewServer()

	chatRepository := chatServiceRepo.NewChatRepository(db)
	chatService := internal.NewService(chatRepository)
	protoChat.RegisterChatServer(server, chatService)

	log.Print("starting server at :8083")
	err = server.Serve(lis)
	if err != nil {
		log.Fatalln("Serve auth error: ", err)
	}
}
