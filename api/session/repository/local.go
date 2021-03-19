package repository

import (
	"backend/api/domain"
)

type sessionRepo struct {

}

func NewSessionRepo() domain.SessionRepo {
	return &sessionRepo{}
}


// будет использоваться для проверки уникальности сессии при создании и для проверки авторизации на сайте в целом
func (repo *sessionRepo) Check(sessionToCheck string, context *domain.CustomContext) (string, bool) {
	number, isItExists := (*context.Sessions)[sessionToCheck]
	if !isItExists {
		return "", false
	}
	return number, true
}

// создание уникальной сессии
func (repo *sessionRepo) Create(ctx *domain.CustomContext, session, uid string) error {

	(*ctx.Sessions)[session] = uid

	return nil
}

func (repo *sessionRepo) UpdateValue(ctx *domain.CustomContext, newValue, oldValue string) error {
	for j, numSession := range *ctx.Sessions {
		if numSession == oldValue {
			(*ctx.Sessions)[j] = newValue
		}
	}

	return nil
}
