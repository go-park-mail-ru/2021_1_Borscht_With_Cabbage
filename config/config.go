package config

<<<<<<< HEAD
import "time"

const (
	Port       = ":9000"
	PostgresDB = "postgres"
	ExpireTime = time.Hour * 24
=======
const (
	Host          = "http://89.208.197.150"
	Repository    = Host + ":5000/"
	Client        = Host + ":3000/"
	DefaultAvatar = Repository + "static/avatar/stas.jpg"
	SessionCookie = "borscht_session"
	Static        = "static/avatar"
>>>>>>> 5db63d9c308ed18420262817abfc96fdeba2a06a
)
