package auth

import (
	"context"
	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/models"
	protoAuth "github.com/borscht/backend/services/proto/auth"
	"github.com/borscht/backend/utils/logger"
	//"github.com/borscht/backend/utils/secure"
)

type ServiceAuth interface {
	Create(ctx context.Context, user models.User) (*models.SuccessUserResponse, error)
	CheckUserExists(ctx context.Context, user models.UserAuth) (*models.SuccessUserResponse, error)
	GetByUid(ctx context.Context, uid int) (*models.SuccessUserResponse, error)
	CreateRestaurant(ctx context.Context, restaurant models.RestaurantInfo) (*models.SuccessRestaurantResponse, error)
	CheckRestaurantExists(ctx context.Context, user models.RestaurantAuth) (*models.SuccessRestaurantResponse, error)
	GetByRid(ctx context.Context, rid int) (*models.SuccessRestaurantResponse, error)
	CheckSession(ctx context.Context, value string) (models.SessionInfo, bool, error)
	CreateSession(ctx context.Context, sessionInfo models.SessionInfo) (string, error)
	DeleteSession(ctx context.Context, session string) error
	//	GetUserData(ctx context.Context) (*models.SuccessUserResponse, error)
}

type service struct {
	authService protoAuth.AuthClient
}

func NewService(authService protoAuth.AuthClient) ServiceAuth {
	return &service{
		authService: authService,
	}
}

// users
func (s service) Create(ctx context.Context, user models.User) (*models.SuccessUserResponse, error) {
	request := protoAuth.User{
		Name:     user.Name,
		Email:    user.Email,
		Phone:    user.Phone,
		Password: user.Password,
	}
	userResponse, err := s.authService.Create(ctx, &request)
	if err != nil {
		return &models.SuccessUserResponse{}, err
	}
	logger.UsecaseLevel().InlineDebugLog(ctx, &user.HashPassword)

	user.Password = ""
	user.Uid = int(userResponse.UID)
	response := models.SuccessUserResponse{
		User: user,
		Role: config.RoleUser,
	}

	return &response, nil
}

func (s service) CheckUserExists(ctx context.Context, user models.UserAuth) (*models.SuccessUserResponse, error) {
	userProto := protoAuth.UserAuth{
		Login:    user.Login,
		Password: user.Password,
	}

	userResult, err := s.authService.CheckUserExists(ctx, &userProto)
	if err != nil {
		return &models.SuccessUserResponse{}, err
	}
	logger.UsecaseLevel().InlineDebugLog(ctx, &userResult)

	userResponse := models.User{
		Uid:    int(userResult.UID),
		Phone:  userResult.Phone,
		Email:  userResult.Email,
		Avatar: userResult.Avatar,
	}
	response := models.SuccessUserResponse{
		User: userResponse,
		Role: config.RoleUser,
	}

	return &response, nil
}

func (s service) GetByUid(ctx context.Context, uid int) (*models.SuccessUserResponse, error) {
	UID := protoAuth.UID{
		Uid: int32(uid),
	}
	user, err := s.authService.GetByUid(ctx, &UID)
	if err != nil {
		return nil, err
	}

	UserResponse := models.User{
		Uid:         uid,
		Email:       user.Email,
		Phone:       user.Phone,
		MainAddress: user.MainAddress,
		Avatar:      user.Avatar,
		Name:        user.Name,
	}

	return &models.SuccessUserResponse{
		User: UserResponse,
		Role: config.RoleUser,
	}, nil
}

// restaurants
func (s service) CreateRestaurant(ctx context.Context, restaurant models.RestaurantInfo) (*models.SuccessRestaurantResponse, error) {
	panic("implement me")
}

func (s service) CheckRestaurantExists(ctx context.Context, user models.RestaurantAuth) (*models.SuccessRestaurantResponse, error) {
	panic("implement me")
}

func (s service) GetByRid(ctx context.Context, rid int) (*models.SuccessRestaurantResponse, error) {
	panic("implement me")
}

// sessions
func (s service) CheckSession(ctx context.Context, value string) (models.SessionInfo, bool, error) {
	panic("implement me")
}

func (s service) CreateSession(ctx context.Context, sessionInfo models.SessionInfo) (string, error) {
	panic("implement me")
}

func (s service) DeleteSession(ctx context.Context, session string) error {
	panic("implement me")
}
