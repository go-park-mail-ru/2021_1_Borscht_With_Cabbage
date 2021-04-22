package http

import (
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

func (h Handler) AddBasket(c echo.Context) error {
	ctx := models.GetContext(c)

	basket := models.BasketForUser{}
	if err := c.Bind(&basket); err != nil {
		sendErr := errors.BadRequestError("error with request data")
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}
	logger.DeliveryLevel().DebugLog(ctx, logger.Fields{"basket": basket})

	result, err := h.OrderUcase.AddBasket(ctx, basket)
	if err != nil {
		models.SendResponseWithError(c, err)
	}

	if result == nil {
		result = &models.BasketForUser{}
	}
	return models.SendResponse(c, result)
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

		logger.DeliveryLevel().InlineDebugLog(ctx, basket)
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

	if basket == nil {
		basket = &models.BasketForUser{}
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

	return models.SendResponse(c, nil)
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

	response := make([]models.Response, 0)
	for i := range orders {
		response = append(response, &orders[i])
	}

	return models.SendMoreResponse(c, response...)
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

	response := make([]models.Response, 0)
	for i := range orders {
		response = append(response, &orders[i])
	}

	return models.SendMoreResponse(c, response...)
}

func (h Handler) SetNewStatus(c echo.Context) error {
	ctx := models.GetContext(c)

	newStatus := models.SetNewStatus{}
	if err := c.Bind(&newStatus); err != nil {
		sendErr := errors.AuthorizationError("error with request data")
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}

	err := h.OrderUcase.SetNewStatus(ctx, newStatus)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	return models.SendResponse(c, nil)
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

	if basket == nil {
		basket = &models.BasketForUser{}
	}
	logger.DeliveryLevel().InfoLog(ctx, logger.Fields{"basket": basket})
	return models.SendResponse(c, basket)
}
