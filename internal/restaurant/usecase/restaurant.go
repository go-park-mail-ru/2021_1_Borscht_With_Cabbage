package usecase

import (
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

func (r *restaurantUsecase) GetVendor(limit, offset int) ([]models.RestaurantResponse, error) {
	return r.restaurantRepo.GetVendor(limit, offset)
}

func (r *restaurantUsecase) GetById(id string) (models.Restaurant, bool, error) {
	return r.restaurantRepo.GetById(id)
}
