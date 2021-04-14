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
	HeadImageRestaurant = "static/restaurant/"
)

type restaurantUsecase struct {
	restaurantRepository restaurantAdmin.AdminRestaurantRepo
	imageRepository      image.ImageRepo
}

func NewRestaurantUsecase(adminRepo restaurantAdmin.AdminRestaurantRepo,
	imageRepo image.ImageRepo) restaurantAdmin.AdminRestaurantUsecase {

	return &restaurantUsecase{
		restaurantRepository: adminRepo,
		imageRepository:      imageRepo,
	}
}

func (a restaurantUsecase) UpdateRestaurantData(ctx context.Context, restaurant models.RestaurantUpdateData) (
	*models.RestaurantInfo, error) {

	restaurantAdmin, ok := ctx.Value("Restaurant").(models.RestaurantInfo)
	if !ok {
		failError := errors.FailServerError("failed to convert to models.Restaurant")
		logger.UsecaseLevel().ErrorLog(ctx, failError)
		return nil, failError
	}

	restaurant.ID = restaurantAdmin.ID
	err := a.restaurantRepository.UpdateRestaurantData(ctx, restaurant)
	if err != nil {
		return nil, err
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
	}
	logger.UsecaseLevel().DebugLog(ctx, logger.Fields{"restaurant": restaurantResponse})

	return restaurantResponse, nil
}

func (a restaurantUsecase) CreateRestaurant(ctx context.Context, restaurant models.RestaurantInfo) (*models.RestaurantInfo, error) {
	restaurant.Avatar = config.DefaultRestaurantImage

	id, err := a.restaurantRepository.CreateRestaurant(ctx, restaurant)
	if err != nil {
		return nil, err
	}
	restaurant.ID = id
	return &restaurant, nil
}

func (a restaurantUsecase) CheckRestaurantExists(ctx context.Context, user models.RestaurantAuth) (*models.RestaurantInfo, error) {
	return a.restaurantRepository.CheckRestaurantExists(ctx, user)
}

func (a restaurantUsecase) GetByRid(ctx context.Context, rid int) (*models.RestaurantInfo, error) {
	return a.restaurantRepository.GetByRid(ctx, rid)
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

	if restaurant.Avatar != config.DefaultRestaurantImage {
		removeFile := strings.Replace(restaurant.Avatar, config.Repository, "", -1)
		err := a.imageRepository.DeleteImage(ctx, removeFile)
		if err != nil {
			return nil, err
		}
	}

	custFilename := HeadImageRestaurant + uid + expansion
	err := a.imageRepository.UploadImage(ctx, custFilename, image)
	if err != nil {
		return nil, err
	}

	custFilename = config.Repository + HeadImageRestaurant + uid + expansion
	err = a.restaurantRepository.UpdateRestaurantImage(ctx, restaurant.ID, custFilename)
	if err != nil {
		return nil, err
	}

	response := &models.RestaurantImageResponse{
		Filename: custFilename,
	}
	return response, nil
}
