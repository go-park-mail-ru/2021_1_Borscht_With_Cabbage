package http

import (
	"fmt"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/order"
	errors "github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
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

func (h Handler) AddToBasket(c echo.Context) error {
	ctx := models.GetContext(c)

	user := c.Get("User")
	if user == nil {
		sendErr := errors.AuthorizationError("error with request data")
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}
	userStruct := user.(models.User)

	dish := models.DishToBasket{}
	if err := c.Bind(&dish); err != nil {
		sendErr := errors.AuthorizationError("error with request data")
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}

	if dish.IsPlus {
		err := h.OrderUcase.AddToBasket(ctx, dish, userStruct.Uid)
		if err != nil {
			return models.SendResponseWithError(c, err)
		}

		basket, err := h.OrderUcase.GetBasket(ctx, userStruct.Uid)
		if err != nil {
			return models.SendResponseWithError(c, err)
		}
		fmt.Println(basket)
		return models.SendResponse(c, basket)
	}

	err := h.OrderUcase.DeleteFromBasket(ctx, dish, userStruct.Uid)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	basket, err := h.OrderUcase.GetBasket(ctx, userStruct.Uid)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	return models.SendResponse(c, basket)
}

func (h Handler) Create(c echo.Context) error {
	ctx := models.GetContext(c)

	user := c.Get("User")
	if user == nil {
		sendErr := errors.AuthorizationError("error with request data")
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}

	order := models.CreateOrder{}
	if err := c.Bind(&order); err != nil {
		sendErr := errors.AuthorizationError("error with request data")
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}

	userStruct := user.(models.User)
	err := h.OrderUcase.Create(ctx, userStruct.Uid, order)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	return models.SendResponse(c, "")
}

func (h Handler) GetUserOrders(c echo.Context) error {
	ctx := models.GetContext(c)

	user := c.Get("User")
	if user == nil {
		sendErr := errors.AuthorizationError("error with request data")
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}

	userStruct := user.(models.User)
	orders, err := h.OrderUcase.GetUserOrders(ctx, userStruct.Uid)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	return models.SendResponse(c, orders)
}

func (h Handler) GetRestaurantOrders(c echo.Context) error {
	ctx := models.GetContext(c)

	restaurant := c.Get("Restaurant")
	if restaurant == nil {
		sendErr := errors.AuthorizationError("error with request data")
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}

	restaurantStruct := restaurant.(models.RestaurantInfo)
	orders, err := h.OrderUcase.GetRestaurantOrders(ctx, restaurantStruct.Title)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	return models.SendResponse(c, orders)
}

func (h Handler) GetBasket(c echo.Context) error {
	ctx := models.GetContext(c)

	user := c.Get("User")
	if user == nil {
		sendErr := errors.AuthorizationError("error with request data")
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}

	userStruct := user.(models.User)
	basket, err := h.OrderUcase.GetBasket(ctx, userStruct.Uid)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	return models.SendResponse(c, basket)
}
