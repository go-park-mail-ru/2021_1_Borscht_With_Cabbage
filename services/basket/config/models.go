package config

type Db struct {
	NameSql  string
	User     string
	Password string
	NameDb   string
}

type Microservice struct {
	Port string
	Host string
}

type Config struct {
	Db            Db
	Microservices map[string]Microservice
}
