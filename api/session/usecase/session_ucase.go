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
func (s *sessionUsecase) Check(session string, ctx *domain.CustomContext) (string, bool) {
	return s.sessionRepo.Check(session, ctx)
}

// создание уникальной сессии
func (s *sessionUsecase) Create(ctx *domain.CustomContext, uid string) (string, error) {
	session := ""

	for {
		session = uuid.New().String()

		_, isItExists := s.Check(session, ctx) // далее в цикле - проверка на уникальность
		if isItExists == false {               // не получили привязанного к сессии пользователя, следовательно, не существует
			break
		}
	}

	err := s.sessionRepo.Create(ctx, session, uid)
	if err != nil {
		return "", err
	}

	return session, nil
}
