package http

import (
	"strconv"

	"github.com/borscht/backend/internal/models"
	restModel "github.com/borscht/backend/internal/restaurant"
	errors "github.com/borscht/backend/utils/errors"
	"github.com/labstack/echo/v4"
)

type RestaurantHandler struct {
	restaurantUsecase restModel.RestaurantUsecase
}

// NewArticleHandler will initialize the articles/ resources endpoint
func NewRestaurantHandler(restUCase restModel.RestaurantUsecase) restModel.RestaurantHandler {
	return &RestaurantHandler{
		restaurantUsecase: restUCase,
	}
}

func (h *RestaurantHandler) GetVendor(c echo.Context) error {
	limit, errLimit := strconv.Atoi(c.QueryParam("limit"))
	offset, errOffset := strconv.Atoi(c.QueryParam("offset"))

	ctx := models.GetContext(c)

	if errLimit != nil {
		return models.SendResponseWithError(c, errors.BadRequestError(errLimit.Error()))
	}
	if errOffset != nil {
		return models.SendResponseWithError(c, errors.BadRequestError(errOffset.Error()))
	}

	result, err := h.restaurantUsecase.GetVendor(ctx, limit, offset)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	return models.SendResponse(c, result)
}

func (h *RestaurantHandler) GetRestaurantPage(c echo.Context) error {
	id := c.Param("id")
	ctx := models.GetContext(c)

	restaurant, err := h.restaurantUsecase.GetById(ctx, id)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	return models.SendResponse(c, restaurant)
}
