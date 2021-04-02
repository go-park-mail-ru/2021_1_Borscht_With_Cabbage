package restaurantAdmin

import (
	"context"

	"github.com/borscht/backend/internal/models"
	"github.com/labstack/echo/v4"
)

type AdminHandler interface {
	Create(c echo.Context) error
	Login(c echo.Context) error
	GetUserData(c echo.Context) error
	EditProfile(c echo.Context) error
}

type AdminUsecase interface {
	Create(ctx context.Context, user models.Restaurant) (int, error)
	CheckRestaurantExists(ctx context.Context, user models.RestaurantAuth) (models.Restaurant, error)
	GetByRid(ctx context.Context, rid int) (models.Restaurant, error)
}

type AdminRepo interface {
	Create(ctx context.Context, user models.Restaurant) (int, error)
	CheckRestaurantExists(ctx context.Context, user models.RestaurantAuth) (models.Restaurant, error)
	GetByRid(ctx context.Context, rid int) (models.Restaurant, error)
}
