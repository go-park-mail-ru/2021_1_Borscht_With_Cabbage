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

func (s *service) Create(ctx context.Context, user *protoAuth.User) (*protoAuth.SuccessUserResponse, error) {
	newUser := models.User{
		Email:    user.Email,
		Phone:    user.Phone,
		Name:     user.Name,
		Password: user.Password,
	}

	uid, err := s.userAuthRepo.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}

	response := protoAuth.SuccessUserResponse{
		Email:  newUser.Email,
		Phone:  newUser.Phone,
		Name:   newUser.Name,
		UID:    int32(uid),
		Avatar: "",
	}

	return &response, nil
}

func (s *service) GetByUid(ctx context.Context, uid *protoAuth.UID) (*protoAuth.SuccessUserResponse, error) {
	userResult, err := s.userAuthRepo.GetByUid(ctx, int(uid.Uid))
	if err != nil {
		return nil, err
	}

	response := protoAuth.SuccessUserResponse{
		Email:       userResult.Email,
		Phone:       userResult.Phone,
		Name:        userResult.Name,
		MainAddress: userResult.MainAddress,
		UID:         int32(uid.Uid),
		Avatar:      "",
	}

	return &response, nil
}

func (s *service) CheckUserExists(ctx context.Context, user *protoAuth.UserAuth) (*protoAuth.SuccessUserResponse, error) {
	userResult, err := s.userAuthRepo.GetByLogin(ctx, user.Login)
	if err != nil {
		return nil, err
	}

	response := protoAuth.SuccessUserResponse{
		Email:       userResult.Email,
		Phone:       userResult.Phone,
		Name:        userResult.Name,
		MainAddress: userResult.MainAddress,
		UID:         int32(userResult.Uid),
		Avatar:      "",
	}

	return &response, nil
}

func (s *service) CreateRestaurant(ctx context.Context, restaurant *protoAuth.User) (*protoAuth.SuccessRestaurantResponse, error) {
	panic("implement me")
}

func (a *service) CheckRestaurantExists(ctx context.Context, restaurantAuth *protoAuth.UserAuth) (*protoAuth.SuccessRestaurantResponse, error) {
	panic("implement me")
}

func (s *service) GetByRid(ctx context.Context, rid *protoAuth.RID) (*protoAuth.SuccessRestaurantResponse, error) {
	panic("implement me")
}

// будет использоваться для проверки уникальности сессии при создании и для проверки авторизации на сайте в целом
func (s *service) CheckSession(ctx context.Context, session *protoAuth.SessionValue) (*protoAuth.SessionInfo, error) {
	sessionInfo, exists, err := s.sessionRepo.Check(ctx, session.Session)
	if err != nil {
		return nil, err
	}

	sessionOutput := protoAuth.SessionInfo{
		Id:     int32(sessionInfo.Id),
		Role:   sessionInfo.Role,
		Exists: exists,
	}
	return &sessionOutput, nil
}

// создание уникальной сессии
func (s *service) CreateSession(ctx context.Context, sessionInfo *protoAuth.SessionInfo) (*protoAuth.SessionValue, error) {
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
		Id:      int(sessionInfo.Id),
		Role:    sessionInfo.Role,
	}
	err := s.sessionRepo.Create(ctx, sessionData)
	if err != nil {
		return nil, err
	}

	sessionOutput := protoAuth.SessionValue{
		Session: session,
	}

	return &sessionOutput, nil
}

func (s *service) DeleteSession(ctx context.Context, session *protoAuth.SessionValue) (*protoAuth.Error, error) {
	err := s.sessionRepo.Delete(ctx, session.Session)
	return nil, err
}