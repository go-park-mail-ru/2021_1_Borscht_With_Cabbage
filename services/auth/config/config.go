package config

import (
	"context"

	"github.com/borscht/backend/utils/logger"
	"github.com/spf13/viper"
)

// значения используемые в микросервисе auth
var (
	ConfigDb     Db
	ConfigStatic Static
	RedisHost    string
	Port         string
)

func ReadConfig() error {
	ctx := context.Background()

	viper.AddConfigPath("/etc/deliveryborscht/conf")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		logger.UtilsLevel().ErrorLog(ctx, err)
		return err
	}

	var config Config
	viper.Unmarshal(&config)

	saveConfig(ctx, config)

	logger.UtilsLevel().InlineDebugLog(ctx, config)
	return nil
}

func saveConfig(ctx context.Context, config Config) {
	ConfigStatic = config.Static
	ConfigDb = config.Db

	RedisHost = config.Redis.Host

	Port = config.Microservices["auth"].Port

	logger.UtilsLevel().InfoLog(ctx, logger.Fields{
		"PostgresDB": ConfigDb.NameSql,
		"DBUser":     ConfigDb.User,
		"DBPass":     ConfigDb.Password,
		"DBName":     ConfigDb.NameDb,
		"DBHost":     ConfigDb.Host,
		"DBPort":     ConfigDb.Port,

		"DefaultUserImage":       ConfigStatic.DefaultUserImage,
		"DefaultRestaurantImage": ConfigStatic.DefaultRestaurantImage,
		"DefaultDishImage":       ConfigStatic.DefaultDishImage,

		"RedisHost": RedisHost,

		"Port": Port,
	})
}
