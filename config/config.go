package config

import (
	"context"

	"github.com/borscht/backend/utils/logger"
	"github.com/spf13/viper"
)

// значения используемые в главном микросервисе
var (
	Protocol    string
	HostAddress string
	Host        string
	ServerPort  string
	Repository  string
	Client      string

	Static        string
	DefaultStatic string

	PostgresDB string
	DBUser     string
	DBPass     string
	DBName     string
	DBHost     string
	DBPort     string

	ChatServiceAddress   string
	AuthServiceAddress   string
	BasketServiceAddress string

	DefaultUserImage       string
	DefaultRestaurantImage string
	DefaultDishImage       string
)

func ReadConfig() error {
	ctx := context.Background()

	viper.SetConfigType("yml")
	viper.AddConfigPath("/etc/deliveryborscht/conf")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		logger.UtilsLevel().ErrorLog(ctx, err)
		return err
	}

	var config Config
	viper.Unmarshal(&config)

	saveConfig(ctx, config)

	logger.UtilsLevel().InlineInfoLog(ctx, config)
	return nil
}

func saveConfig(ctx context.Context, config Config) {
	Protocol = config.Protocol
	HostAddress = config.Server.Host
	Host = Protocol + HostAddress
	ServerPort = config.Server.Port
	Repository = config.Static.Repository
	Client = config.Client.Host + ":" + config.Client.Port + "/"

	Static = config.Static.Folder
	DefaultStatic = config.Static.Default

	PostgresDB = config.Db.NameSql
	DBUser = config.Db.User
	DBPass = config.Db.Password
	DBName = config.Db.NameDb
	DBHost = config.Db.Host
	DBPort = config.Db.Port

	ChatServiceAddress = config.Microservices["chat"].GetFullHost()
	AuthServiceAddress = config.Microservices["auth"].GetFullHost()
	BasketServiceAddress = config.Microservices["basket"].GetFullHost()

	DefaultUserImage = config.Static.DefaultUserImage
	DefaultRestaurantImage = config.Static.DefaultRestaurantImage
	DefaultDishImage = config.Static.DefaultDishImage

	logger.UtilsLevel().InfoLog(ctx, logger.Fields{
		"Protocol":    Protocol,
		"HostAddress": HostAddress,
		"Host":        Host,
		"ServerPort":  ServerPort,
		"Repository":  Repository,
		"Client":      Client,

		"PostgresDB": PostgresDB,
		"DBUser":     DBUser,
		"DBPass":     DBPass,
		"DBName":     DBName,
		"DBHost":     DBHost,
		"DBPort":     DBPort,

		"Static":                 Static,
		"DefaultStatic":          DefaultStatic,
		"DefaultUserImage":       DefaultUserImage,
		"DefaultRestaurantImage": DefaultRestaurantImage,
		"DefaultDishImage":       DefaultDishImage,

		"ChatServiceAddress":   ChatServiceAddress,
		"AuthServiceAddress":   AuthServiceAddress,
		"BasketServiceAddress": BasketServiceAddress,
	})
}
