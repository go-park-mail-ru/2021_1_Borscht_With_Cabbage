package usecase

import (
	"context"

	"github.com/borscht/backend/internal/models"
	restModel "github.com/borscht/backend/internal/restaurant"
)

type restaurantUsecase struct {
	restaurantRepo restModel.RestaurantRepo
}

func NewRestaurantUsecase(repo restModel.RestaurantRepo) restModel.RestaurantUsecase {
	return &restaurantUsecase{
		restaurantRepo: repo,
	}
}

func (r *restaurantUsecase) GetVendor(ctx context.Context, limit, offset int) (*models.ArrayRestaurants, error) {
	restaurants, err := r.restaurantRepo.GetVendor(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return &models.ArrayRestaurants{Restaurants: restaurants}, nil
}

func (r *restaurantUsecase) GetById(ctx context.Context, id string) (*models.RestaurantWithDishes, error) {
	return r.restaurantRepo.GetById(ctx, id)
}
