package config

type Db struct {
	NameSql  string
	User     string
	Password string
	NameDb   string
	Host     string
	Port     string
}

type Redis struct {
	Host string
}

type Static struct {
	DefaultUserImage       string
	DefaultRestaurantImage string
	DefaultDishImage       string
}

type Microservice struct {
	Port string
	Host string
}

type Config struct {
	Db            Db
	Redis         Redis
	Static        Static
	Microservices map[string]Microservice
}
