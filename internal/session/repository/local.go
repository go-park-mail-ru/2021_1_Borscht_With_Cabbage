package repository

import (
	"context"

	"github.com/borscht/backend/database/local"
	sessionModel "github.com/borscht/backend/internal/session"
)

type sessionRepo struct {
	db local.Database
}

func NewSessionRepo() sessionModel.SessionRepo {
	return &sessionRepo{
		db: local.GetInstance(),
	}
}

// будет использоваться для проверки уникальности сессии при создании и для проверки авторизации на сайте в целом
func (repo *sessionRepo) Check(ctx context.Context, sessionToCheck string) (string, bool) {
	number, isItExists := (*repo.db.GetModels().Sessions)[sessionToCheck]
	if !isItExists {
		return "", false
	}
	return number, true
}

// создание уникальной сессии
func (repo *sessionRepo) Create(ctx context.Context, session, uid string) error {

	(*repo.db.GetModels().Sessions)[session] = uid

	return nil
}

func (repo *sessionRepo) Update(ctx context.Context, newValue, oldValue string) error {
	for j, numSession := range *repo.db.GetModels().Sessions {
		if numSession == oldValue {
			(*repo.db.GetModels().Sessions)[j] = newValue
		}
	}

	return nil
}

func (repo *sessionRepo) Delete(ctx context.Context, session string) {
	_, ok := (*repo.db.GetModels().Sessions)[session]
	if ok {
		delete(*repo.db.GetModels().Sessions, session)
	}
}
