package internal

import (
	"context"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/services/auth"
	protoAuth "github.com/borscht/backend/services/proto/auth"
	"github.com/google/uuid"
)

type service struct {
	userAuthRepo       auth.UserAuthRepo
	restaurantAuthRepo auth.RestaurantAuthRepo
	sessionRepo        auth.SessionRepo
}

func NewService(userAuthRepo auth.UserAuthRepo, restaurantAuthRepo auth.RestaurantAuthRepo, sessionRepo auth.SessionRepo) *service { // todo repositories here
	return &service{
		userAuthRepo:       userAuthRepo,
		restaurantAuthRepo: restaurantAuthRepo,
		sessionRepo:        sessionRepo,
	}
}

func (s *service) CreateUser(ctx context.Context, user *protoAuth.User) (protoAuth.SuccessUserResponse, error) {
	newUser := models.User{
		Email:    user.Email,
		Phone:    user.Phone,
		Name:     user.Name,
		Password: user.Password,
	}

	uid, err := s.userAuthRepo.Create(ctx, newUser)
	if err != nil {
		return protoAuth.SuccessUserResponse{}, err
	}

	response := protoAuth.SuccessUserResponse{
		Email:  newUser.Email,
		Phone:  newUser.Phone,
		Name:   newUser.Name,
		UID:    int32(uid),
		Avatar: "",
	}

	return response, nil
}

func (s *service) GetByUid(ctx context.Context, uid int) (protoAuth.SuccessUserResponse, error) {
	userResult, err := s.userAuthRepo.GetByUid(ctx, uid)
	if err != nil {
		return protoAuth.SuccessUserResponse{}, err
	}

	response := protoAuth.SuccessUserResponse{
		Email:       userResult.Email,
		Phone:       userResult.Phone,
		Name:        userResult.Name,
		MainAddress: userResult.MainAddress,
		UID:         int32(uid),
		Avatar:      "",
	}

	return response, nil
}

func (s *service) CheckUserExists(ctx context.Context, user protoAuth.UserAuth) (protoAuth.SuccessUserResponse, error) {
	userResult, err := s.userAuthRepo.GetByLogin(ctx, user.Login)
	if err != nil {
		return protoAuth.SuccessUserResponse{}, err
	}

	response := protoAuth.SuccessUserResponse{
		Email:       userResult.Email,
		Phone:       userResult.Phone,
		Name:        userResult.Name,
		MainAddress: userResult.MainAddress,
		UID:         int32(userResult.Uid),
		Avatar:      "",
	}

	return response, nil
}

func (s *service) CreateRestaurant(ctx context.Context, restaurant models.RestaurantInfo) (*models.SuccessRestaurantResponse, error) {

}

func (s *service) GetByRid(ctx context.Context, rid int) (*models.SuccessRestaurantResponse, error) {

}

func (a *service) CheckRestaurantExists(ctx context.Context, restaurantAuth models.RestaurantAuth) (*models.SuccessRestaurantResponse, error) {

}

// будет использоваться для проверки уникальности сессии при создании и для проверки авторизации на сайте в целом
func (s *service) CheckSession(ctx context.Context, session string) (models.SessionInfo, bool, error) {
	return s.sessionRepo.Check(ctx, session)
}

// создание уникальной сессии
func (s *service) CreateSession(ctx context.Context, sessionInfo models.SessionInfo) (string, error) {
	session := ""
	for {
		session = uuid.New().String()

		_, isItExists, _ := s.sessionRepo.Check(ctx, session) // далее в цикле - проверка на уникальность
		if isItExists == false {                              // не получили привязанного к сессии пользователя, следовательно, не существует
			break
		}
	}

	sessionData := models.SessionData{
		Session: session,
		Id:      sessionInfo.Id,
		Role:    sessionInfo.Role,
	}
	err := s.sessionRepo.Create(ctx, sessionData)
	if err != nil {
		return "", err
	}

	return session, nil
}

func (s *service) DeleteSession(ctx context.Context, session string) error {
	return s.sessionRepo.Delete(ctx, session)

}
