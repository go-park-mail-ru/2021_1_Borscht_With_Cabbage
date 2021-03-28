package usecase

import (
	sessionModel "github.com/borscht/backend/internal/session"
	"github.com/google/uuid"
)

type sessionUsecase struct {
	sessionRepo sessionModel.SessionRepo
}

func NewSessionUsecase(repo sessionModel.SessionRepo) sessionModel.SessionUsecase {
	return &sessionUsecase{
		sessionRepo: repo,
	}
}

// будет использоваться для проверки уникальности сессии при создании и для проверки авторизации на сайте в целом
func (s *sessionUsecase) Check(session string) (int32, bool, string) {
	return s.sessionRepo.Check(session)
}

// создание уникальной сессии
func (s *sessionUsecase) Create(uid int32, role string) (string, error) {
	session := ""
	for {
		session = uuid.New().String()

		_, isItExists, _ := s.sessionRepo.Check(session) // далее в цикле - проверка на уникальность
		if isItExists == false {                         // не получили привязанного к сессии пользователя, следовательно, не существует
			break
		}
	}

	err := s.sessionRepo.Create(session, uid, role)
	if err != nil {
		return "", err
	}

	return session, nil
}

func (s *sessionUsecase) Delete(session string) error {
	return s.sessionRepo.Delete(session)

}
