package usecase

import (
	"context"

	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/restaurantAdmin"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
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
	restaurant, ok := ctx.Value("Restaurant").(models.Restaurant)
	if !ok {
		failError := errors.FailServerError("failed to convert to models.Restaurant")
		logger.RepoLevel().ErrorLog(ctx, failError)
		return nil, failError
	}

	section.Restaurant = restaurant.ID

	id, err := s.sectionRepository.AddSection(ctx, section)
	if err != nil {
		return nil, err
	}

	section.ID = id
	return &section, nil
}

func (s sectionUsecase) UpdateSection(ctx context.Context, section models.Section) (*models.Section, error) {
	if section.ID == 0 {
		requestError := errors.BadRequestError("No id at the section")
		logger.UsecaseLevel().ErrorLog(ctx, requestError)
		return nil, requestError
	}

	// проверка прав на update
	ok := s.checkRightsForSection(ctx, section.ID)
	if !ok {
		requestError := errors.BadRequestError("No rights to update a section")
		logger.UsecaseLevel().ErrorLog(ctx, requestError)
		return nil, requestError
	}

	restaurant, ok := ctx.Value("Restaurant").(models.Restaurant)
	if !ok {
		failError := errors.FailServerError("failed to convert to models.Restaurant")
		logger.RepoLevel().ErrorLog(ctx, failError)
		return nil, failError
	}

	section.Restaurant = restaurant.ID

	err := s.sectionRepository.UpdateSection(ctx, section)
	if err != nil {
		return nil, err
	}
	return &section, nil
}

func (s sectionUsecase) DeleteSection(ctx context.Context, sid int) error {
	ok := s.checkRightsForSection(ctx, sid)
	if !ok {
		requestError := errors.BadRequestError("No rights to delete a section")
		logger.UsecaseLevel().ErrorLog(ctx, requestError)
		return requestError
	}

	return s.sectionRepository.DeleteSection(ctx, sid)
}

func (s sectionUsecase) checkRightsForSection(ctx context.Context, idSection int) bool {
	section, err := s.sectionRepository.GetSection(ctx, idSection)
	if err != nil {
		return false
	}

	restaurant, ok := ctx.Value("Restaurant").(models.Restaurant)
	if !ok {
		failError := errors.FailServerError("failed to convert to models.Restaurant")
		logger.UsecaseLevel().ErrorLog(ctx, failError)
		return false
	}

	if restaurant.ID != section.Restaurant {
		return false
	}

	return true
}
