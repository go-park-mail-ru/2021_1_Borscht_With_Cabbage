package http

import (
	"strconv"
	"strings"

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
	ctx := models.GetContext(c)

	limit, errLimit := strconv.Atoi(c.QueryParam("limit"))
	offset, errOffset := strconv.Atoi(c.QueryParam("offset"))
	categories := strings.Split(c.QueryParam("category"), ",")
	logger.DeliveryLevel().DebugLog(ctx, logger.Fields{"categories": categories, "size": len(categories)})

	if errLimit != nil {
		requestError := errors.BadRequestError(errLimit.Error())
		logger.DeliveryLevel().ErrorLog(ctx, requestError)
		return models.SendResponseWithError(c, requestError)
	}
	if errOffset != nil {
		requestError := errors.BadRequestError(errOffset.Error())
		logger.DeliveryLevel().ErrorLog(ctx, requestError)
		return models.SendResponseWithError(c, requestError)
	}

	result, err := h.restaurantUsecase.GetVendor(ctx, models.RestaurantRequest{
		Limit:      limit,
		Offset:     offset,
		Categories: categories,
	})
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

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
