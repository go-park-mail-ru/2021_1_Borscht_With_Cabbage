package usecase

import (
	"context"
	"github.com/borscht/backend/internal/models"
	restModel "github.com/borscht/backend/internal/restaurant"
	"github.com/borscht/backend/internal/restaurantAdmin"
	"github.com/borscht/backend/utils/logger"
)

type restaurantUsecase struct {
	restaurantRepo      restModel.RestaurantRepo
	restaurantAdminRepo restaurantAdmin.AdminRestaurantRepo
}

func NewRestaurantUsecase(repo restModel.RestaurantRepo,
	repoAdmin restaurantAdmin.AdminRestaurantRepo) restModel.RestaurantUsecase {

	return &restaurantUsecase{
		restaurantRepo:      repo,
		restaurantAdminRepo: repoAdmin,
	}
}

func (r *restaurantUsecase) GetVendor(ctx context.Context, params restModel.GetVendorParams) (
	restaurants []models.RestaurantInfo, err error) {

	restaurants, err = r.restaurantRepo.GetVendor(ctx, params)
	if err != nil {
		return nil, err
	}

	for i, value := range restaurants {
		address, err := r.restaurantAdminRepo.GetAddress(ctx, value.ID)
		if err != nil {
			return nil, err
		}
		logger.UsecaseLevel().InlineDebugLog(ctx, address)
		restaurants[i].Address = *address
	}
	logger.UsecaseLevel().InlineInfoLog(ctx, "finish GetVendor")
	return restaurants, nil
}

func (r *restaurantUsecase) GetById(ctx context.Context, coordinates models.Coordinates, id int) (
	restaurant *models.RestaurantWithDishes, err error) {
	restaurant, err = r.restaurantRepo.GetById(ctx, id, coordinates)
	if err != nil {
		return nil, err
	}

	address, err := r.restaurantAdminRepo.GetAddress(ctx, restaurant.ID)
	if err != nil {
		return nil, err
	}
	restaurant.Address = *address

	return restaurant, nil
}

func (r *restaurantUsecase) GetReviews(ctx context.Context, id int) ([]models.RestaurantReview, error) {
	return r.restaurantRepo.GetReviews(ctx, id)
}
