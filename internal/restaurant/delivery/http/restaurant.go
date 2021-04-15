package http

import (
	"strconv"

	"github.com/borscht/backend/internal/models"
	restModel "github.com/borscht/backend/internal/restaurant"
	errors "github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
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

	logger.DeliveryLevel().InfoLog(ctx, logger.Fields{"restaurant": result})

	response := make([]models.Response, 0)
	for i := range result {
		response = append(response, &result[i])
	}
	logger.DeliveryLevel().InfoLog(ctx, logger.Fields{"restaurant": &response})
	return models.SendMoreResponse(c, response...)
}

func (h *RestaurantHandler) GetRestaurantPage(c echo.Context) error {
	ctx := models.GetContext(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		badRequest := errors.BadRequestError(err.Error())
		logger.DeliveryLevel().ErrorLog(ctx, badRequest)
		return models.SendResponseWithError(c, badRequest)
	}

	restaurant, err := h.restaurantUsecase.GetById(ctx, id)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	return models.SendResponse(c, restaurant)
}
