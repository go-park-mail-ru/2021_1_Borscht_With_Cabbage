package restaurantAdmin

import (
	"context"

	"github.com/borscht/backend/internal/models"
	"github.com/labstack/echo/v4"
)

type AdminHandler interface {
	Create(c echo.Context) error
	Update(c echo.Context) error
	Login(c echo.Context) error
	GetUserData(c echo.Context) error
	EditProfile(c echo.Context) error
	AddDish(c echo.Context) error
	DeleteDish(c echo.Context) error
}

type AdminUsecase interface {
	Create(ctx context.Context, restaurant models.Restaurant) (int, error)
	Update(ctx context.Context, restaurant models.RestaurantUpdate) (models.RestaurantResponse, error)
	CheckRestaurantExists(ctx context.Context, user models.RestaurantAuth) (models.Restaurant, error)
	GetByRid(ctx context.Context, rid int) (models.Restaurant, error)
	AddDish(ctx context.Context, dish models.Dish) (*models.DishResponse, error)
	DeleteDish(ctx context.Context, did int) error
}

type AdminRepo interface {
	Create(ctx context.Context, user models.Restaurant) (int, error)
	Update(ctx context.Context, restaurant models.RestaurantUpdate) error
	CheckRestaurantExists(ctx context.Context, user models.RestaurantAuth) (models.Restaurant, error)
	GetByRid(ctx context.Context, rid int) (models.Restaurant, error)
	AddDish(ctx context.Context, dish models.Dish) (int, error)
	DeleteDish(ctx context.Context, did int) error
	GetDish(ctx context.Context, did int) (models.Dish, error)
}
