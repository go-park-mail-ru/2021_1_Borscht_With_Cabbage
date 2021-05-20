package config

import (
	"context"

	"github.com/borscht/backend/utils/logger"
	"github.com/spf13/viper"
)

// значения используемые в микросервисе basket
var (
	PostgresDB string
	DBUser     string
	DBPass     string
	DBName     string
	Port       string
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

	Port = config.Microservices["basket"].Port

	logger.UtilsLevel().InfoLog(ctx, logger.Fields{
		"PostgresDB": PostgresDB,
		"DBUser":     DBUser,
		"DBPass":     DBPass,
		"DBName":     DBName,
		"Port":       Port,
	})
}