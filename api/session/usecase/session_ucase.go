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
func (s *sessionUsecase) Check(session string) (int32, bool) {
	return s.sessionRepo.Check(session)
}

// создание уникальной сессии
func (s *sessionUsecase) Create(uid string) (string, error) {
	session := ""

	for {
		session = uuid.New().String()

		_, isItExists := s.Check(session) // далее в цикле - проверка на уникальность
		if isItExists == false {          // не получили привязанного к сессии пользователя, следовательно, не существует
			break
		}
	}

	err := s.sessionRepo.Create(session, uid)
	if err != nil {
		return "", err
	}

	return session, nil
}

func (s *sessionUsecase) UpdateValue(newValue, oldValue string) error {
	return s.sessionRepo.UpdateValue(newValue, oldValue)
}

func (s *sessionUsecase) Delete(session string) error {
	return s.sessionRepo.Delete(session)
}
