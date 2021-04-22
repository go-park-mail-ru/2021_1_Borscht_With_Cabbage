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
	SetNewStatus(c echo.Context) error
	AddToBasket(c echo.Context) error
	GetBasket(c echo.Context) error
	AddBasket(c echo.Context) error
}

type OrderUsecase interface {
	Create(ctx context.Context, uid int, orderParams models.CreateOrder) error
	GetUserOrders(ctx context.Context, uid int) ([]models.Order, error)
	GetRestaurantOrders(ctx context.Context, restaurantName string) ([]models.Order, error)
	SetNewStatus(ctx context.Context, newStatus models.SetNewStatus) error
	AddToBasket(ctx context.Context, dish models.DishToBasket, uid int) error
	DeleteFromBasket(ctx context.Context, dish models.DishToBasket, uid int) error
	GetBasket(ctx context.Context, uid int) (*models.BasketForUser, error)
	AddBasket(ctx context.Context, basket models.BasketForUser) (*models.BasketForUser, error)
}

type OrderRepo interface {
	Create(ctx context.Context, uid int, orderParams models.CreateOrder) error
	GetUserOrders(ctx context.Context, uid int) ([]models.Order, error)
	GetRestaurantOrders(ctx context.Context, restaurantName string) ([]models.Order, error)
	SetNewStatus(ctx context.Context, newStatus models.SetNewStatus) error
	AddToBasket(ctx context.Context, dish models.DishToBasket, uid int) error
	DeleteFromBasket(ctx context.Context, dish models.DishToBasket, uid int) error
	GetBasket(ctx context.Context, uid int) (*models.BasketForUser, error)
	AddBasket(ctx context.Context, userID, restaurantID int) (basketID int, err error)
	DeleteBasket(ctx context.Context, userID, basketID int) error
	AddDishToBasket(ctx context.Context, basketID int, dish models.DishInBasket) error
}
