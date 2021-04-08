package usecase

import (
	"context"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/image"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/restaurantAdmin"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
	"github.com/borscht/backend/utils/uniq"
)

const (
	HeadImage = "static/dish/"
)

type dishUsecase struct {
	dishRepository  restaurantAdmin.AdminDishRepo
	imageRepository image.ImageRepo
}

func NewDishUsecase(adminRepo restaurantAdmin.AdminDishRepo,
	imageRepo image.ImageRepo) restaurantAdmin.AdminDishUsecase {

	return &dishUsecase{
		dishRepository:  adminRepo,
		imageRepository: imageRepo,
	}
}

func (a dishUsecase) GetAllDishes(ctx context.Context) ([]models.Dish, error) {
	restaurantAdmin, ok := ctx.Value("Restaurant").(models.Restaurant)
	if !ok {
		failError := errors.FailServerError("failed to convert to models.Restaurant")
		logger.UsecaseLevel().ErrorLog(ctx, failError)
		return nil, failError
	}

	return a.dishRepository.GetAllDishes(ctx, restaurantAdmin.ID)
}

func (a dishUsecase) UpdateDish(ctx context.Context, dish models.Dish) (*models.Dish, error) {
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

	restaurant, ok := ctx.Value("Restaurant").(models.Restaurant)
	if !ok {
		failError := errors.FailServerError("failed to convert to models.Restaurant")
		logger.RepoLevel().ErrorLog(ctx, failError)
		return nil, failError
	}

	dish.Restaurant = restaurant.ID

	err := a.dishRepository.UpdateDish(ctx, dish)
	if err != nil {
		return nil, err
	}
	return &dish, nil
}

func (a dishUsecase) DeleteDish(ctx context.Context, did int) error {
	ok := a.checkRightsForDish(ctx, did)
	if !ok {
		requestError := errors.BadRequestError("No rights to delete a dish")
		logger.UsecaseLevel().ErrorLog(ctx, requestError)
		return requestError
	}

	// удаление изображения
	oldDish, err := a.dishRepository.GetDish(ctx, did)
	if oldDish.Image != config.DefaultAvatar {
		removeFile := strings.Replace(oldDish.Image, config.Repository, "", -1)
		err = a.imageRepository.DeleteImage(ctx, removeFile)
		if err != nil {
			return err
		}
	}

	return a.dishRepository.DeleteDish(ctx, did)
}

func (a dishUsecase) AddDish(ctx context.Context, dish models.Dish) (*models.Dish, error) {
	if dish.Section == 0 {
		requestError := errors.BadRequestError("No section at the dish")
		logger.UsecaseLevel().ErrorLog(ctx, requestError)
		return nil, requestError
	}

	restaurant, ok := ctx.Value("Restaurant").(models.Restaurant)
	if !ok {
		failError := errors.FailServerError("failed to convert to models.Restaurant")
		logger.RepoLevel().ErrorLog(ctx, failError)
		return nil, failError
	}

	dish.Restaurant = restaurant.ID

	id, err := a.dishRepository.AddDish(ctx, dish)
	if err != nil {
		return nil, err
	}
	responseDish := dish
	responseDish.ID = id
	responseDish.Image = config.DefaultAvatar
	return &responseDish, nil
}

func (a dishUsecase) UploadDishImage(ctx context.Context, image *multipart.FileHeader, idDish int) (*models.DishImageResponse, error) {
	if idDish == 0 {
		requestError := errors.BadRequestError("No id at the dish")
		logger.UsecaseLevel().ErrorLog(ctx, requestError)
		return nil, requestError
	}

	ok := a.checkRightsForDish(ctx, idDish)
	if !ok {
		requestError := errors.BadRequestError("No rights to update a dish")
		logger.UsecaseLevel().ErrorLog(ctx, requestError)
		return nil, requestError
	}

	// парсим расширение
	expansion := filepath.Ext(image.Filename)

	uid, localErr := uniq.GetUniqFilename(ctx, image.Filename)
	if localErr != nil {
		return nil, localErr
	}

	// удаление изображения
	oldDish, err := a.dishRepository.GetDish(ctx, idDish)
	if oldDish.Image != config.DefaultAvatar {
		removeFile := strings.Replace(oldDish.Image, config.Repository, "", -1)
		err = a.imageRepository.DeleteImage(ctx, removeFile)
		if err != nil {
			return nil, err
		}
	}

	custFilename := HeadImage + uid + expansion
	err = a.imageRepository.UploadImage(ctx, custFilename, image)
	if err != nil {
		return nil, err
	}

	custFilename = config.Repository + HeadImage + uid + expansion
	err = a.dishRepository.UpdateDishImage(ctx, idDish, custFilename)
	if err != nil {
		return nil, err
	}

	response := &models.DishImageResponse{
		ID:       idDish,
		Filename: custFilename,
	}
	return response, nil
}

func (a dishUsecase) checkRightsForDish(ctx context.Context, idDish int) bool {
	dish, err := a.dishRepository.GetDish(ctx, idDish)
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
