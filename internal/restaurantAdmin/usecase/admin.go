package usecase

import (
	"context"

	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/restaurantAdmin"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
)

type adminUsecase struct {
	adminRepository restaurantAdmin.AdminRepo
}

func NewAdminUsecase(repo restaurantAdmin.AdminRepo) restaurantAdmin.AdminUsecase {
	return &adminUsecase{
		adminRepository: repo,
	}
}

func (a adminUsecase) Update(ctx context.Context, restaurant models.RestaurantUpdate) (
	*models.RestaurantResponse, error) {
	// TODO: сохранение фотки

	restaurantAdmin, ok := ctx.Value("Restaurant").(models.Restaurant)
	if !ok {
		failError := errors.FailServerError("failed to convert to models.Restaurant")
		logger.UsecaseLevel().ErrorLog(ctx, failError)
		return nil, failError
	}

	restaurant.ID = restaurantAdmin.ID
	err := a.adminRepository.Update(ctx, restaurant)
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

func (a adminUsecase) GetAllDishes(ctx context.Context) ([]models.Dish, error) {
	restaurantAdmin, ok := ctx.Value("Restaurant").(models.Restaurant)
	if !ok {
		failError := errors.FailServerError("failed to convert to models.Restaurant")
		logger.UsecaseLevel().ErrorLog(ctx, failError)
		return nil, failError
	}

	return a.adminRepository.GetAllDishes(ctx, restaurantAdmin.ID)
}

func (a adminUsecase) UpdateDish(ctx context.Context, dish models.Dish) (*models.DishResponse, error) {
	//TODO: добавление изображения блюда в хранилище

	if dish.ID == 0 {
		requestError := errors.BadRequestError("No id at the dish")
		logger.UsecaseLevel().ErrorLog(ctx, requestError)
		return nil, requestError
	}

	// проверка прав на update
	oldDish, err := a.adminRepository.GetDish(ctx, dish.ID)
	if err != nil {
		return nil, err
	}
	restaurant, ok := ctx.Value("Restaurant").(models.Restaurant)
	if !ok {
		failError := errors.FailServerError("failed to convert to models.Restaurant")
		logger.UsecaseLevel().ErrorLog(ctx, failError)
		return nil, failError
	}
	if oldDish.Restaurant != restaurant.ID {
		requestError := errors.BadRequestError("No rights to update a dish")
		logger.UsecaseLevel().ErrorLog(ctx, requestError)
		return nil, requestError
	}

	err = a.adminRepository.UpdateDish(ctx, dish)
	if err != nil {
		return nil, err
	}
	responseDish := &models.DishResponse{
		ID:    dish.ID,
		Name:  dish.Name,
		Image: dish.Image,
	}
	return responseDish, nil
}

func (a adminUsecase) DeleteDish(ctx context.Context, did int) error {
	//TODO: удаление изображения блюда из хранилища

	dish, err := a.adminRepository.GetDish(ctx, did)
	if err != nil {
		return err
	}

	restaurant, ok := ctx.Value("Restaurant").(models.Restaurant)
	if !ok {
		failError := errors.FailServerError("failed to convert to models.Restaurant")
		logger.UsecaseLevel().ErrorLog(ctx, failError)
		return failError
	}

	if restaurant.ID != dish.Restaurant {
		requestError := errors.BadRequestError("No rights to delete a dish")
		logger.UsecaseLevel().ErrorLog(ctx, requestError)
		return requestError
	}

	return a.adminRepository.DeleteDish(ctx, did)
}

func (a adminUsecase) AddDish(ctx context.Context, dish models.Dish) (*models.DishResponse, error) {
	//TODO: добавление изображения блюда в хранилище

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
