package usecase

import (
	"context"
	"testing"

	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/image/mocks"
	"github.com/borscht/backend/internal/models"
	adminMock "github.com/borscht/backend/internal/restaurantAdmin/mocks"
	"github.com/borscht/backend/utils/logger"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestNewDishUsecase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dishRepoMock := adminMock.NewMockAdminDishRepo(ctrl)
	sectionRepoMock := adminMock.NewMockAdminSectionRepo(ctrl)
	imageRepoMock := mocks.NewMockImageRepo(ctrl)

	dishUsecase := NewDishUsecase(dishRepoMock, sectionRepoMock, imageRepoMock)
	if dishUsecase == nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestDishUsecase_GetAllDishes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dishRepoMock := adminMock.NewMockAdminDishRepo(ctrl)
	sectionRepoMock := adminMock.NewMockAdminSectionRepo(ctrl)
	imageRepoMock := mocks.NewMockImageRepo(ctrl)
	dishUsecase := NewDishUsecase(dishRepoMock, sectionRepoMock, imageRepoMock)

	restaurant := models.RestaurantInfo{
		ID: 1,
	}
	c := context.Background()
	ctx := context.WithValue(c, "Restaurant", restaurant)

	sections := []models.Section{
		{ID: 1, Name: "yum"},
	}
	dishes := []models.Dish{
		{ID: 1, Name: "dish1"},
		{ID: 2, Name: "dish2"},
	}

	sectionRepoMock.EXPECT().GetAllSections(ctx, 1).Return(sections, nil)
	dishRepoMock.EXPECT().GetAllDishes(ctx, sections[0].ID).Return(dishes, nil)

	sectionWithDishes := make([]models.SectionWithDishes, 0)
	var err error
	sectionWithDishes, err = dishUsecase.GetAllDishes(ctx)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	require.EqualValues(t, sectionWithDishes[0].SectionName, "yum")
}

func TestDishUsecase_UpdateDishData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dishRepoMock := adminMock.NewMockAdminDishRepo(ctrl)
	sectionRepoMock := adminMock.NewMockAdminSectionRepo(ctrl)
	imageRepoMock := mocks.NewMockImageRepo(ctrl)
	dishUsecase := NewDishUsecase(dishRepoMock, sectionRepoMock, imageRepoMock)

	restaurant := models.RestaurantInfo{
		ID: 1,
	}
	c := context.Background()
	ctx := context.WithValue(c, "Restaurant", restaurant)

	dish := models.Dish{
		ID: 1, Name: "dish1", Restaurant: 1,
	}
	dishFull := dish
	dishFull.Restaurant = 1

	dishRepoMock.EXPECT().GetDish(ctx, dish.ID).Return(&dishFull, nil)
	dishRepoMock.EXPECT().UpdateDishData(ctx, dish).Return(nil)

	dishResponse := new(models.Dish)
	var err error
	dishResponse, err = dishUsecase.UpdateDishData(ctx, dish)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	require.EqualValues(t, dishResponse.ID, 1)
}

func TestDishUsecase_UpdateDishData_RightsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dishRepoMock := adminMock.NewMockAdminDishRepo(ctrl)
	sectionRepoMock := adminMock.NewMockAdminSectionRepo(ctrl)
	imageRepoMock := mocks.NewMockImageRepo(ctrl)
	dishUsecase := NewDishUsecase(dishRepoMock, sectionRepoMock, imageRepoMock)

	c := context.Background()
	ctx := context.WithValue(c, "Restaurant", "e")

	dish := models.Dish{
		ID: 1, Name: "dish1", Restaurant: 1,
	}

	logger.InitLogger()
	dishRepoMock.EXPECT().GetDish(ctx, dish.ID).Return(&dish, nil)

	_, err := dishUsecase.UpdateDishData(ctx, dish)
	if err == nil {
		t.Errorf("unexpected err")
		return
	}
}

func TestDishUsecase_DeleteDish(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dishRepoMock := adminMock.NewMockAdminDishRepo(ctrl)
	sectionRepoMock := adminMock.NewMockAdminSectionRepo(ctrl)
	imageRepoMock := mocks.NewMockImageRepo(ctrl)
	dishUsecase := NewDishUsecase(dishRepoMock, sectionRepoMock, imageRepoMock)

	restaurant := models.RestaurantInfo{
		ID: 1,
	}
	c := context.Background()
	ctx := context.WithValue(c, "Restaurant", restaurant)

	dish := models.Dish{
		ID: 1, Restaurant: 1,
		Image: config.ConfigStatic.DefaultDishImage,
	}
	dishFull := dish
	dishFull.Restaurant = 1

	dishRepoMock.EXPECT().GetDish(ctx, dish.ID).Return(&dishFull, nil)
	dishRepoMock.EXPECT().GetDish(ctx, dish.ID).Return(&dish, nil)
	dishRepoMock.EXPECT().DeleteDish(ctx, dish.ID).Return(nil)

	deleted := new(models.DeleteSuccess)
	var err error
	deleted, err = dishUsecase.DeleteDish(ctx, dish.ID)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	require.EqualValues(t, deleted.ID, 1)
}

func TestDishUsecase_AddDish(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dishRepoMock := adminMock.NewMockAdminDishRepo(ctrl)
	sectionRepoMock := adminMock.NewMockAdminSectionRepo(ctrl)
	imageRepoMock := mocks.NewMockImageRepo(ctrl)
	dishUsecase := NewDishUsecase(dishRepoMock, sectionRepoMock, imageRepoMock)

	restaurant := models.RestaurantInfo{
		ID: 1,
	}
	c := context.Background()
	ctx := context.WithValue(c, "Restaurant", restaurant)

	dish := models.Dish{
		Name: "dish1", Restaurant: 1, Section: 1,
	}
	dishID := 1

	dishRepoMock.EXPECT().AddDish(ctx, dish).Return(dishID, nil)

	dishResponse := new(models.Dish)
	var err error
	dishResponse, err = dishUsecase.AddDish(ctx, dish)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	require.EqualValues(t, dishResponse.ID, 1)
	require.EqualValues(t, dishResponse.Name, "dish1")
	require.EqualValues(t, dishResponse.Image, config.ConfigStatic.DefaultDishImage)
}

func TestDishUsecase_UploadDishImage(t *testing.T) {
	// TODO
}
