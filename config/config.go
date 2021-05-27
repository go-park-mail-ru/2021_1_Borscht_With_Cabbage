package config

import (
	"context"

	"github.com/borscht/backend/utils/logger"
	"github.com/spf13/viper"
)

// значения используемые в главном микросервисе
var (
	ConfigDb     Db
	ConfigStatic StaticProject
	Config       ConfigProject
	HostAddress  string
	Client       string

	ChatServiceAddress   string
	AuthServiceAddress   string
	BasketServiceAddress string
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

	var config ConfigProject
	viper.Unmarshal(&config)

	saveConfig(ctx, config)

	logger.UtilsLevel().InlineDebugLog(ctx, config)
	return nil
}

func saveConfig(ctx context.Context, config ConfigProject) {
	Config = config
	ConfigStatic = config.Static
	ConfigDb = config.Db
	HostAddress = config.Server.Host
	Client = config.Protocol + config.Client.Host
	if config.Client.Port != "" {
		Client += ":" + config.Client.Port
	}
	Client += "/"

	ChatServiceAddress = config.Microservices["chat"].GetFullHost()
	AuthServiceAddress = config.Microservices["auth"].GetFullHost()
	BasketServiceAddress = config.Microservices["basket"].GetFullHost()

	logger.UtilsLevel().InfoLog(ctx, logger.Fields{
		"Protocol":    Config.Protocol,
		"HostAddress": HostAddress,
		"Host":        Config.Server.Host,
		"ServerPort":  Config.Server.Port,
		"Repository":  ConfigStatic.Repository,
		"Client":      Client,

		"PostgresDB": ConfigDb.NameSql,
		"DBUser":     ConfigDb.User,
		"DBPass":     ConfigDb.Password,
		"DBName":     ConfigDb.NameDb,
		"DBHost":     ConfigDb.Host,
		"DBPort":     ConfigDb.Port,

		"Static":                 ConfigStatic.Folder,
		"DefaultStatic":          ConfigStatic.Default,
		"DefaultUserImage":       ConfigStatic.DefaultUserImage,
		"DefaultRestaurantImage": ConfigStatic.DefaultRestaurantImage,
		"DefaultDishImage":       ConfigStatic.DefaultDishImage,

		"ChatServiceAddress":   ChatServiceAddress,
		"AuthServiceAddress":   AuthServiceAddress,
		"BasketServiceAddress": BasketServiceAddress,
	})
}
