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
	HeadImageDish = "static/dish/"
)

type dishUsecase struct {
	FolderSaveImages  string
	dishRepository    restaurantAdmin.AdminDishRepo
	sectionRepository restaurantAdmin.AdminSectionRepo
	imageRepository   image.ImageRepo
}

func NewDishUsecase(adminRepo restaurantAdmin.AdminDishRepo,
	sectionRepo restaurantAdmin.AdminSectionRepo,
	imageRepo image.ImageRepo) restaurantAdmin.AdminDishUsecase {

	return &dishUsecase{
		FolderSaveImages:  config.Static,
		dishRepository:    adminRepo,
		sectionRepository: sectionRepo,
		imageRepository:   imageRepo,
	}
}

func (a dishUsecase) GetAllDishes(ctx context.Context) ([]models.SectionWithDishes, error) {
	restaurant, ok := ctx.Value("Restaurant").(models.RestaurantInfo)
	if !ok {
		failError := errors.FailServerError("failed to convert to models.Restaurant")
		logger.UsecaseLevel().ErrorLog(ctx, failError)
		return nil, failError
	}

	sections, err := a.sectionRepository.GetAllSections(ctx, restaurant.ID)
	if err != nil {
		return nil, err
	}

	response := []models.SectionWithDishes{}
	for _, section := range sections {
		dishes, err := a.dishRepository.GetAllDishes(ctx, section.ID)
		if err != nil {
			return nil, err
		}

		sectionWithDishes := models.SectionWithDishes{
			Dishes:      dishes,
			SectionName: section.Name,
			SectionId:   section.ID,
		}

		response = append(response, sectionWithDishes)
	}

	return response, nil
}

func (a dishUsecase) UpdateDishData(ctx context.Context, dish models.Dish) (*models.Dish, error) {
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

	restaurant, ok := ctx.Value("Restaurant").(models.RestaurantInfo)
	if !ok {
		failError := errors.FailServerError("failed to convert to models.Restaurant")
		logger.RepoLevel().ErrorLog(ctx, failError)
		return nil, failError
	}

	dish.Restaurant = restaurant.ID

	err := a.dishRepository.UpdateDishData(ctx, dish)
	if err != nil {
		return nil, err
	}
	return &dish, nil
}

func (a dishUsecase) DeleteDish(ctx context.Context, did int) (*models.DeleteSuccess, error) {
	ok := a.checkRightsForDish(ctx, did)
	if !ok {
		requestError := errors.BadRequestError("No rights to delete a dish")
		logger.UsecaseLevel().ErrorLog(ctx, requestError)
		return nil, requestError
	}

	// удаление изображения
	oldDish, err := a.dishRepository.GetDish(ctx, did)
	if oldDish.Image != config.DefaultDishImage {
		removeFile := a.FolderSaveImages + strings.Replace(oldDish.Image, config.Repository, "", -1)
		err = a.imageRepository.DeleteImage(ctx, removeFile)
		if err != nil {
			return nil, err
		}
	}

	return &models.DeleteSuccess{ID: did}, a.dishRepository.DeleteDish(ctx, did)
}

func (a dishUsecase) AddDish(ctx context.Context, dish models.Dish) (*models.Dish, error) {
	if dish.Section == 0 {
		requestError := errors.BadRequestError("No section at the dish")
		logger.UsecaseLevel().ErrorLog(ctx, requestError)
		return nil, requestError
	}

	restaurant, ok := ctx.Value("Restaurant").(models.RestaurantInfo)
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
	responseDish.Image = config.DefaultDishImage
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
	if oldDish.Image != config.DefaultDishImage {
		removeFile := a.FolderSaveImages + strings.Replace(oldDish.Image, config.Repository, "", -1)
		err = a.imageRepository.DeleteImage(ctx, removeFile)
		if err != nil {
			return nil, err
		}
	}

	custFilename := a.FolderSaveImages + HeadImageDish + uid + expansion
	err = a.imageRepository.UploadImage(ctx, custFilename, image)
	if err != nil {
		return nil, err
	}

	custFilename = config.Repository + HeadImageDish + uid + expansion
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

	restaurant, ok := ctx.Value("Restaurant").(models.RestaurantInfo)
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
