package usecase

import (
	"context"
	"github.com/borscht/backend/internal/models"
	adminMock "github.com/borscht/backend/internal/restaurantAdmin/mocks"
	"github.com/borscht/backend/utils/logger"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewSectionUsecase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sectionRepoMock := adminMock.NewMockAdminSectionRepo(ctrl)

	sectionUsecase := NewSectionUsecase(sectionRepoMock)
	if sectionUsecase == nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestSectionUsecase_AddSection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sectionRepoMock := adminMock.NewMockAdminSectionRepo(ctrl)

	sectionUsecase := NewSectionUsecase(sectionRepoMock)
	if sectionUsecase == nil {
		t.Errorf("incorrect result")
		return
	}

	restaurantAdmin := models.RestaurantInfo{
		ID: 1,
	}
	c := context.Background()
	ctx := context.WithValue(c, "Restaurant", restaurantAdmin)

	section := models.Section{}
	sectionWithRest := section
	sectionWithRest.Restaurant = restaurantAdmin.ID

	sectionRepoMock.EXPECT().AddSection(ctx, sectionWithRest).Return(1, nil)

	var err error
	sectionResponse := new(models.Section)
	sectionResponse, err = sectionUsecase.AddSection(ctx, section)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	require.EqualValues(t, sectionResponse.ID, 1)
}

func TestSectionUsecase_UpdateSection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sectionRepoMock := adminMock.NewMockAdminSectionRepo(ctrl)

	sectionUsecase := NewSectionUsecase(sectionRepoMock)
	if sectionUsecase == nil {
		t.Errorf("incorrect result")
		return
	}

	restaurantAdmin := models.RestaurantInfo{
		ID: 1,
	}
	c := context.Background()
	ctx := context.WithValue(c, "Restaurant", restaurantAdmin)

	section := models.Section{
		ID: 1,
	}
	sectionFull := section
	sectionFull.ID = 1
	sectionFull.Restaurant = 1

	sectionRepoMock.EXPECT().GetSection(ctx, 1).Return(&sectionFull, nil)
	sectionRepoMock.EXPECT().UpdateSection(ctx, sectionFull).Return(nil)

	logger.InitLogger()
	var err error
	sectionResponse := new(models.Section)
	sectionResponse, err = sectionUsecase.UpdateSection(ctx, section)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	require.EqualValues(t, sectionResponse.ID, 1)
}

func TestSectionUsecase_UpdateSection_NoRights(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sectionRepoMock := adminMock.NewMockAdminSectionRepo(ctrl)

	sectionUsecase := NewSectionUsecase(sectionRepoMock)
	if sectionUsecase == nil {
		t.Errorf("incorrect result")
		return
	}

	c := context.Background()
	ctx := context.WithValue(c, "Restaurant", 1)

	section := models.Section{
		ID: 1,
	}
	sectionFull := section
	sectionFull.ID = 1
	sectionFull.Restaurant = 1

	sectionRepoMock.EXPECT().GetSection(ctx, 1).Return(&sectionFull, nil)

	logger.InitLogger()
	_, err := sectionUsecase.UpdateSection(ctx, section)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
}

func TestSectionUsecase_UpdateSection_SectionId0(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sectionRepoMock := adminMock.NewMockAdminSectionRepo(ctrl)

	sectionUsecase := NewSectionUsecase(sectionRepoMock)
	if sectionUsecase == nil {
		t.Errorf("incorrect result")
		return
	}

	c := context.Background()
	ctx := context.WithValue(c, "Restaurant", 1)

	section := models.Section{
		ID: 0,
	}

	logger.InitLogger()
	_, err := sectionUsecase.UpdateSection(ctx, section)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
}

func TestSectionUsecase_DeleteSection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sectionRepoMock := adminMock.NewMockAdminSectionRepo(ctrl)

	sectionUsecase := NewSectionUsecase(sectionRepoMock)
	if sectionUsecase == nil {
		t.Errorf("incorrect result")
		return
	}

	restaurantAdmin := models.RestaurantInfo{
		ID: 1,
	}
	c := context.Background()
	ctx := context.WithValue(c, "Restaurant", restaurantAdmin)

	section := models.Section{
		ID: 1,
	}
	sectionFull := section
	sectionFull.ID = 1
	sectionFull.Restaurant = 1

	sectionRepoMock.EXPECT().GetSection(ctx, 1).Return(&sectionFull, nil)
	sectionRepoMock.EXPECT().DeleteSection(ctx, section.ID).Return(nil)

	logger.InitLogger()
	var err error
	deleted := new(models.DeleteSuccess)
	deleted, err = sectionUsecase.DeleteSection(ctx, section.ID)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	require.EqualValues(t, deleted.ID, 1)
}
