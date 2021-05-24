package http

import (
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/order"
	"github.com/borscht/backend/internal/services/basket"
	errors "github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
	"github.com/borscht/backend/utils/notifications"
	"github.com/labstack/echo/v4"
	"strconv"
)

type Handler struct {
	OrderUcase    order.OrderUsecase
	BasketService basket.ServiceBasket
	Notificator   notifications.OrderNotificator
}

func NewOrderHandler(orderUcase order.OrderUsecase, basketService basket.ServiceBasket, notificator notifications.OrderNotificator) order.OrderHandler {
	handler := &Handler{
		OrderUcase:    orderUcase,
		BasketService: basketService,
		Notificator:   notificator,
	}

	return handler
}

func (h Handler) AddBaskets(c echo.Context) error {
	ctx := models.GetContext(c)
	user, ok := ctx.Value("User").(models.User)
	if !ok {
		failError := errors.FailServerError("failed to get user")
		logger.UsecaseLevel().ErrorLog(ctx, failError)
		return models.SendResponseWithError(c, failError)
	}

	baskets := models.BasketsForUser{}
	if err := c.Bind(&baskets); err != nil {
		sendErr := errors.BadRequestError("error with request data")
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}
	for i := range baskets.Baskets {
		baskets.Baskets[i].UID = user.Uid
	}
	logger.DeliveryLevel().DebugLog(ctx, logger.Fields{"baskets": baskets.Baskets})

	result, err := h.BasketService.AddBaskets(ctx, baskets.Baskets)
	if err != nil {
		models.SendResponseWithError(c, err)
	}

	if result == nil {
		result = &[]models.BasketForUser{}
	}
	return models.SendResponse(c, &models.BasketsForUser{
		Baskets: *result,
	})
}

func (h Handler) AddToBasket(c echo.Context) error {
	ctx := models.GetContext(c)
	var err error

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
	dish.RestaurantID, err = strconv.Atoi(c.Param("restID"))
	if err != nil {
		sendErr := errors.AuthorizationError("error with restaurant id")
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}

	if dish.IsPlus {
		err := h.BasketService.AddToBasket(ctx, dish, userStruct.Uid)
		if err != nil {
			return models.SendResponseWithError(c, err)
		}

		basket, err := h.BasketService.GetBasket(ctx, userStruct.Uid, dish.RestaurantID)
		if err != nil {
			return models.SendResponseWithError(c, err)
		}

		logger.DeliveryLevel().InlineDebugLog(ctx, basket)
		return models.SendResponse(c, basket)
	}

	err = h.BasketService.DeleteFromBasket(ctx, dish, userStruct.Uid)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	basket, err := h.BasketService.GetBasket(ctx, userStruct.Uid, dish.RestaurantID)
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
	var err error

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

	order.BasketID, err = strconv.Atoi(c.Param("basketID"))
	if err != nil {
		sendErr := errors.AuthorizationError("error with restaurant id")
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}

	userStruct := user.(models.User)
	err = h.OrderUcase.Create(ctx, userStruct.Uid, order)
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

	uid, err := h.OrderUcase.SetNewStatus(ctx, newStatus)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	err = h.Notificator.OrderStatusChangedNotification(newStatus.Status, uid)
	if err != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Order notification error")
	}

	return models.SendResponse(c, nil)
}

func (h Handler) CreateReview(c echo.Context) error {
	ctx := models.GetContext(c)

	newReview := models.SetNewReview{}
	if err := c.Bind(&newReview); err != nil {
		sendErr := errors.AuthorizationError("error with request data")
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}

	err := h.OrderUcase.CreateReview(ctx, newReview)
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

	rid, err := strconv.Atoi(c.Param("restID"))
	if err != nil {
		sendErr := errors.AuthorizationError("error with restaurant id")
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}

	basket, err := h.BasketService.GetBasket(ctx, userStruct.Uid, rid)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	if basket == nil {
		basket = &models.BasketForUser{}
	}
	logger.DeliveryLevel().InfoLog(ctx, logger.Fields{"basket": basket})
	return models.SendResponse(c, basket)
}

func (h Handler) GetBaskets(c echo.Context) error {
	ctx := models.GetContext(c)

	user := c.Get("User")
	if user == nil {
		sendErr := errors.AuthorizationError("error with request data")
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}
	userStruct := user.(models.User)

	getBasketsParams := models.GetBasketParams{
		Uid:       userStruct.Uid,
		Latitude:  c.QueryParam("latitude"),
		Longitude: c.QueryParam("longitude"),
	}

	baskets, err := h.BasketService.GetBaskets(ctx, getBasketsParams)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	response := models.BasketsForUser{}
	if baskets != nil {
		response = models.BasketsForUser{
			Baskets: baskets,
		}
	}

	logger.DeliveryLevel().InfoLog(ctx, logger.Fields{"basket": baskets})
	return models.SendResponse(c, &response)
}

func (h Handler) DeleteBaskets(c echo.Context) error {
	ctx := models.GetContext(c)

	user := c.Get("User")
	if user == nil {
		sendErr := errors.AuthorizationError("error with request data")
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}
	userStruct := user.(models.User)

	err := h.BasketService.DeleteBaskets(ctx, userStruct.Uid)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	return models.SendResponse(c, nil)
}

func (h Handler) DeleteBasket(c echo.Context) error {
	ctx := models.GetContext(c)

	bid, err := strconv.Atoi(c.Param("basketID"))
	if err != nil {
		sendErr := errors.AuthorizationError("error with basket id")
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}

	err = h.BasketService.DeleteBasket(ctx, bid)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	return models.SendResponse(c, nil)
}
