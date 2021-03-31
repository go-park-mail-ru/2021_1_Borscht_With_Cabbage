package http

import (
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/order"
	errors "github.com/borscht/backend/utils"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	OrderUcase order.OrderUsecase
}

func NewOrderHandler(orderUcase order.OrderUsecase) order.OrderHandler {
	handler := &Handler{
		OrderUcase: orderUcase,
	}

	return handler
}

func (h Handler) Create(c echo.Context) error {
	user := c.Get("User")
	if user == nil {
		userError := errors.Authorization("not authorized")
		return models.SendResponseWithError(c, userError)
	}

	order := models.CreateOrder{}
	if err := c.Bind(order); err != nil {
		sendErr := errors.Authorization("error with request data")
		return models.SendResponseWithError(c, sendErr)
	}

	userStruct := user.(models.User)
	err := h.OrderUcase.Create(userStruct.Uid, order)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	return models.SendResponse(c, "")
}

func (h Handler) GetUserOrders(c echo.Context) error {
	user := c.Get("User")
	if user == nil {
		userError := errors.Authorization("not authorized")
		return models.SendResponseWithError(c, userError)
	}

	userStruct := user.(models.User)
	orders, err := h.OrderUcase.GetUserOrders(userStruct.Uid)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	return models.SendResponse(c, orders)
}

func (h Handler) GetRestaurantOrders(c echo.Context) error {
	restaurant := c.Get("Restaurant")
	if restaurant == nil {
		userError := errors.Authorization("not authorized")
		return models.SendResponseWithError(c, userError)
	}

	restaurantStruct := restaurant.(models.Restaurant)
	orders, err := h.OrderUcase.GetRestaurantOrders(restaurantStruct.Name)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	return models.SendResponse(c, orders)
}
