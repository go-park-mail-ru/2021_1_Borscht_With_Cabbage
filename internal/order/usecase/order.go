package usecase

import (
	"context"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/order"
	"github.com/borscht/backend/internal/restaurantAdmin"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
)

type orderUsecase struct {
	orderRepository      order.OrderRepo
	restaurantRepository restaurantAdmin.AdminRestaurantRepo
}

func NewOrderUsecase(repoOredr order.OrderRepo,
	repoRestaurant restaurantAdmin.AdminRestaurantRepo) order.OrderUsecase {

	return &orderUsecase{
		orderRepository:      repoOredr,
		restaurantRepository: repoRestaurant,
	}
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

func (o orderUsecase) SetNewStatus(ctx context.Context, newStatus models.SetNewStatus) error {
	restaurant, ok := ctx.Value("Restaurant").(models.RestaurantInfo)
	if !ok {
		failError := errors.FailServerError("failed to convert to models.Restaurant")
		logger.UsecaseLevel().ErrorLog(ctx, failError)
		return failError
	}

	newStatus.Restaurant = restaurant.Title
	return o.orderRepository.SetNewStatus(ctx, newStatus)
}

func (o orderUsecase) CreateReview(ctx context.Context, newReview models.SetNewReview) error {
	user, ok := ctx.Value("User").(models.User)
	if !ok {
		failError := errors.FailServerError("failed to convert to models.Restaurant")
		logger.UsecaseLevel().ErrorLog(ctx, failError)
		return failError
	}

	newReview.User = user.Uid
	return o.orderRepository.CreateReview(ctx, newReview)
}
