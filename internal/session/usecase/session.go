package usecase

import (
	"context"

	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/models"
	sessionModel "github.com/borscht/backend/internal/session"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
	"github.com/google/uuid"
)

const headSession = "sessions:"
const headKey = "key:"

type sessionUsecase struct {
	sessionRepo sessionModel.SessionRepo
}

func NewSessionUsecase(repo sessionModel.SessionRepo) sessionModel.SessionUsecase {
	return &sessionUsecase{
		sessionRepo: repo,
	}
}

// будет использоваться для проверки уникальности сессии при создании и для проверки авторизации на сайте в целом
func (s *sessionUsecase) Check(ctx context.Context, session string) (models.SessionInfo, bool, error) {
	return s.sessionRepo.Check(ctx, headSession+session)
}

// создание уникальной сессии
func (s *sessionUsecase) Create(ctx context.Context, sessionInfo models.SessionInfo) (string, error) {
	session := ""
	for {
		session = uuid.New().String()

		_, isItExists, _ := s.sessionRepo.Check(ctx, headSession+session) // далее в цикле - проверка на уникальность
		if !isItExists {                                                  // не получили привязанного к сессии пользователя, следовательно, не существует
			break
		}
	}

	sessionData := models.SessionData{
		Session:         headSession + session,
		Id:              sessionInfo.Id,
		Role:            sessionInfo.Role,
		LifeTimeSeconds: config.LifetimeSecond,
	}
	err := s.sessionRepo.Create(ctx, sessionData)
	if err != nil {
		return "", err
	}

	return session, nil
}

func (s *sessionUsecase) Delete(ctx context.Context, session string) error {
	return s.sessionRepo.Delete(ctx, headSession+session)

}

func (s *sessionUsecase) CheckKey(ctx context.Context, session string) (models.SessionInfo, bool, error) {
	return s.sessionRepo.Check(ctx, headKey+session)
}

func getSessionInfo(ctx context.Context) (*models.SessionInfo, error) {
	userInterface := ctx.Value("User")
	if userInterface != nil {
		user, ok := userInterface.(models.User)
		if !ok {
			failError := errors.FailServerError("failed to convert to models.User")
			logger.UsecaseLevel().ErrorLog(ctx, failError)
			return nil, failError
		}

		return &models.SessionInfo{
			Id:   user.Uid,
			Role: config.RoleUser,
		}, nil
	}

	restaurantInterface := ctx.Value("Restaurant")
	if restaurantInterface != nil {
		restaurant, ok := restaurantInterface.(models.RestaurantInfo)
		if !ok {
			failError := errors.FailServerError("failed to convert to models.Restaurant")
			logger.UsecaseLevel().ErrorLog(ctx, failError)
			return nil, failError
		}

		return &models.SessionInfo{
			Id:   restaurant.ID,
			Role: config.RoleAdmin,
		}, nil
	}

	return nil, errors.BadRequestError("not authorization")
}

func (s *sessionUsecase) CreateKey(ctx context.Context) (string, error) {
	sessionInfo, err := getSessionInfo(ctx)
	if err != nil {
		return "", err
	}
	session := ""
	for {
		session = uuid.New().String()

		_, isItExists, _ := s.sessionRepo.Check(ctx, headKey+session) // далее в цикле - проверка на уникальность
		if !isItExists {                                              // не получили привязанного к сессии пользователя, следовательно, не существует
			break
		}
	}

	sessionData := models.SessionData{
		Session:         headKey + session,
		Id:              sessionInfo.Id,
		Role:            sessionInfo.Role,
		LifeTimeSeconds: 60,
	}
	err = s.sessionRepo.Create(ctx, sessionData)
	if err != nil {
		return "", err
	}

	return session, nil
}
