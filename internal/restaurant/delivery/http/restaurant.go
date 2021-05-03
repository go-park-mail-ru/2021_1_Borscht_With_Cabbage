package http

import (
	"context"
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

	params := make([]string, 0)
	params = append(params, c.QueryParam("limit"))
	params = append(params, c.QueryParam("offset"))
	params = append(params, c.QueryParam("time"))
	params = append(params, c.QueryParam("receipt"))
	paramsNumber, err := AtoiParams(ctx, params...)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	rating, parseErr := strconv.ParseFloat(c.QueryParam("rating"), 64)
	if parseErr != nil {
		requestError := errors.BadRequestError(err.Error())
		logger.DeliveryLevel().ErrorLog(ctx, requestError)
		return models.SendResponseWithError(c, requestError)
	}

	categories := strings.Split(c.QueryParam("category"), ",")
	logger.DeliveryLevel().DebugLog(ctx, logger.Fields{"categories": categories, "size": len(categories)})

	result, err := h.restaurantUsecase.GetVendor(ctx, models.RestaurantRequest{
		Limit:      paramsNumber[0],
		Offset:     paramsNumber[1],
		Categories: categories,
		Time:       paramsNumber[2],
		Receipt:    paramsNumber[3],
		Rating:     rating,
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

func AtoiParams(ctx context.Context, params ...string) ([]int, error) {
	result := make([]int, 0)
	for _, value := range params {
		valueNumber, err := strconv.Atoi(value)
		if err != nil {
			requestError := errors.BadRequestError(err.Error())
			logger.DeliveryLevel().ErrorLog(ctx, requestError)
			return nil, requestError
		}

		result = append(result, valueNumber)
	}
	return result, nil
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
