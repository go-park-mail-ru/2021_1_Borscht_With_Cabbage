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

//func (o orderUsecase) AddToBasket(ctx context.Context, dish models.DishToBasket, uid int) error {
//	return o.orderRepository.AddToBasket(ctx, dish, uid)
//}
//
//func (o orderUsecase) DeleteFromBasket(ctx context.Context, dish models.DishToBasket, uid int) error {
//	return o.orderRepository.DeleteFromBasket(ctx, dish, uid)
//}

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

//func (o orderUsecase) GetBasket(ctx context.Context, uid int) (*models.BasketForUser, error) {
//	restaurant, err := o.orderRepository.GetBasket(ctx, uid)
//	if err != nil {
//		return nil, err
//	}
//	if restaurant == nil {
//		return restaurant, nil
//	}
//
//	address, err := o.restaurantRepository.GetAddress(ctx, restaurant.RID)
//	if err != nil {
//		return nil, err
//	}
//	if address != nil {
//		logger.UsecaseLevel().DebugLog(ctx, logger.Fields{"address": *address})
//		restaurant.Address = *address
//	}
//	return restaurant, nil
//}

//func (o orderUsecase) AddBasket(ctx context.Context, basket models.BasketForUser) (*models.BasketForUser, error) {
//	user, ok := ctx.Value("User").(models.User)
//	if !ok {
//		failError := errors.FailServerError("failed to convert to models.Restaurant")
//		logger.UsecaseLevel().ErrorLog(ctx, failError)
//		return nil, failError
//	}
//
//	// пока что удаляем предыдущую корзину, в будущем надо будет изменить логику
//	basketOld, err := o.orderRepository.GetBasket(ctx, user.Uid)
//	if err != nil {
//		return nil, err
//	}
//
//	// у пользователя уже есть корзина, удаляем ее
//	if basketOld != nil {
//		err = o.orderRepository.DeleteBasket(ctx, user.Uid, basketOld.BID)
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	newBasketId, err := o.orderRepository.AddBasket(ctx, user.Uid, basket.RID)
//	if err != nil {
//		return nil, err
//	}
//
//	// TODO: попробовать сделать одной транзакцией это
//	// а то может записать только половину корзины
//	for _, value := range basket.Foods {
//		err = o.orderRepository.AddDishToBasket(ctx, newBasketId, value)
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	return o.orderRepository.GetBasket(ctx, user.Uid)
//}
