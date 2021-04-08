package restaurantAdmin

import (
	"context"
	"mime/multipart"

	"github.com/borscht/backend/internal/models"
	"github.com/labstack/echo/v4"
)

type AdminRestaurantHandler interface {
	Create(c echo.Context) error
	Update(c echo.Context) error
	Login(c echo.Context) error
	GetUserData(c echo.Context) error
	EditProfile(c echo.Context) error
}

type AdminDishHandler interface {
	AddDish(c echo.Context) error
	DeleteDish(c echo.Context) error
	UpdateDish(c echo.Context) error
	GetAllDishes(c echo.Context) error
	UploadDishImage(c echo.Context) error
}

type AdminSectionHandler interface {
	AddSection(c echo.Context) error
	DeleteSection(c echo.Context) error
	UpdateSection(c echo.Context) error
}

type AdminRestaurantUsecase interface {
	Create(ctx context.Context, restaurant models.Restaurant) (*models.Restaurant, error)
	Update(ctx context.Context, restaurant models.RestaurantUpdate) (*models.RestaurantResponse, error)
	CheckRestaurantExists(ctx context.Context, user models.RestaurantAuth) (*models.Restaurant, error)
	GetByRid(ctx context.Context, rid int) (models.Restaurant, error)
}

type AdminDishUsecase interface {
	AddDish(ctx context.Context, dish models.Dish) (*models.Dish, error)
	DeleteDish(ctx context.Context, did int) error
	UpdateDish(ctx context.Context, dish models.Dish) (*models.Dish, error)
	GetAllDishes(ctx context.Context) ([]models.Dish, error)
	UploadDishImage(ctx context.Context, image *multipart.FileHeader, idDish int) (*models.DishImageResponse, error)
}

type AdminSectionUsecase interface {
	AddSection(ctx context.Context, section models.Section) (*models.Section, error)
	UpdateSection(ctx context.Context, section models.Section) (*models.Section, error)
	DeleteSection(ctx context.Context, sid int) error
}

type AdminRestaurantRepo interface {
	Create(ctx context.Context, user models.Restaurant) (int, error)
	Update(ctx context.Context, restaurant models.RestaurantUpdate) error
	CheckRestaurantExists(ctx context.Context, user models.RestaurantAuth) (*models.Restaurant, error)
	GetByRid(ctx context.Context, rid int) (models.Restaurant, error)
}

type AdminDishRepo interface {
	AddDish(ctx context.Context, dish models.Dish) (int, error)
	DeleteDish(ctx context.Context, did int) error
	GetDish(ctx context.Context, did int) (*models.Dish, error)
	UpdateDish(ctx context.Context, dish models.Dish) error
	GetAllDishes(ctx context.Context, idRestaurant int) ([]models.Dish, error)
	UpdateDishImage(ctx context.Context, idDish int, filename string) error
}

type AdminSectionRepo interface {
	AddSection(ctx context.Context, section models.Section) (int, error)
	UpdateSection(ctx context.Context, section models.Section) error
	DeleteSection(ctx context.Context, sid int) error
	GetSection(ctx context.Context, sid int) (*models.Section, error)
}
