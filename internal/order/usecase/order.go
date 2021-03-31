package usecase

import (
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

func (o orderUsecase) Create(uid int) (string, error) {
	panic("implement me")
}

func (o orderUsecase) GetUserOrder(session string) {
	// get user
	panic("implement me")
}
