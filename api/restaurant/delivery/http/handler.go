package http

import (
	"backend/api/domain"
	errors "backend/utils"
	"github.com/labstack/echo/v4"
	"strconv"
)

type RestaurantHandler struct {
	restaurantUsecase domain.RestaurantUsecase
}

// NewArticleHandler will initialize the articles/ resources endpoint
func NewRestaurantHandler(e *echo.Echo, restUCase domain.RestaurantUsecase) {
	handler := &RestaurantHandler{
		restaurantUsecase: restUCase,
	}

	e.GET("/:id", handler.GetRestaurantPage)
	e.GET("/", handler.GetVendor)
	e.GET("/restaurants", handler.GetVendor)
}

func (h *RestaurantHandler) GetVendor(c echo.Context) error {
	cc := c.(*domain.CustomContext)

	Limit, errLimit := strconv.Atoi(c.QueryParam("limit"))
	Offset, errOffset := strconv.Atoi(c.QueryParam("offset"))

	if errLimit != nil {
		return cc.SendERR(errors.BadRequest(errLimit.Error()))
	}
	if errOffset != nil {
		return cc.SendERR(errors.BadRequest(errOffset.Error()))
	}


	result, err := h.restaurantUsecase.GetSlice(cc, Limit, Offset)
	if err != nil {
		return cc.SendERR(err)
	}
	return cc.SendOK(result)
}

func (h *RestaurantHandler) GetRestaurantPage(c echo.Context) error {
	cc := c.(*domain.CustomContext)

	id := c.Param("id")

	restaurant, isItExists := h.restaurantUsecase.GetById(cc, id)
	if !isItExists {
		return cc.SendERR(errors.BadRequest("error with request data"))
	}

	return cc.SendOK(restaurant)
}
