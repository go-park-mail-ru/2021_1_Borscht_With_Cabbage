package usecase

import (
	"context"
	"path/filepath"

	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/restaurantAdmin"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
	"github.com/borscht/backend/utils/uniq"
)

const (
	HeadImage = "static/dish/"
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

func (a adminUsecase) UpdateDish(ctx context.Context, dish models.Dish) (*models.Dish, error) {
	//TODO: добавление изображения блюда в хранилище

	if dish.ID == 0 {
		requestError := errors.BadRequestError("No id at the dish")
		logger.UsecaseLevel().ErrorLog(ctx, requestError)
		return nil, requestError
	}

	// проверка прав на update
	ok := a.checkRightsForDish(ctx, dish.ID)
	if !ok {
		requestError := errors.BadRequestError("No rights to delete a dish")
		logger.UsecaseLevel().ErrorLog(ctx, requestError)
		return nil, requestError
	}

	err := a.adminRepository.UpdateDish(ctx, dish)
	if err != nil {
		return nil, err
	}
	return &dish, nil
}

func (a adminUsecase) DeleteDish(ctx context.Context, did int) error {
	//TODO: удаление изображения блюда из хранилища

	ok := a.checkRightsForDish(ctx, did)
	if !ok {
		requestError := errors.BadRequestError("No rights to delete a dish")
		logger.UsecaseLevel().ErrorLog(ctx, requestError)
		return requestError
	}

	return a.adminRepository.DeleteDish(ctx, did)
}

func (a adminUsecase) AddDish(ctx context.Context, dish models.Dish) (*models.Dish, error) {
	//TODO: добавление изображения блюда в хранилище

	id, err := a.adminRepository.AddDish(ctx, dish)
	if err != nil {
		return nil, err
	}
	responseDish := dish
	responseDish.ID = id
	responseDish.Image = config.DefaultAvatar // TODO: убрать
	return &responseDish, nil
}

func (a adminUsecase) UploadDishImage(ctx context.Context, image models.DishImage) (string, error) {
	if image.IdDish == 0 {
		requestError := errors.BadRequestError("No id at the dish")
		logger.UsecaseLevel().ErrorLog(ctx, requestError)
		return "", requestError
	}

	ok := a.checkRightsForDish(ctx, image.IdDish)
	if !ok {
		requestError := errors.BadRequestError("No rights to delete a dish")
		logger.UsecaseLevel().ErrorLog(ctx, requestError)
		return "", requestError
	}

	// парсим расширение
	expansion := filepath.Ext(image.Image.Filename)

	uid, localErr := uniq.GetUniqFilename(ctx, image.Image.Filename)
	if localErr != nil {
		return "", localErr
	}

	image.CustFilename = HeadImage + uid + expansion
	err := a.adminRepository.UploadDishImage(ctx, image)
	if err != nil {
		return "", err
	}

	image.CustFilename = config.Repository + HeadImage + uid + expansion
	err = a.adminRepository.UpdateDishImage(ctx, image)
	if err != nil {
		return "", err
	}

	return image.CustFilename, nil
}

func (a adminUsecase) Create(ctx context.Context, restaurant models.Restaurant) (*models.Restaurant, error) {
	restaurant.Avatar = config.DefaultAvatar

	id, err := a.adminRepository.Create(ctx, restaurant)
	if err != nil {
		return nil, err
	}
	restaurant.ID = id
	return &restaurant, nil
}

func (a adminUsecase) CheckRestaurantExists(ctx context.Context, user models.RestaurantAuth) (*models.Restaurant, error) {
	return a.adminRepository.CheckRestaurantExists(ctx, user)
}

func (a adminUsecase) GetByRid(ctx context.Context, rid int) (models.Restaurant, error) {
	return a.adminRepository.GetByRid(ctx, rid)
}

func (a adminUsecase) checkRightsForDish(ctx context.Context, idDish int) bool {
	dish, err := a.adminRepository.GetDish(ctx, idDish)
	if err != nil {
		return false
	}

	restaurant, ok := ctx.Value("Restaurant").(models.Restaurant)
	if !ok {
		failError := errors.FailServerError("failed to convert to models.Restaurant")
		logger.UsecaseLevel().ErrorLog(ctx, failError)
		return false
	}

	if restaurant.ID != dish.Restaurant {
		return false
	}

	return true
}
