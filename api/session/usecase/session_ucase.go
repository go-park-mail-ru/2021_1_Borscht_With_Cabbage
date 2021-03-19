package usecase

import (
	"backend/api/domain"
	"github.com/google/uuid"
)

type sessionUsecase struct {
	sessionRepo domain.SessionRepo
}

func NewSessionUsecase(repo domain.SessionRepo) domain.SessionUsecase {
	return &sessionUsecase{
		sessionRepo: repo,
	}
}


// будет использоваться для проверки уникальности сессии при создании и для проверки авторизации на сайте в целом
func (u *sessionUsecase) Check(sessionToCheck string, context *domain.CustomContext) (string, bool) {
	number, isItExists := (*context.Sessions)[sessionToCheck]
	if !isItExists {
		return "", false
	}
	return number, true
}

// создание уникальной сессии
func (u *sessionUsecase) Create(ctx *domain.CustomContext, uid string) (string, error) {
	session := ""

	for {
		session = uuid.New().String()

		_, isItExists := u.Check(session, ctx) // далее в цикле - проверка на уникальность
		if isItExists == false {               // не получили привязанного к сессии пользователя, следовательно, не существует
			break
		}
	}

	err := u.sessionRepo.Create(ctx, session, uid)
	if err != nil {
		return "", err
	}

	return session, nil
}
