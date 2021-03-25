package http

import (
	"github.com/borscht/backend/internal/models"
	_restModel "github.com/borscht/backend/internal/restaurant"
	_errors "github.com/borscht/backend/utils"
	"github.com/labstack/echo/v4"
	"strconv"
)

type RestaurantHandler struct {
	restaurantUsecase _restModel.RestaurantUsecase
}

// NewArticleHandler will initialize the articles/ resources endpoint
func NewRestaurantHandler(restUCase _restModel.RestaurantUsecase) _restModel.RestaurantHandler {
	return &RestaurantHandler{
		restaurantUsecase: restUCase,
	}
}

func (h *RestaurantHandler) GetVendor(c echo.Context) error {
	cc := c.(*models.CustomContext)
	limit, errLimit := strconv.Atoi(c.QueryParam("limit"))
	offset, errOffset := strconv.Atoi(c.QueryParam("offset"))

	if errLimit != nil {
		return cc.SendResponseWithError(_errors.BadRequest(errLimit.Error()))
	}
	if errOffset != nil {
		return cc.SendResponseWithError(_errors.BadRequest(errOffset.Error()))
	}

	result, err := h.restaurantUsecase.GetVendor(limit, offset)
	if err != nil {
		return cc.SendResponseWithError(err)
	}
	return cc.SendResponse(result)
}

func (h *RestaurantHandler) GetRestaurantPage(c echo.Context) error {
	cc := c.(*models.CustomContext)

	id := c.Param("id")

	restaurant, isItExists, err := h.restaurantUsecase.GetById(id)
	if err != nil {
		return cc.SendResponseWithError(err)
	}
	if !isItExists {
		return cc.SendResponseWithError(_errors.BadRequest("error with request data"))
	}

	return cc.SendResponse(restaurant)
}
