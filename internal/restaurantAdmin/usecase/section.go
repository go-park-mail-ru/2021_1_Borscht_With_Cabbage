package usecase

import (
	"context"

	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/restaurantAdmin"
)

type sectionUsecase struct {
	sectionRepository restaurantAdmin.AdminSectionRepo
}

func NewSectionUsecase(sectionRepo restaurantAdmin.AdminSectionRepo) restaurantAdmin.AdminSectionUsecase {

	return &sectionUsecase{
		sectionRepository: sectionRepo,
	}
}

func (s sectionUsecase) AddSection(ctx context.Context, section models.Section) (*models.Section, error) {
	return nil, nil
}

func (s sectionUsecase) UpdateSection(ctx context.Context, section models.Section) (*models.Section, error) {
	return nil, nil
}

func (s sectionUsecase) DeleteSection(ctx context.Context, sid int) error {
	return nil
}
