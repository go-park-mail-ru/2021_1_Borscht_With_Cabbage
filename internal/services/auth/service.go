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
	restaurantToService := protoAuth.User{
		Name:     restaurant.Title,
		Email:    restaurant.AdminEmail,
		Phone:    restaurant.AdminPhone,
		Password: restaurant.AdminPassword,
	}

	restaurantResult, err := s.authService.CreateRestaurant(ctx, &restaurantToService)
	if err != nil {
		return nil, err
	}
	restaurantResponse := models.RestaurantInfo{
		Title:  restaurantResult.Title,
		Avatar: restaurant.Avatar,
	}

	return &models.SuccessRestaurantResponse{
		RestaurantInfo: restaurantResponse,
		Role:           config.RoleAdmin,
	}, nil
}

func (s service) CheckRestaurantExists(ctx context.Context, restaurant models.RestaurantAuth) (*models.SuccessRestaurantResponse, error) {
	authParametres := protoAuth.UserAuth{
		Login:    restaurant.Login,
		Password: restaurant.Password,
	}

	restaurantResult, err := s.authService.CheckRestaurantExists(ctx, &authParametres)
	if err != nil {
		return nil, err
	}
	restaurantResponse := models.RestaurantInfo{
		Title:  restaurantResult.Title,
		Avatar: restaurantResult.Avatar,
	}

	return &models.SuccessRestaurantResponse{
		RestaurantInfo: restaurantResponse,
		Role:           config.RoleAdmin,
	}, nil
}

func (s service) GetByRid(ctx context.Context, rid int) (*models.SuccessRestaurantResponse, error) {
	RID := protoAuth.RID{
		Uid: int32(rid),
	}

	restaurantResult, err := s.authService.GetByRid(ctx, &RID)
	if err != nil {
		return nil, err
	}
	restaurantResponse := models.RestaurantInfo{
		Title:  restaurantResult.Title,
		Avatar: restaurantResult.Avatar,
	}

	return &models.SuccessRestaurantResponse{
		RestaurantInfo: restaurantResponse,
		Role:           config.RoleAdmin,
	}, nil
}

// sessions
func (s service) CheckSession(ctx context.Context, value string) (models.SessionInfo, bool, error) {
	session := protoAuth.SessionValue{
		Session: value,
	}

	sessionInfo, err := s.authService.CheckSession(ctx, &session)
	if err != nil {
		return models.SessionInfo{}, false, err
	}

	sessionOutput := models.SessionInfo{
		Id:   int(sessionInfo.Id),
		Role: sessionInfo.Role,
	}

	return sessionOutput, true, nil
}

func (s service) CreateSession(ctx context.Context, sessionInfo models.SessionInfo) (string, error) {
	session := protoAuth.SessionInfo{
		Id:   int32(sessionInfo.Id),
		Role: sessionInfo.Role,
	}

	sessionValue, err := s.authService.CreateSession(ctx, &session)
	if err != nil {
		return "", err
	}

	return sessionValue.Session, nil
}

func (s service) DeleteSession(ctx context.Context, session string) error {
	sessionValue := protoAuth.SessionValue{
		Session: session,
	}

	_, err := s.authService.DeleteSession(ctx, &sessionValue)
	if err != nil {
		return nil
	}

	return nil
}
