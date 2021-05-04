package restaurant

import (
	"context"

	"github.com/borscht/backend/internal/models"
	"github.com/labstack/echo/v4"
)

const (
	MinutesInHour = 60
	CourierSpeed  = 5
	CookingTime   = 30
)

type RestaurantHandler interface {
	GetRestaurantPage(c echo.Context) error
	GetVendor(c echo.Context) error
	GetReviews(c echo.Context) error
}

type RestaurantUsecase interface {
	GetVendor(ctx context.Context, request models.RestaurantRequest) ([]models.RestaurantInfo, error)
	GetById(ctx context.Context, id int) (*models.RestaurantWithDishes, error)
	GetReviews(ctx context.Context, id int) ([]models.RestaurantReview, error)
}

type RestaurantRepo interface {
	GetVendorWithCategory(ctx context.Context, request models.RestaurantRequest) ([]models.RestaurantInfo, error)
	GetVendor(ctx context.Context, request models.RestaurantRequest) ([]models.RestaurantInfo, error)
	GetById(ctx context.Context, id int) (*models.RestaurantWithDishes, error)
	GetReviews(ctx context.Context, id int) ([]models.RestaurantReview, error)
}
