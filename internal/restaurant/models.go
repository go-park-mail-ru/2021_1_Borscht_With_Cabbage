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
	GetRecommendations(c echo.Context) error
}

type RestaurantUsecase interface {
	GetVendor(ctx context.Context, request models.RestaurantRequest) ([]models.RestaurantInfo, error)
	GetById(ctx context.Context, coordinates models.Coordinates, id int) (*models.RestaurantWithDishes, error)
	GetReviews(ctx context.Context, id int) ([]models.RestaurantReview, error)
	GetRecommendations(ctx context.Context, params models.RecommendationsParams) ([]models.RestaurantInfo, error)
}

type RestaurantRepo interface {
	GetAllCategories(ctx context.Context) ([]string, error)
	GetVendor(ctx context.Context, request models.RestaurantRequest) ([]models.RestaurantInfo, error)
	GetById(ctx context.Context, id int, userCoordinates models.Coordinates) (*models.RestaurantWithDishes, error)
	GetReviews(ctx context.Context, id int) ([]models.RestaurantReview, error)
	GetRecommendations(ctx context.Context, params models.RecommendationsParams) ([]models.RestaurantInfo, error)
}
