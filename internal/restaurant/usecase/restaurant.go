package usecase

import (
	"github.com/borscht/backend/internal/models"
	_restModel "github.com/borscht/backend/internal/restaurant"
)

type restaurantUsecase struct {
	restaurantRepo _restModel.RestaurantRepo
}

func NewRestaurantUsecase(repo _restModel.RestaurantRepo) _restModel.RestaurantUsecase {
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
