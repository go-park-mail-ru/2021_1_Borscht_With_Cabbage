package config

const (
	// Host          = "http://89.208.197.150"
	Host       = "http://127.0.0.1"
	Repository = Host + ":5000/"
	Client     = Host + ":3000/"

	DefaultAvatar  = Repository + "static/avatar/stas.jpg"
	SessionCookie  = "borscht_session"
	LifetimeSecond = 86400
	Static         = "static/avatar"

	PostgresDB = "postgres"
	DBUser     = "delivery" // todo ...
	DBPass     = "password"
	DBName     = "delivery_psql"
	CORStime   = 3600
	RoleAdmin  = "admin"
	RoleUser   = "user"

	RedisHost = "localhost:6379"
)
