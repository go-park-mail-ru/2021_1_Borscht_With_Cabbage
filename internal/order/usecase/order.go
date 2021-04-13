package usecase

import (
	"context"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/order"
)

type orderUsecase struct {
	orderRepository order.OrderRepo
}

func NewOrderUsecase(repo order.OrderRepo) order.OrderUsecase {
	return &orderUsecase{
		orderRepository: repo,
	}
}

func (o orderUsecase) AddToBasket(ctx context.Context, dish models.DishToBasket, uid int) error {
	return o.orderRepository.AddToBasket(ctx, dish, uid)
}

func (o orderUsecase) DeleteFromBasket(ctx context.Context, dish models.DishToBasket, uid int) error {
	return o.orderRepository.DeleteFromBasket(ctx, dish, uid)
}

func (o orderUsecase) Create(ctx context.Context, uid int, orderParams models.CreateOrder) error {
	return o.orderRepository.Create(ctx, uid, orderParams)
}

func (o orderUsecase) GetUserOrders(ctx context.Context, uid int) ([]models.Order, error) {
	return o.orderRepository.GetUserOrders(ctx, uid)
}

func (o orderUsecase) GetRestaurantOrders(ctx context.Context, restaurantName string) ([]models.Order, error) {
	return o.orderRepository.GetRestaurantOrders(ctx, restaurantName)
}

func (o orderUsecase) GetBasket(ctx context.Context, uid int) (models.BasketForUser, error) {
	return o.orderRepository.GetBasket(ctx, uid)
}
