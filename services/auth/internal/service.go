package internal

import (
	"context"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/secure"

	"github.com/borscht/backend/configProject"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/services/auth"
	"github.com/borscht/backend/services/auth/config"
	protoAuth "github.com/borscht/backend/services/proto/auth"
	"github.com/borscht/backend/utils/logger"
	"github.com/google/uuid"
)

const headSession = "sessions:"
const headKey = "key:"

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

func convertToSuccessUserResponse(user models.User, address models.Address, id int) protoAuth.SuccessUserResponse {
	return protoAuth.SuccessUserResponse{
		Email:       user.Email,
		Phone:       user.Phone,
		Name:        user.Name,
		Password:    user.Password,
		UID:         int32(id),
		Avatar:      user.Avatar,
		AddressName: address.Name,
		Latitude:    float32(address.Latitude),
		Longitude:   float32(address.Longitude),
		Radius:      int32(address.Radius),
		Role:        configProject.RoleUser,
	}
}

func convertToSuccessRestaurantResponse(restaurant models.RestaurantInfo, address models.Address) protoAuth.SuccessRestaurantResponse {
	return protoAuth.SuccessRestaurantResponse{
		RID:          int32(restaurant.ID),
		Title:        restaurant.Title,
		Email:        restaurant.AdminEmail,
		Phone:        restaurant.AdminPhone,
		DeliveryCost: int32(restaurant.DeliveryCost),
		AvgCheck:     int32(restaurant.AvgCheck),
		Description:  restaurant.Description,
		Rating:       float32(restaurant.Rating),
		Avatar:       restaurant.Avatar,
		Role:         configProject.RoleAdmin,
	}
}

func (s *service) CreateUser(ctx context.Context, user *protoAuth.User) (*protoAuth.SuccessUserResponse, error) {
	newUser := models.User{
		Email:    user.Email,
		Phone:    user.Phone,
		Name:     user.Name,
		Password: user.Password,
		Avatar:   config.ConfigStatic.DefaultUserImage,
	}

	newUser.HashPassword = secure.HashPassword(ctx, secure.GetSalt(), newUser.Password)
	uid, err := s.userAuthRepo.Create(ctx, newUser)
	if err != nil {
		return &protoAuth.SuccessUserResponse{}, err
	}

	response := convertToSuccessUserResponse(newUser, models.Address{}, uid)

	return &response, nil
}

func (s *service) GetByUid(ctx context.Context, uid *protoAuth.UID) (*protoAuth.SuccessUserResponse, error) {
	userResult, err := s.userAuthRepo.GetByUid(ctx, int(uid.Uid))
	if err != nil {
		return &protoAuth.SuccessUserResponse{}, err
	}

	address, err := s.userAuthRepo.GetAddress(ctx, int(uid.Uid))
	if err != nil {
		logger.UsecaseLevel().DebugLog(ctx, logger.Fields{"address error": err})
		return &protoAuth.SuccessUserResponse{}, err
	}
	if address != nil {
		userResult.Address = *address
	}

	response := convertToSuccessUserResponse(userResult, *address, int(uid.Uid))

	return &response, nil
}

func (s *service) CheckUserExists(ctx context.Context, user *protoAuth.UserAuth) (*protoAuth.SuccessUserResponse, error) {
	logger.UsecaseLevel().InlineDebugLog(ctx, user)
	userResult, err := s.userAuthRepo.GetByLogin(ctx, user.Login)
	if err != nil {
		return &protoAuth.SuccessUserResponse{}, err
	}
	logger.UsecaseLevel().InlineDebugLog(ctx, userResult)

	if !secure.CheckPassword(ctx, userResult.HashPassword, user.Password) {
		logger.UsecaseLevel().InlineInfoLog(ctx, "Password faild")
		return nil, errors.AuthorizationError("user not found")
	}

	address, err := s.userAuthRepo.GetAddress(ctx, userResult.Uid)
	if err != nil {
		logger.UsecaseLevel().DebugLog(ctx, logger.Fields{"address error": err})
		return &protoAuth.SuccessUserResponse{}, err
	}

	response := convertToSuccessUserResponse(*userResult, *address, userResult.Uid)

	return &response, nil
}

