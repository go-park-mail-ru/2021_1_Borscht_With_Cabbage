package repository

import (
	"database/sql"
	sessionModel "github.com/borscht/backend/internal/session"
	errors "github.com/borscht/backend/utils"
)

type sessionRepo struct {
	DB *sql.DB
}

func NewSessionRepo(db *sql.DB) sessionModel.SessionRepo {
	return &sessionRepo{
		DB: db,
	}
}

// будет использоваться для проверки уникальности сессии при создании и для проверки авторизации на сайте в целом
func (repo *sessionRepo) Check(sessionToCheck string) (int32, bool) {
	var uid int32
	err := repo.DB.QueryRow("select uid from sessions where session=$1", sessionToCheck).Scan(&uid)

	if err != sql.ErrNoRows { // если она не уникальна
		return uid, true
	}

	return 0, false
}

// создание уникальной сессии
func (repo *sessionRepo) Create(session string, uid int32, role string) error {
	err := repo.DB.QueryRow("insert into sessions (session, uid, role) values ($1, $2, $3)", session, uid, role)
	if err != nil {
		return nil
	}

	return errors.FailServer("session saving failed")
}

func (repo *sessionRepo) Delete(session string) error {
	_, err := repo.DB.Exec("delete from sessions where session=$1", session)
	if err != nil {
		return errors.FailServer("session not found")
	}

	return nil
}
