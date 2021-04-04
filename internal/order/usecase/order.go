package usecase

import (
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

func (o orderUsecase) AddToBasket(dish models.DishToBasket, uid int) error {
	return o.orderRepository.AddToBasket(dish, uid)
}

func (o orderUsecase) Create(uid int, orderParams models.CreateOrder) error {
	return o.orderRepository.Create(uid, orderParams)
}

func (o orderUsecase) GetUserOrders(uid int) ([]models.Order, error) {
	return o.orderRepository.GetUserOrders(uid)
}

func (o orderUsecase) GetRestaurantOrders(restaurantName string) ([]models.Order, error) {
	return o.orderRepository.GetRestaurantOrders(restaurantName)
}
