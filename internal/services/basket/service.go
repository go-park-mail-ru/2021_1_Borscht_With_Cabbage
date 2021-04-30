package basket

import (
	"context"
	"github.com/borscht/backend/internal/models"
	protoBasket "github.com/borscht/backend/services/proto/basket"
)

type BasketAuth interface {
	AddToBasket(ctx context.Context, dish models.DishToBasket, uid int) error
	DeleteFromBasket(ctx context.Context, dish models.DishToBasket, uid int) error
	GetBasket(ctx context.Context, uid int) (*models.BasketForUser, error)
	AddBasket(ctx context.Context, basket models.BasketForUser) (*models.BasketForUser, error)
}

type service struct {
	basketService protoBasket.BasketClient
}

func NewService(basketService protoBasket.BasketClient) BasketAuth {
	return &service{
		basketService: basketService,
	}
}

func (s service) AddToBasket(ctx context.Context, dish models.DishToBasket, uid int) error {
	panic("implement me")
}

func (s service) DeleteFromBasket(ctx context.Context, dish models.DishToBasket, uid int) error {
	panic("implement me")
}

func (s service) GetBasket(ctx context.Context, uid int) (*models.BasketForUser, error) {
	panic("implement me")
}

func (s service) AddBasket(ctx context.Context, basket models.BasketForUser) (*models.BasketForUser, error) {
	panic("implement me")
}
