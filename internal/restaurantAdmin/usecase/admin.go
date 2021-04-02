package usecase

import (
	"context"

	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/restaurantAdmin"
	"github.com/borscht/backend/utils"
)

type adminUsecase struct {
	adminRepository restaurantAdmin.AdminRepo
}

func NewAdminUsecase(repo restaurantAdmin.AdminRepo) restaurantAdmin.AdminUsecase {
	return &adminUsecase{
		adminRepository: repo,
	}
}

func (a adminUsecase) DeleteDish(ctx context.Context, did int) error {
	dish, err := a.adminRepository.GetDish(ctx, did)
	if err != nil {
		return err
	}

	restaurant, ok := ctx.Value("Restaurant").(models.Restaurant)
	if !ok {
		return utils.FailServer(ctx, "failed to convert to models.Restaurant")
	}

	if restaurant.ID != dish.Restaurant {
		return utils.BadRequest(ctx, "No rights to delete a dish")
	}

	return a.adminRepository.DeleteDish(ctx, did)
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
