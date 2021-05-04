package config

const (
	// HostAddress = "89.208.197.150"
	HostAddress = "127.0.0.1"
	Host        = "http://" + HostAddress
	ServerPort  = "5000"
	Repository  = Host + ":" + ServerPort + "/"
	Client      = Host + ":3000/"

	Static        = "static"
	DefaultStatic = "default"

	PostgresDB = "postgres"
	DBUser     = "labzunova" // todo ...
	DBPass     = "1111"
	DBName     = "postgres"

	RedisHost = "localhost:6379"

	ChatServiceAddress   = HostAddress + ":8083"
	AuthServiceAddress   = HostAddress + ":8081"
	BasketServiceAddress = HostAddress + ":8082"
)
