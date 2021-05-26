package usecase

import (
	"context"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/borscht/backend/config"
	"github.com/borscht/backend/configProject"
	"github.com/borscht/backend/internal/image"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/restaurantAdmin"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
	"github.com/borscht/backend/utils/uniq"
)

const (
	HeadImageRestaurant = "static/restaurant/"
)

type restaurantUsecase struct {
	FolderSaveImages     string
	restaurantRepository restaurantAdmin.AdminRestaurantRepo
	imageRepository      image.ImageRepo
}

func NewRestaurantUsecase(adminRepo restaurantAdmin.AdminRestaurantRepo,
	imageRepo image.ImageRepo) restaurantAdmin.AdminRestaurantUsecase {

	return &restaurantUsecase{
		FolderSaveImages: config.ConfigStatic.Folder,

		restaurantRepository: adminRepo,
		imageRepository:      imageRepo,
	}
}

func convertAddressToStr(latitude, longitude string) (float64, float64, error) {
	latitudeNum, errLat := strconv.ParseFloat(latitude, 64)
	longitudeNum, errLon := strconv.ParseFloat(longitude, 64)
	if errLat != nil || errLon != nil {
		return 0, 0, errors.NewErrorWithMessage("converting address to float error")
	}
	return latitudeNum, longitudeNum, nil
}

func (a restaurantUsecase) AddCategories(ctx context.Context, categories models.Categories) error {
	restaurantAdmin, ok := ctx.Value("Restaurant").(models.RestaurantInfo)
	if !ok {
		failError := errors.FailServerError("failed to convert to models.Restaurant")
		logger.UsecaseLevel().ErrorLog(ctx, failError)
		return failError
	}
	err := a.restaurantRepository.DeleteAllCategories(ctx, restaurantAdmin.ID)
	if err != nil {
		return err
	}

	return a.restaurantRepository.AddCategories(ctx, restaurantAdmin.ID, categories.CategoriesID)
}

func (a restaurantUsecase) correcRestaurantData(ctx context.Context, newRestaurant *models.RestaurantUpdateData,
	oldRestaurant *models.RestaurantInfo) {

	newRestaurant.ID = oldRestaurant.ID
	if newRestaurant.Title == "" {
		newRestaurant.Title = oldRestaurant.Title
	}
	if newRestaurant.AdminEmail == "" {
		newRestaurant.AdminEmail = oldRestaurant.AdminEmail
	}
	if newRestaurant.AdminPhone == "" {
		newRestaurant.AdminPhone = oldRestaurant.AdminPhone
	}
	logger.UsecaseLevel().DebugLog(ctx, logger.Fields{
		"Address new": newRestaurant.Address,
		"Address old": oldRestaurant.Address,
	})
	if newRestaurant.Address.Name == "" {
		newRestaurant.Address.Name = oldRestaurant.Address.Name
		newRestaurant.Address.Latitude = oldRestaurant.Address.Latitude
		newRestaurant.Address.Longitude = oldRestaurant.Address.Longitude
		newRestaurant.Address.Radius = oldRestaurant.Address.Radius
	}
	if newRestaurant.DeliveryCost == 0 {
		newRestaurant.DeliveryCost = oldRestaurant.DeliveryCost
	}
	if newRestaurant.Description == "" {
		newRestaurant.Description = oldRestaurant.Description
	}
	if newRestaurant.Address.Radius == 0 {
		newRestaurant.Address.Radius = oldRestaurant.Address.Radius
	}
}

func (a restaurantUsecase) GetCategories(ctx context.Context, rid int) (*models.Categories, error) {
	categories, err := a.restaurantRepository.GetCategories(ctx, rid)
	if err != nil {
		return nil, err
	}
	return &models.Categories{CategoriesID: categories}, nil
}

func (a restaurantUsecase) AddAddress(ctx context.Context, rid int, address models.Address) error {
	return a.restaurantRepository.AddAddress(ctx, rid, address)
}

