package config

import (
	"context"

	"github.com/borscht/backend/utils/logger"
	"github.com/spf13/viper"
)

// значения используемые в микросервисе auth
var (
	PostgresDB             string
	DBUser                 string
	DBPass                 string
	DBName                 string
	DBHost                 string
	DBPort                 string
	RedisHost              string
	Port                   string
	DefaultUserImage       string
	DefaultRestaurantImage string
	DefaultDishImage       string
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
	PostgresDB = config.Db.NameSql
	DBUser = config.Db.User
	DBPass = config.Db.Password
	DBName = config.Db.NameDb

	RedisHost = config.Redis.Host

	DefaultUserImage = config.Static.DefaultUserImage
	DefaultRestaurantImage = config.Static.DefaultRestaurantImage
	DefaultDishImage = config.Static.DefaultDishImage

	Port = config.Microservices["auth"].Port

	logger.UtilsLevel().InfoLog(ctx, logger.Fields{
		"PostgresDB":             PostgresDB,
		"DBUser":                 DBUser,
		"DBPass":                 DBPass,
		"DBName":                 DBName,
		"DBHost":                 DBHost,
		"DBPort":                 DBPort,
		"RedisHost":              RedisHost,
		"DefaultUserImage":       DefaultUserImage,
		"DefaultRestaurantImage": DefaultRestaurantImage,
		"DefaultDishImage":       DefaultDishImage,
		"Port":                   Port,
	})
}
