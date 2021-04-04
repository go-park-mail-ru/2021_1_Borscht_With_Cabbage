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

	if err := c.Bind(dish); err != nil {
		sendErr := errors.AuthorizationError("error with request data")
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}

	err := h.OrderUcase.AddToBasket(ctx, dish, userStruct.Uid)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	return models.SendResponse(c, "")
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
	//if err := c.Bind(order); err != nil {
	//	sendErr := errors.AuthorizationError("error with request data")
	//	logger.DeliveryLevel().ErrorLog(ctx, sendErr)
	//	return models.SendResponseWithError(c, sendErr)
	//}

	order.Address = "улица Пупкина"

	userStruct := user.(models.User)
	err := h.OrderUcase.Create(ctx, userStruct.Uid, order)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	return models.SendResponse(c, "")
}

func (h Handler) GetUserOrders(c echo.Context) error {
	//orders := make([]models.Order, 0)
	//dishes := make([]models.DishInOrder, 0)
	//dish := models.DishInOrder{
	//	Name:   "Солянка",
	//	Price:  200,
	//	Number: 2,
	//}
	//dishes = append(dishes, dish)
	//dishes = append(dishes, dish)
	//
	//testOrder1 := models.Order{
	//	OID:          1,
	//	Restaurant:   "rest1",
	//	Address:      "Проспект мира 15,56",
	//	OrderTime:    "12.10.2021 15:45",
	//	DeliveryCost: 200,
	//	DeliveryTime: "1 час",
	//	Summary:      "1900",
	//	Status:       models.StatusOrderAdded,
	//	Foods:        dishes,
	//}
	//
	//testOrder2 := models.Order{
	//	OID:          2,
	//	Restaurant:   "rest2",
	//	Address:      "Бауманская 2",
	//	OrderTime:    "1.01.2021 15:45",
	//	DeliveryCost: 100,
	//	DeliveryTime: "1,5 часа",
	//	Summary:      "100",
	//	Status:       models.StatusOrderAdded,
	//	Foods:        dishes,
	//}
	//orders = append(orders, testOrder1)
	//orders = append(orders, testOrder2)
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

	restaurantStruct := restaurant.(models.Restaurant)
	orders, err := h.OrderUcase.GetRestaurantOrders(ctx, restaurantStruct.Name)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	return models.SendResponse(c, orders)
}