func (a restaurantUsecase) GetAddress(ctx context.Context, rid int) (*models.Address, error) {
	return a.restaurantRepository.GetAddress(ctx, rid)
}

func (a restaurantUsecase) UpdateRestaurantData(ctx context.Context, restaurant models.RestaurantUpdateData) (
	*models.SuccessRestaurantResponse, error) {

	restaurantAdmin, ok := ctx.Value("Restaurant").(models.RestaurantInfo)
	if !ok {
		failError := errors.FailServerError("failed to convert to models.Restaurant")
		logger.UsecaseLevel().ErrorLog(ctx, failError)
		return nil, failError
	}

	a.correcRestaurantData(ctx, &restaurant, &restaurantAdmin)
	err := a.restaurantRepository.UpdateRestaurantData(ctx, restaurant)
	if err != nil {
		return nil, err
	}

	err = a.restaurantRepository.UpdateAddress(ctx, restaurant.ID, restaurant.Address)
	if err != nil {
		return nil, err
	}

	if restaurant.Categories != nil && len(restaurant.Categories) != 0 {
		categories := models.Categories{
			CategoriesID: restaurant.Categories,
		}
		err = a.AddCategories(ctx, categories)
		if err != nil {
			return nil, err
		}
	} else {
		categories, err := a.GetCategories(ctx, restaurant.ID)
		if err != nil {
			return nil, err
		}
		restaurant.Categories = categories.CategoriesID
	}

	logger.UsecaseLevel().DebugLog(ctx, logger.Fields{"restaurant": restaurant})
	restaurantResponse := &models.RestaurantInfo{
		ID:           restaurant.ID,
		AdminEmail:   restaurant.AdminEmail,
		AdminPhone:   restaurant.AdminPhone,
		Title:        restaurant.Title,
		Description:  restaurant.Description,
		Rating:       restaurantAdmin.Rating,
		AvgCheck:     restaurantAdmin.AvgCheck,
		DeliveryCost: restaurant.DeliveryCost,
		Avatar:       restaurantAdmin.Avatar,
		Address:      restaurant.Address,
		Categories:   restaurant.Categories,
	}
	logger.UsecaseLevel().DebugLog(ctx, logger.Fields{"restaurant": restaurantResponse})

	return &models.SuccessRestaurantResponse{
		RestaurantInfo: *restaurantResponse,
		Role:           configProject.RoleAdmin,
	}, nil
}

func (a restaurantUsecase) UploadRestaurantImage(ctx context.Context, image *multipart.FileHeader) (*models.RestaurantImageResponse, error) {
	// парсим расширение
	expansion := filepath.Ext(image.Filename)

	uid, localErr := uniq.GetUniqFilename(ctx, image.Filename)
	if localErr != nil {
		return nil, localErr
	}

	// удаление изображения
	restaurant, ok := ctx.Value("Restaurant").(models.RestaurantInfo)
	if !ok {
		failError := errors.FailServerError("failed to convert to models.Restaurant")
		logger.UsecaseLevel().ErrorLog(ctx, failError)
		return nil, failError
	}

	if restaurant.Avatar != config.ConfigStatic.DefaultRestaurantImage {
		removeFile := a.FolderSaveImages + strings.Replace(restaurant.Avatar, config.ConfigStatic.Repository, "", -1)
		err := a.imageRepository.DeleteImage(ctx, removeFile)
		if err != nil {
			return nil, err
		}
	}

	custFilename := a.FolderSaveImages + HeadImageRestaurant + uid + expansion
	err := a.imageRepository.UploadImage(ctx, custFilename, image)
	if err != nil {
		return nil, err
	}

	custFilename = config.ConfigStatic.Repository + HeadImageRestaurant + uid + expansion
	err = a.restaurantRepository.UpdateRestaurantImage(ctx, restaurant.ID, custFilename)
	if err != nil {
		return nil, err
	}

	response := &models.RestaurantImageResponse{
		Filename: custFilename,
	}
	return response, nil
}
