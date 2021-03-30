package repository

import (
	"github.com/borscht/backend/config"
	sessionModel "github.com/borscht/backend/internal/session"
	errors "github.com/borscht/backend/utils"

	"encoding/json"

	"github.com/gomodule/redigo/redis"
)

const headKey = "sessions:"

type sessionData struct {
	Uid  int32  `json:"uid"`
	Role string `json:"role"`
}

type sessionID struct {
	ID string
}
type sessionRepo struct {
	redisConn redis.Conn
}

func NewSessionRepo(conn redis.Conn) sessionModel.SessionRepo {
	return &sessionRepo{
		redisConn: conn,
	}
}

// будет использоваться для проверки уникальности сессии при создании и для проверки авторизации на сайте в целом
<<<<<<< HEAD:internal/session/repository/local.go

func (repo *sessionRepo) Check(sessionToCheck string) (int32, bool, string) {
	mkey := headKey + sessionToCheck
	data, err := redis.Bytes(repo.redisConn.Do("GET", mkey))
	if err != nil {
		return 0, false, ""
	}
	sess := &sessionData{}
	err = json.Unmarshal(data, sess)
	if err != nil {
		return 0, false, ""
	}
	return sess.Uid, true, sess.Role
}

// создание уникальной сессии
<<<<<<< HEAD:internal/session/repository/local.go
func (repo *sessionRepo) Create(session string, uid int32, role string) error {
	id := sessionID{session}
	dataSerialized, err := json.Marshal(sessionData{
		Uid:  uid,
		Role: role,
	})

=======
func (repo *sessionRepo) Create(session string, id int, role string) error {
	var err error
	switch role {
	case config.RoleUser:
		err = repo.DB.QueryRow("insert into usersessions (session, uid) values ($1, $2)", session, id).Scan()
	case config.RoleAdmin:
		err = repo.DB.QueryRow("insert into adminsessions (session, rid) values ($1, $2)", session, id).Scan()
	default:
		return errors.FailServer("role error")
	}
>>>>>>> restaurantAuth:internal/session/repository/session.go
	if err != nil {
		return errors.FailServer(err.Error())
	}

	mkey := headKey + id.ID

	result, err := redis.String(repo.redisConn.Do("SET", mkey, dataSerialized, "EX", config.LifetimeSecond))
	if err != nil {
		return errors.FailServer(err.Error())
	}
	if result != "OK" {
		return errors.FailServer("result not OK")
	}
	return nil
}

func (repo *sessionRepo) Delete(session string) error {
	mkey := headKey + session
	_, err := redis.Int(repo.redisConn.Do("DEL", mkey))
	if err != nil {
		return errors.FailServer("redis error:" + err.Error())
	}

	return nil
}
