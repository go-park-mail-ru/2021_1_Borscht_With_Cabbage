package repository

import (
	"database/sql"
	"github.com/borscht/backend/config"
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

func (repo *sessionRepo) Check(sessionToCheck string) (int32, bool, string) {
	var id int32
	err := repo.DB.QueryRow("select uid from usersessions where session=$1", sessionToCheck).Scan(&id)
	if err != sql.ErrNoRows { // если она не уникальная
		return id, true, config.RoleUser
	}

	err = repo.DB.QueryRow("select rid from adminsessions where session=$1", sessionToCheck).Scan(&id)
	if err != sql.ErrNoRows {
		return id, true, config.RoleAdmin
	}

	return 0, false, ""
}

// создание уникальной сессии
func (repo *sessionRepo) Create(session string, id int32, role string) error {
	var err error
	switch role {
	case config.RoleUser:
		err = repo.DB.QueryRow("insert into usersessions (session, uid) values ($1, $2)", session, id).Scan()
	case config.RoleAdmin:
		err = repo.DB.QueryRow("insert into adminsessions (session, rid) values ($1, $2)", session, id).Scan()
	default:
		return errors.FailServer("role error")
	}
	if err != nil {
		return nil
	}

	return errors.FailServer("session saving failed")
}

func (repo *sessionRepo) Delete(session string) error {
	_, err := repo.DB.Exec("delete from usersessions where session=$1", session)
	_, errr := repo.DB.Exec("delete from adminsessions where session=$1", session)

	if err != nil || errr != nil {
		return errors.FailServer("session not found")
	}

	return nil
}
