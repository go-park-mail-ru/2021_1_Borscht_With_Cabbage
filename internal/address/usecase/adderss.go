package usecase

import (
	"context"

	"github.com/borscht/backend/internal/address"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/restaurantAdmin"
	"github.com/borscht/backend/internal/user"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
)

type AddressUsecase struct {
	UserRepo       user.UserRepo
	RestaurantRepo restaurantAdmin.AdminRestaurantRepo
}

func NewAddressUsecase(user user.UserRepo,
	restaurant restaurantAdmin.AdminRestaurantRepo) address.AddressUsecase {

	return &AddressUsecase{
		UserRepo:       user,
		RestaurantRepo: restaurant,
	}
}

func (a AddressUsecase) UpdateMainAddress(ctx context.Context, address models.Address) error {
	restaurantInterface := ctx.Value("Restaurant")
	if restaurantInterface != nil {
		restaurant, ok := restaurantInterface.(models.RestaurantInfo)
		if !ok {
			failError := errors.FailServerError("failed to convert to models.Restaurant")
			logger.UsecaseLevel().ErrorLog(ctx, failError)
			return failError
		}

		return a.RestaurantRepo.UpdateMainAddress(ctx, restaurant.ID, address.Address)
	}

	user, ok := ctx.Value("User").(models.User)
	if !ok {
		failError := errors.FailServerError("failed to convert to models.User")
		logger.UsecaseLevel().ErrorLog(ctx, failError)
		return failError
	}

	return a.UserRepo.UpdateMainAddress(ctx, user.Uid, address.Address)
}
