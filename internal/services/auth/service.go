package auth

import (
	"context"

	"github.com/borscht/backend/config"
	"github.com/borscht/backend/configProject"
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
	CreateKey(ctx context.Context, sessionInfo models.SessionInfo) (string, error)
	CheckKey(ctx context.Context, sessionToCheck string) (models.SessionInfo, bool, error)
}

type service struct {
	authService protoAuth.AuthClient
}

func NewService(authService protoAuth.AuthClient) ServiceAuth {
	return &service{
		authService: authService,
	}
}

func (s service) CheckKey(ctx context.Context, value string) (models.SessionInfo, bool, error) {
	session := protoAuth.SessionValue{
		Session: value,
	}

	sessionInfo, err := s.authService.CheckKey(ctx, &session)
	if err != nil {
		return models.SessionInfo{}, false, err
	}

	sessionOutput := models.SessionInfo{
		Id:   int(sessionInfo.Id),
		Role: sessionInfo.Role,
	}

	return sessionOutput, true, nil
}

func (s service) CreateKey(ctx context.Context, sessionInfo models.SessionInfo) (string, error) {
	session := protoAuth.SessionInfo{
		Id:   int32(sessionInfo.Id),
		Role: sessionInfo.Role,
	}

	sessionValue, err := s.authService.CreateKey(ctx, &session)
	if err != nil {
		return "", err
	}

	return sessionValue.Session, nil
}

// users
func (s service) Create(ctx context.Context, user models.User) (*models.SuccessUserResponse, error) {
	request := protoAuth.User{
		Name:     user.Name,
		Email:    user.Email,
		Phone:    user.Phone,
		Password: user.Password,
	}
	userResponse, err := s.authService.CreateUser(ctx, &request)
	if err != nil {
		return &models.SuccessUserResponse{}, err
	}
	logger.UsecaseLevel().InlineDebugLog(ctx, &user.HashPassword)

	user.Password = ""
	user.Uid = int(userResponse.UID)
	response := models.SuccessUserResponse{
		User: user,
		Role: configProject.RoleUser,
	}
	response.Avatar = config.ConfigStatic.DefaultUserImage

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

	Address := models.Address{
		Name:      userResult.AddressName,
		Longitude: userResult.Longitude,
		Latitude:  userResult.Latitude,
		Radius:    int(userResult.Radius),
	}
	userResponse := models.User{
		Uid:     int(userResult.UID),
		Name:    userResult.Name,
		Phone:   userResult.Phone,
		Email:   userResult.Email,
		Avatar:  userResult.Avatar,
		Address: Address,
	}
	response := models.SuccessUserResponse{
		User: userResponse,
		Role: configProject.RoleUser,
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

	Address := models.Address{
		Name:      user.AddressName,
		Longitude: user.Longitude,
		Latitude:  user.Latitude,
		Radius:    int(user.Radius),
	}
	UserResponse := models.User{
		Uid:     uid,
		Email:   user.Email,
		Phone:   user.Phone,
		Avatar:  user.Avatar,
		Name:    user.Name,
		Address: Address,
	}

	return &models.SuccessUserResponse{
		User: UserResponse,
		Role: configProject.RoleUser,
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
		ID:         int(restaurantResult.RID),
		Title:      restaurantResult.Title,
		Avatar:     config.ConfigStatic.DefaultRestaurantImage,
		AdminEmail: restaurant.AdminEmail,
		AdminPhone: restaurant.AdminPhone,
		Address:    restaurant.Address,
	}

	return &models.SuccessRestaurantResponse{
		RestaurantInfo: restaurantResponse,
		Role:           configProject.RoleAdmin,
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
		ID:           int(restaurantResult.RID),
		Title:        restaurantResult.Title,
		AdminPhone:   restaurantResult.Phone,
		AdminEmail:   restaurantResult.Email,
		DeliveryCost: int(restaurantResult.DeliveryCost),
		Description:  restaurantResult.Description,
		Rating:       float64(restaurantResult.Rating),
		AvgCheck:     int(restaurantResult.AvgCheck),
		Avatar:       restaurantResult.Avatar,
	}

	return &models.SuccessRestaurantResponse{
		RestaurantInfo: restaurantResponse,
		Role:           configProject.RoleAdmin,
	}, nil
}

func (s service) GetByRid(ctx context.Context, rid int) (*models.SuccessRestaurantResponse, error) {
	RID := protoAuth.RID{
		Rid: int32(rid),
	}

	restaurantResult, err := s.authService.GetByRid(ctx, &RID)
	if err != nil {
		return nil, err
	}

	Address := models.Address{
		Name:      restaurantResult.AddressName,
		Longitude: restaurantResult.Longitude,
		Latitude:  restaurantResult.Latitude,
		Radius:    int(restaurantResult.Radius),
	}
	restaurantResponse := models.RestaurantInfo{
		ID:         rid,
		Title:      restaurantResult.Title,
		AdminPhone: restaurantResult.Phone,
		AdminEmail: restaurantResult.Email,
		Avatar:     restaurantResult.Avatar,
		Address:    Address,
	}

	return &models.SuccessRestaurantResponse{
		RestaurantInfo: restaurantResponse,
		Role:           configProject.RoleAdmin,
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
