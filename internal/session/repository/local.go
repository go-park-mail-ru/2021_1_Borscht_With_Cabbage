package repository

import (
	"context"
	"database/sql"
	"github.com/borscht/backend/database/local"
	_sessionModel "github.com/borscht/backend/internal/session"
	errors "github.com/borscht/backend/utils"
)

type sessionRepo struct {
	db local.Database
	DB *sql.DB
}

func NewSessionRepo() _sessionModel.SessionRepo {
	return &sessionRepo{
		db: local.GetInstance(),
	}
}


// будет использоваться для проверки уникальности сессии при создании и для проверки авторизации на сайте в целом
func (repo *sessionRepo) Check(sessionToCheck string) (int32, bool) {
	//number, isItExists := (*repo.db.GetModels().Sessions)[sessionToCheck]
	//if !isItExists {
	//	return "", false
	//}
	//return number, true
	var uid int32
	err := repo.DB.QueryRow("select uid from sessions where session=&1", sessionToCheck).Scan(&uid)
	if err != nil { // если она уникальная
		return uid, true
	}
	return 0, false
}

// создание уникальной сессии
func (repo *sessionRepo) Create(session, uid int32) error {
	//
	//(*repo.db.GetModels().Sessions)[session] = uid
	//
	//return nil

	err := repo.DB.QueryRow("insert into sessions (session, uid) values ($1, $2)", session, uid)
	if err != nil {
		return nil
	}
	return errors.FailServer("session saving failed")
}

func (repo *sessionRepo) Delete(session string) error {
	//_, ok := (*repo.db.GetModels().Sessions)[session]
	//if ok {
	//	delete(*repo.db.GetModels().Sessions, session)
	//}
	err := repo.DB.QueryRow("delete from sessions where session=$1", session)
	if err != nil {
		return errors.FailServer("session not found")
	}
	return nil
}
