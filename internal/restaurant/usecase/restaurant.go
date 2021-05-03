package usecase

import (
	"context"
	"github.com/borscht/backend/utils/errors"
	"strconv"

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

func (r *restaurantUsecase) GetVendor(ctx context.Context, params restModel.GetVendorParams, longitude, latitude string) (
	restaurants []models.RestaurantInfo, err error) {

	longitudeNum, longitudeErr := strconv.ParseFloat(longitude, 64)
	latitudeNum, latitudeErr := strconv.ParseFloat(latitude, 64)
	if longitudeErr != nil || latitudeErr != nil {
		failError := errors.FailServerError("err while parsing longitude and latitude")
		logger.RepoLevel().ErrorLog(ctx, failError)
		return []models.RestaurantInfo{}, failError
	}
	params.Longitude = longitudeNum
	params.Latitude = latitudeNum

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

func (r *restaurantUsecase) GetById(ctx context.Context, id int) (
	restaurant *models.RestaurantWithDishes, err error) {
	restaurant, err = r.restaurantRepo.GetById(ctx, id)
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
