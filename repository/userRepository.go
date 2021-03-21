package repository

import "database/sql"

type UserRepository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (repo *UserRepository) GetUserBySession(session string) {

}
