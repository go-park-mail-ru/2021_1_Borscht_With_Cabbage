package http

import (
	"github.com/borscht/backend/internal/models"
	adminModel "github.com/borscht/backend/internal/restaurantAdmin"
	errors "github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
	"github.com/labstack/echo/v4"
)

type SectionHandler struct {
	SectionUsecase adminModel.AdminSectionUsecase
}

func NewSectionHandler(sectionUCase adminModel.AdminSectionUsecase) adminModel.AdminSectionHandler {
	return &SectionHandler{
		SectionUsecase: sectionUCase,
	}
}

func (s SectionHandler) AddSection(c echo.Context) error {
	ctx := models.GetContext(c)

	newSection := new(models.Section)
	if err := c.Bind(newSection); err != nil {
		sendErr := errors.BadRequestError(err.Error())
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}

	response, err := s.SectionUsecase.AddSection(ctx, *newSection)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	return models.SendResponse(c, response)
}

func (s SectionHandler) DeleteSection(c echo.Context) error {
	ctx := models.GetContext(c)

	section := new(models.Section)
	if err := c.Bind(section); err != nil {
		sendErr := errors.BadRequestError(err.Error())
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}

	response, err := s.SectionUsecase.DeleteSection(ctx, section.ID)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	return models.SendResponse(c, response)
}

func (s SectionHandler) UpdateSection(c echo.Context) error {
	ctx := models.GetContext(c)

	updateSection := new(models.Section)
	if err := c.Bind(updateSection); err != nil {
		sendErr := errors.BadRequestError(err.Error())
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}

	response, err := s.SectionUsecase.UpdateSection(ctx, *updateSection)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	return models.SendResponse(c, response)
}
