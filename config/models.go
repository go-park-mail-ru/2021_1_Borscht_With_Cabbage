package config

type Db struct {
	NameSql  string
	User     string
	Password string
	NameDb   string
	Host     string
	Port     string
}

type ServerHost struct {
	Host string
	Port string
}

type ClientHost struct {
	Host string
	Port string
}

type Microservice struct {
	Port string
	Host string
}

func (s Microservice) GetFullHost() string {
	return s.Host + ":" + s.Port
}

type StaticProject struct {
	Folder                 string
	Default                string
	DefaultUserImage       string
	DefaultRestaurantImage string
	DefaultDishImage       string
}

type Config struct {
	Protocol      string
	Server        ServerHost
	Client        ClientHost
	Db            Db
	Static        StaticProject
	Microservices map[string]Microservice
}
