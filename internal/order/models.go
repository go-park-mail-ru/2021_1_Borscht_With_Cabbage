package order

import (
	"context"
	"github.com/borscht/backend/internal/models"
	"github.com/labstack/echo/v4"
)

type OrderHandler interface {
	Create(c echo.Context) error
	GetUserOrders(c echo.Context) error
	GetRestaurantOrders(c echo.Context) error
	AddToBasket(c echo.Context) error
}

type OrderUsecase interface {
	Create(ctx context.Context, uid int, orderParams models.CreateOrder) error
	GetUserOrders(ctx context.Context, uid int) ([]models.Order, error)
	GetRestaurantOrders(ctx context.Context, restaurantName string) ([]models.Order, error)
	AddToBasket(ctx context.Context, dish models.DishToBasket, uid int) error
}

type OrderRepo interface {
	Create(ctx context.Context, uid int, orderParams models.CreateOrder) error
	GetUserOrders(ctx context.Context, uid int) ([]models.Order, error)
	GetRestaurantOrders(ctx context.Context, restaurantName string) ([]models.Order, error)
	AddToBasket(ctx context.Context, dish models.DishToBasket, uid int) error
}
