package usecase

import (
	"context"

	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/restaurantAdmin"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
)

type restaurantUsecase struct {
	restaurantRepository restaurantAdmin.AdminRestaurantRepo
}

func NewRestaurantUsecase(adminRepo restaurantAdmin.AdminRestaurantRepo) restaurantAdmin.AdminRestaurantUsecase {

	return &restaurantUsecase{
		restaurantRepository: adminRepo,
	}
}

func (a restaurantUsecase) Update(ctx context.Context, restaurant models.RestaurantUpdate) (
	*models.RestaurantResponse, error) {

	restaurantAdmin, ok := ctx.Value("Restaurant").(models.Restaurant)
	if !ok {
		failError := errors.FailServerError("failed to convert to models.Restaurant")
		logger.UsecaseLevel().ErrorLog(ctx, failError)
		return nil, failError
	}

	restaurant.ID = restaurantAdmin.ID
	err := a.restaurantRepository.Update(ctx, restaurant)
	if err != nil {
		return nil, err
	}

	restaurantResponse := &models.RestaurantResponse{
		ID:           restaurant.ID,
		Title:        restaurant.Title,
		Description:  restaurant.Description,
		Rating:       restaurantAdmin.Rating,
		AvgCheck:     restaurantAdmin.AvgCheck,
		DeliveryCost: restaurant.DeliveryCost,
		Avatar:       restaurant.Avatar,
	}

	return restaurantResponse, nil
}

func (a restaurantUsecase) Create(ctx context.Context, restaurant models.Restaurant) (*models.Restaurant, error) {
	restaurant.Avatar = config.DefaultAvatar

	id, err := a.restaurantRepository.Create(ctx, restaurant)
	if err != nil {
		return nil, err
	}
	restaurant.ID = id
	return &restaurant, nil
}

func (a restaurantUsecase) CheckRestaurantExists(ctx context.Context, user models.RestaurantAuth) (*models.Restaurant, error) {
	return a.restaurantRepository.CheckRestaurantExists(ctx, user)
}

func (a restaurantUsecase) GetByRid(ctx context.Context, rid int) (models.Restaurant, error) {
	return a.restaurantRepository.GetByRid(ctx, rid)
}
