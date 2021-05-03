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

func NewRestaurantHandler(restUCase restModel.RestaurantUsecase) restModel.RestaurantHandler {
	return &RestaurantHandler{
		restaurantUsecase: restUCase,
	}
}

func (h *RestaurantHandler) GetVendor(c echo.Context) error {
	limit, errLimit := strconv.Atoi(c.QueryParam("limit"))
	offset, errOffset := strconv.Atoi(c.QueryParam("offset"))
	latitude := c.QueryParam("latitude")
	longitude := c.QueryParam("longitude")
	name := c.QueryParam("name")

	ctx := models.GetContext(c)

	params := restModel.GetVendorParams{
		Limit:     limit,
		Offset:    offset,
		Address:   true,
		Name:      name,
		Latitude:  latitude,
		Longitude: longitude,
	}
	if longitude == "" || latitude == "" { // адрес не передан
		params.Address = false
	}
	logger.DeliveryLevel().InfoLog(ctx, logger.Fields{"getVendor params": params})

	if errLimit != nil {
		return models.SendResponseWithError(c, errors.BadRequestError(errLimit.Error()))
	}
	if errOffset != nil {
		return models.SendResponseWithError(c, errors.BadRequestError(errOffset.Error()))
	}

	result, err := h.restaurantUsecase.GetVendor(ctx, params)
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

func (h *RestaurantHandler) GetReviews(c echo.Context) error {
	ctx := models.GetContext(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		badRequest := errors.BadRequestError(err.Error())
		logger.DeliveryLevel().ErrorLog(ctx, badRequest)
		return models.SendResponseWithError(c, badRequest)
	}

	reviews, err := h.restaurantUsecase.GetReviews(ctx, id)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	logger.DeliveryLevel().InfoLog(ctx, logger.Fields{"reviews": reviews})

	response := make([]models.Response, 0)
	for i := range reviews {
		response = append(response, &reviews[i])
	}

	return models.SendMoreResponse(c, response...)
}
