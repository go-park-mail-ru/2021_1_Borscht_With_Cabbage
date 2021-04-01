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
	orders := make([]models.Order, 0)
	dishes := make([]models.Dish, 0)
	dish := models.Dish{
		Name:  "Солянка",
		Price: 200,
	}
	dishes = append(dishes, dish)

	testOrder1 := models.Order{
		OID:          1,
		Restaurant:   "rest1",
		Address:      "Проспект мира 15,56",
		OrderTime:    "12.10.2021 15:45",
		DeliveryCost: 200,
		DeliveryTime: "1 час",
		Summary:      "1900",
		Status:       models.StatusOrderAdded,
		Foods:        dishes,
	}
	orders = append(orders, testOrder1)
	orders = append(orders, testOrder1)

	//user := c.Get("User")
	//if user == nil {
	//	userError := errors.Authorization("not authorized")
	//	return models.SendResponseWithError(c, userError)
	//}
	//
	//userStruct := user.(models.User)
	//orders, err := h.OrderUcase.GetUserOrders(userStruct.Uid)
	//if err != nil {
	//	return models.SendResponseWithError(c, err)
	//}

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
