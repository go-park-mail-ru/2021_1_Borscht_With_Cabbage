package usecase

import "backend/api/domain"

type restaurantUsecase struct {
	restaurantRepo domain.RestaurantRepo
}

func NewRestaurantUsecase(repo domain.RestaurantRepo) domain.RestaurantUsecase {
	return &restaurantUsecase{
		restaurantRepo: repo,
	}
}

func (r *restaurantUsecase) GetSlice(ctx *domain.CustomContext, limit, offset int) ([]domain.RestaurantResponse, error) {
	return r.restaurantRepo.GetSlice(ctx, limit, offset)
}

func (r *restaurantUsecase) GetById(ctx *domain.CustomContext, id string) (domain.Restaurant, bool) {
	return r.restaurantRepo.GetById(ctx, id)
}
