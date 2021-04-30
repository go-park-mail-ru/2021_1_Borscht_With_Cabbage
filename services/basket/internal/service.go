package internal

import (
	"context"
	basketServiceRepo "github.com/borscht/backend/services/basket/repository"
	protoBasket "github.com/borscht/backend/services/proto/basket"
)

type service struct {
	basketRepo basketServiceRepo.BasketRepo
}

func NewService(basketRepo basketServiceRepo.BasketRepo) *service {
	return &service{
		basketRepo: basketRepo,
	}
}

func (s service) AddToBasket(ctx context.Context, basket *protoBasket.DishToBasket) (*protoBasket.Nothing, error) {
	panic("implement me")
}

func (s service) DeleteFromBasket(ctx context.Context, basket *protoBasket.DishToBasket) (*protoBasket.Nothing, error) {
	panic("implement me")
}

func (s service) GetBasket(ctx context.Context, uid *protoBasket.UID) (*protoBasket.BasketInfo, error) {
	panic("implement me")
}

func (s service) AddBasket(ctx context.Context, info *protoBasket.BasketInfo) (*protoBasket.BasketInfo, error) {
	panic("implement me")
}
