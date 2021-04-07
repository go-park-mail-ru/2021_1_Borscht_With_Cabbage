package http

import (
	adminModel "github.com/borscht/backend/internal/restaurantAdmin"
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
	return nil
}

func (s SectionHandler) DeleteSection(c echo.Context) error {
	return nil
}

func (s SectionHandler) UpdateSection(c echo.Context) error {
	return nil
}
