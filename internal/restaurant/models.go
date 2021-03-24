package restaurant

import (
	"github.com/borscht/backend/internal/models"
	"github.com/labstack/echo/v4"
)

type RestaurantHandler interface {
	GetRestaurantPage(c echo.Context) error
	GetVendor(c echo.Context) error
}

type RestaurantUsecase interface {
	GetVendor(limit, offset int) ([]models.RestaurantResponse, error)
	GetById(id string) (models.Restaurant, bool, error)
}

type RestaurantRepo interface {
	GetVendor(limit, offset int) ([]models.RestaurantResponse, error)
	GetById(id string) (models.Restaurant, bool, error)
}
