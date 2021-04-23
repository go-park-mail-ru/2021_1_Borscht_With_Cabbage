package restaurant

import (
	"context"

	"github.com/borscht/backend/internal/models"
	"github.com/labstack/echo/v4"
)

type RestaurantHandler interface {
	GetRestaurantPage(c echo.Context) error
	GetVendor(c echo.Context) error
}

type RestaurantUsecase interface {
	GetVendor(ctx context.Context, request models.RestaurantRequest) ([]models.RestaurantInfo, error)
	GetById(ctx context.Context, id int) (*models.RestaurantWithDishes, error)
}

type RestaurantRepo interface {
	GetVendorWithCategory(ctx context.Context, request models.RestaurantRequest) ([]models.RestaurantInfo, error)
	GetVendor(ctx context.Context, request models.RestaurantRequest) ([]models.RestaurantInfo, error)
	GetById(ctx context.Context, id int) (*models.RestaurantWithDishes, error)
}
