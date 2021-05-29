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
	CreateReview(c echo.Context) error
	AddToBasket(c echo.Context) error
	GetBasket(c echo.Context) error
	GetBaskets(c echo.Context) error
	AddBaskets(c echo.Context) error
	DeleteBaskets(c echo.Context) error
	DeleteBasket(c echo.Context) error
}

type OrderUsecase interface {
	Create(ctx context.Context, uid int, orderParams models.CreateOrder, coordinates models.Coordinates) error
	GetUserOrders(ctx context.Context, uid int) ([]models.Order, error)
	GetRestaurantOrders(ctx context.Context, restaurantName string) ([]models.Order, error)
	SetNewStatus(ctx context.Context, newStatus models.SetNewStatus) (int, error)
	CreateReview(ctx context.Context, newReview models.SetNewReview) error
}

type OrderRepo interface {
	Create(ctx context.Context, uid int, orderParams models.CreateOrder, coordinates models.Coordinates) error
	GetUserOrders(ctx context.Context, uid int) ([]models.Order, error)
	GetRestaurantOrders(ctx context.Context, restaurantName string) ([]models.Order, error)
	SetNewStatus(ctx context.Context, newStatus models.SetNewStatus) (int, error)
	CreateReview(ctx context.Context, newReview models.SetNewReview) error
}
