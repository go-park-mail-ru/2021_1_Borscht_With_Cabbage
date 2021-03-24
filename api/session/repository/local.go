package repository

import (
	"backend/api/domain"
	errors "backend/utils"
	"database/sql"
)

type sessionRepo struct {
	DB *sql.DB
}

func NewSessionRepo() domain.SessionRepo {
	return &sessionRepo{}
}

// будет использоваться для проверки уникальности сессии при создании и для проверки авторизации на сайте в целом
func (repo *sessionRepo) Check(sessionToCheck string) (int32, bool) {
	var uid int32
	err := repo.DB.QueryRow("select uid from sessions where session=&1", sessionToCheck).Scan(&uid)
	if err != nil { // если она уникальная
		return uid, true
	}
	return 0, false
}

// создание уникальной сессии
func (repo *sessionRepo) Create(session, uid string) error {
	err := repo.DB.QueryRow("insert into sessions (session, uid) values ($1, $2)", session, uid)
	if err != nil {
		return nil
	}
	return errors.FailServer("session saving failed")
}

func (repo *sessionRepo) UpdateValue(newValue, oldValue string) error {
	err := repo.DB.QueryRow("update sessions set session = $1 where session = $2", newValue, oldValue)
	if err != nil {
		return nil
	}
	return errors.FailServer("session not found")
}

func (repo *sessionRepo) Delete(session string) error {
	err := repo.DB.QueryRow("delete from sessions where session=$1", session)
	if err != nil {
		return errors.FailServer("session not found")
	}
	return nil
}
