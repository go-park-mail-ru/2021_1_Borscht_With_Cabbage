package basketServiceRepo

import (
	"context"
	"database/sql"
	"github.com/borscht/backend/internal/models"
)

type BasketRepo interface {
	AddToBasket(ctx context.Context, dish models.DishToBasket, uid int) error
	DeleteFromBasket(ctx context.Context, dish models.DishToBasket, uid int) error
	GetBasket(ctx context.Context, uid int) (*models.BasketForUser, error)
	AddBasket(ctx context.Context, userID, restaurantID int) (basketID int, err error)
	DeleteBasket(ctx context.Context, userID, basketID int) error
	AddDishToBasket(ctx context.Context, basketID int, dish models.DishInBasket) error
}

type basketRepository struct {
	DB *sql.DB
}

func NewBasketRepository(db *sql.DB) BasketRepo {
	return &basketRepository{
		DB: db,
	}
}

func (b basketRepository) AddToBasket(ctx context.Context, dish models.DishToBasket, uid int) error {
	panic("implement me")
}

func (b basketRepository) DeleteFromBasket(ctx context.Context, dish models.DishToBasket, uid int) error {
	panic("implement me")
}

func (b basketRepository) GetBasket(ctx context.Context, uid int) (*models.BasketForUser, error) {
	panic("implement me")
}

func (b basketRepository) AddBasket(ctx context.Context, userID, restaurantID int) (basketID int, err error) {
	panic("implement me")
}

func (b basketRepository) DeleteBasket(ctx context.Context, userID, basketID int) error {
	panic("implement me")
}

func (b basketRepository) AddDishToBasket(ctx context.Context, basketID int, dish models.DishInBasket) error {
	panic("implement me")
}
