package restaurant

import (
	"context"

	"github.com/borscht/backend/internal/models"
	"github.com/labstack/echo/v4"
)

type GetVendorParams struct {
	Limit     int
	Offset    int
	Address   bool
	Name      string
	Longitude string
	Latitude  string
}

type RestaurantHandler interface {
	GetRestaurantPage(c echo.Context) error
	GetVendor(c echo.Context) error
	GetReviews(c echo.Context) error
}

type RestaurantUsecase interface {
	GetVendor(ctx context.Context, params GetVendorParams) ([]models.RestaurantInfo, error)
	GetById(ctx context.Context, id int) (*models.RestaurantWithDishes, error)
	GetReviews(ctx context.Context, id int) ([]models.RestaurantReview, error)
}

type RestaurantRepo interface {
	GetVendor(ctx context.Context, params GetVendorParams) ([]models.RestaurantInfo, error)
	GetById(ctx context.Context, id int) (*models.RestaurantWithDishes, error)
	GetReviews(ctx context.Context, id int) ([]models.RestaurantReview, error)
}
