package usecase

import (
	"context"

	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/restaurantAdmin"
)

type adminUsecase struct {
	adminRepository restaurantAdmin.AdminRepo
}

func NewAdminUsecase(repo restaurantAdmin.AdminRepo) restaurantAdmin.AdminUsecase {
	return &adminUsecase{
		adminRepository: repo,
	}
}

func (a adminUsecase) AddDish(ctx context.Context, dish models.Dish) (*models.DishResponse, error) {
	id, err := a.adminRepository.AddDish(ctx, dish)
	if err != nil {
		return nil, err
	}
	responseDish := &models.DishResponse{
		ID:    id,
		Name:  dish.Name,
		Image: dish.Image,
	}
	return responseDish, nil
}

func (a adminUsecase) Create(ctx context.Context, user models.Restaurant) (int, error) {
	return a.adminRepository.Create(ctx, user)
}

func (a adminUsecase) CheckRestaurantExists(ctx context.Context, user models.RestaurantAuth) (models.Restaurant, error) {
	return a.adminRepository.CheckRestaurantExists(ctx, user)
}

func (a adminUsecase) GetByRid(ctx context.Context, rid int) (models.Restaurant, error) {
	return a.adminRepository.GetByRid(ctx, rid)
}
