package config

import "time"

const (
	Port       = ":9000"
	PostgresDB = "postgres"
	ExpireTime = time.Hour * 24
)
