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

	request := models.RestaurantRequest{
		Limit:      paramsNumber[0],
		Offset:     paramsNumber[1],
		Categories: categories,
		Time:       paramsNumber[2],
		Receipt:    paramsNumber[3],
		Rating:     rating,
		Address:    true,
	}
	request.LatitudeUser, _ = strconv.ParseFloat(c.QueryParam("latitude"), 64)
	request.LongitudeUser, _ = strconv.ParseFloat(c.QueryParam("longitude"), 64)

	if request.LatitudeUser == 0 || request.LongitudeUser == 0 { // адрес не передан
		request.Address = false
	}
	logger.DeliveryLevel().InfoLog(ctx, logger.Fields{"getVendor params": request})

	result, err := h.restaurantUsecase.GetVendor(ctx, request)
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
	latitude, _ := strconv.ParseFloat(c.QueryParam("latitude"), 64)
	longitude, _ := strconv.ParseFloat(c.QueryParam("longitude"), 64)

	coordinates := models.Coordinates{Latitude: latitude, Longitude: longitude}
	ctx := models.GetContext(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		badRequest := errors.BadRequestError(err.Error())
		logger.DeliveryLevel().ErrorLog(ctx, badRequest)
		return models.SendResponseWithError(c, badRequest)
	}

	restaurant, err := h.restaurantUsecase.GetById(ctx, coordinates, id)
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

func (h *RestaurantHandler) GetRecommendations(c echo.Context) error {
	ctx := models.GetContext(c)

	latitude, _ := strconv.ParseFloat(c.QueryParam("latitude"), 64)
	longitude, _ := strconv.ParseFloat(c.QueryParam("longitude"), 64)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		badRequest := errors.BadRequestError(err.Error())
		logger.DeliveryLevel().ErrorLog(ctx, badRequest)
		return models.SendResponseWithError(c, badRequest)
	}

	recommendationParams := models.RecommendationsParams{
		Id:            id,
		LatitudeUser:  latitude,
		LongitudeUser: longitude,
	}

	recommendations, err := h.restaurantUsecase.GetRecommendations(ctx, recommendationParams)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	logger.DeliveryLevel().InfoLog(ctx, logger.Fields{"recommendations": recommendations})

	response := make([]models.Response, 0)
	for i := range recommendations {
		response = append(response, &recommendations[i])
	}
	return models.SendMoreResponse(c, response...)
}
