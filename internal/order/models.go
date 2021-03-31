package order

import (
	"github.com/borscht/backend/internal/models"
	"github.com/labstack/echo/v4"
)

type OrderHandler interface {
	Create(c echo.Context) error
	GetUserOrders(c echo.Context) error
	GetRestaurantOrders(c echo.Context) error
}

type OrderUsecase interface {
	Create(uid int, orderParams models.CreateOrder) error
	GetUserOrders(uid int) ([]models.Order, error)
	GetRestaurantOrders(restaurantName string) ([]models.Order, error)
}

type OrderRepo interface {
	Create(uid int, orderParams models.CreateOrder) error
	GetUserOrders(uid int) ([]models.Order, error)
	GetRestaurantOrders(restaurantName string) ([]models.Order, error)
}
