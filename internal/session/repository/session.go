package repository

import (
	"context"

	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/models"
	sessionModel "github.com/borscht/backend/internal/session"
	errors "github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"

	"encoding/json"

	"github.com/gomodule/redigo/redis"
)

const headKey = "sessions:"

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
func (repo *sessionRepo) Check(ctx context.Context, sessionToCheck string) (models.SessionInfo, bool, error) {
	mkey := headKey + sessionToCheck
	data, err := redis.Bytes(repo.redisConn.Do("GET", mkey))
	if err != nil {
		return models.SessionInfo{}, false, err
	}
	sess := &models.SessionInfo{}
	err = json.Unmarshal(data, sess)
	if err != nil {
		return models.SessionInfo{}, false, err
	}
	return *sess, true, nil
}

// создание уникальной сессии
func (repo *sessionRepo) Create(ctx context.Context, sessionData models.SessionData) error {
	id := sessionID{sessionData.Session}
	dataSerialized, err := json.Marshal(models.SessionInfo{
		Id:   sessionData.Id,
		Role: sessionData.Role,
	})
	if err != nil {
		custError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, custError)
		return custError
	}

	mkey := headKey + id.ID

	result, err := redis.String(repo.redisConn.Do("SET", mkey, dataSerialized, "EX", config.LifetimeSecond))
	if err != nil {
		custError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, custError)
		return custError
	}
	if result != "OK" {
		custError := errors.FailServerError("result not OK")
		logger.RepoLevel().ErrorLog(ctx, custError)
		return custError
	}
	return nil
}

func (repo *sessionRepo) Delete(ctx context.Context, session string) error {
	mkey := headKey + session
	_, err := redis.Int(repo.redisConn.Do("DEL", mkey))
	if err != nil {
		custError := errors.FailServerError("redis error:" + err.Error())
		logger.RepoLevel().ErrorLog(ctx, custError)
		return custError
	}

	return nil
}
