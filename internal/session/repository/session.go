package repository

import (
	"context"

	"github.com/borscht/backend/config"
	sessionModel "github.com/borscht/backend/internal/session"
	errors "github.com/borscht/backend/utils"

	"encoding/json"

	"github.com/gomodule/redigo/redis"
)

const headKey = "sessions:"

type sessionData struct {
	Uid  int    `json:"uid"`
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
func (repo *sessionRepo) Check(ctx context.Context, sessionToCheck string) (int, bool, string) {
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
func (repo *sessionRepo) Create(ctx context.Context, session string, uid int, role string) error {
	id := sessionID{session}
	dataSerialized, err := json.Marshal(sessionData{
		Uid:  uid,
		Role: role,
	})
	if err != nil {
		return errors.FailServer(ctx, err.Error())
	}

	mkey := headKey + id.ID

	result, err := redis.String(repo.redisConn.Do("SET", mkey, dataSerialized, "EX", config.LifetimeSecond))
	if err != nil {
		return errors.FailServer(ctx, err.Error())
	}
	if result != "OK" {
		return errors.FailServer(ctx, "result not OK")
	}
	return nil
}

func (repo *sessionRepo) Delete(ctx context.Context, session string) error {
	mkey := headKey + session
	_, err := redis.Int(repo.redisConn.Do("DEL", mkey))
	if err != nil {
		return errors.FailServer(ctx, "redis error:"+err.Error())
	}

	return nil
}