func (s *service) CreateRestaurant(ctx context.Context, restaurant *protoAuth.User) (*protoAuth.SuccessRestaurantResponse, error) {
	newRestaurant := models.RestaurantInfo{
		Title:         restaurant.Name,
		AdminEmail:    restaurant.Email,
		AdminPhone:    restaurant.Phone,
		AdminPassword: restaurant.Password,
		Avatar:        config.ConfigStatic.DefaultRestaurantImage,
	}

	newRestaurant.AdminHashPassword = secure.HashPassword(ctx, secure.GetSalt(), newRestaurant.AdminPassword)
	rid, err := s.restaurantAuthRepo.CreateRestaurant(ctx, newRestaurant)
	if err != nil {
		return &protoAuth.SuccessRestaurantResponse{}, err
	}

	response := protoAuth.SuccessRestaurantResponse{
		Title: restaurant.Name,
		Email: restaurant.Email,
		Phone: restaurant.Phone,
		Role:  configProject.RoleAdmin,
		RID:   int32(rid),
	}

	return &response, nil
}

func (s *service) CheckRestaurantExists(ctx context.Context, restaurantAuth *protoAuth.UserAuth) (*protoAuth.SuccessRestaurantResponse, error) {
	logger.UsecaseLevel().InlineDebugLog(ctx, restaurantAuth)
	restaurant, err := s.restaurantAuthRepo.GetByLogin(ctx, restaurantAuth.Login)
	if err != nil {
		return &protoAuth.SuccessRestaurantResponse{}, err
	}

	logger.UsecaseLevel().InlineDebugLog(ctx, restaurant)

	if !secure.CheckPassword(ctx, restaurant.AdminHashPassword, restaurantAuth.Password) {
		logger.UsecaseLevel().InlineInfoLog(ctx, "Password faild")
		return nil, errors.AuthorizationError("restaurant not found")
	}

	response := convertToSuccessRestaurantResponse(*restaurant, models.Address{})

	return &response, nil
}

func (s *service) GetByRid(ctx context.Context, rid *protoAuth.RID) (*protoAuth.SuccessRestaurantResponse, error) {
	restaurant, err := s.restaurantAuthRepo.GetByRid(ctx, int(rid.Rid))
	if err != nil {
		return &protoAuth.SuccessRestaurantResponse{}, err
	}

	address, err := s.restaurantAuthRepo.GetAddress(ctx, int(rid.Rid))
	if err != nil {
		logger.UsecaseLevel().DebugLog(ctx, logger.Fields{"address error": err})
		return &protoAuth.SuccessRestaurantResponse{}, err
	}

	response := convertToSuccessRestaurantResponse(*restaurant, *address)

	return &response, nil
}

func (s *service) CheckKey(ctx context.Context, session *protoAuth.SessionValue) (*protoAuth.SessionInfo, error) {
	sessionInfo, exists, err := s.sessionRepo.Check(ctx, headKey+session.Session)
	if err != nil {
		return &protoAuth.SessionInfo{}, err
	}

	sessionOutput := protoAuth.SessionInfo{
		Id:     int32(sessionInfo.Id),
		Role:   sessionInfo.Role,
		Exists: exists,
	}
	return &sessionOutput, nil
}

// создание уникальной сессии
func (s *service) CreateKey(ctx context.Context, sessionInfo *protoAuth.SessionInfo) (*protoAuth.SessionValue, error) {
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
		Id:              int(sessionInfo.Id),
		Role:            sessionInfo.Role,
		LifeTimeSeconds: 60,
	}
	err := s.sessionRepo.Create(ctx, sessionData)
	if err != nil {
		return &protoAuth.SessionValue{}, err
	}

	sessionOutput := protoAuth.SessionValue{
		Session: session,
	}

	return &sessionOutput, nil
}

// будет использоваться для проверки уникальности сессии при создании и для проверки авторизации на сайте в целом
func (s *service) CheckSession(ctx context.Context, session *protoAuth.SessionValue) (*protoAuth.SessionInfo, error) {
	sessionInfo, exists, err := s.sessionRepo.Check(ctx, headSession+session.Session)
	if err != nil {
		return &protoAuth.SessionInfo{}, err
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

		_, isItExists, _ := s.sessionRepo.Check(ctx, headSession+session) // далее в цикле - проверка на уникальность
		if !isItExists {                                                  // не получили привязанного к сессии пользователя, следовательно, не существует
			break
		}
	}

	sessionData := models.SessionData{
		Session:         headSession + session,
		Id:              int(sessionInfo.Id),
		Role:            sessionInfo.Role,
		LifeTimeSeconds: configProject.LifetimeSecond,
	}
	err := s.sessionRepo.Create(ctx, sessionData)
	if err != nil {
		return &protoAuth.SessionValue{}, err
	}

	sessionOutput := protoAuth.SessionValue{
		Session: session,
	}

	return &sessionOutput, nil
}

func (s *service) DeleteSession(ctx context.Context, session *protoAuth.SessionValue) (*protoAuth.Error, error) {
	err := s.sessionRepo.Delete(ctx, headSession+session.Session)
	return &protoAuth.Error{}, err
}
