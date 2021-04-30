package internal

import (
	"context"
	"github.com/borscht/backend/internal/models"
	basketServiceRepo "github.com/borscht/backend/services/basket/repository"
	protoBasket "github.com/borscht/backend/services/proto/basket"
	"github.com/borscht/backend/utils/logger"
)

type service struct {
	basketRepo basketServiceRepo.BasketRepo
}

func NewService(basketRepo basketServiceRepo.BasketRepo) *service {
	return &service{
		basketRepo: basketRepo,
	}
}

func (s service) AddToBasket(ctx context.Context, dish *protoBasket.DishToBasket) (*protoBasket.Nothing, error) {
	dishToBasket := models.DishToBasket{
		SameBasket: dish.SameBasket,
		IsPlus:     dish.IsPlus,
		DishID:     int(dish.Did),
	}

	err := s.basketRepo.AddToBasket(ctx, dishToBasket, int(dish.Uid))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s service) DeleteFromBasket(ctx context.Context, dish *protoBasket.DishToBasket) (*protoBasket.Nothing, error) {
	dishFromBasket := models.DishToBasket{
		SameBasket: dish.SameBasket,
		IsPlus:     dish.IsPlus,
		DishID:     int(dish.Did),
	}

	err := s.basketRepo.DeleteFromBasket(ctx, dishFromBasket, int(dish.Uid))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s service) GetBasket(ctx context.Context, uid *protoBasket.UID) (*protoBasket.BasketInfo, error) {
	basket, err := s.basketRepo.GetBasket(ctx, int(uid.Uid))
	if err != nil {
		return nil, err
	}
	if basket == nil {
		return nil, nil
	}

	address, err := s.basketRepo.GetAddress(ctx, basket.RID)
	if err != nil {
		return nil, err
	}
	if address != nil {
		logger.UsecaseLevel().DebugLog(ctx, logger.Fields{"address": *address})
		basket.Address = *address
	}

	dishesProto := make([]*protoBasket.DishFromBasket, 0)
	for _, dish := range basket.Foods {
		dishProto := protoBasket.DishFromBasket{
			ID:     int32(dish.ID),
			Name:   dish.Name,
			Price:  int32(dish.Price),
			Number: int32(dish.Number),
			Image:  dish.Image,
		}
		dishesProto = append(dishesProto, &dishProto)
	}
	addressProto := protoBasket.Address{
		Longitude:   basket.Address.Longitude,
		Latitude:    basket.Address.Latitude,
		AddressName: basket.Address.Name,
		Radius:      int32(basket.Address.Radius),
	}
	basketResponse := protoBasket.BasketInfo{
		Bid:             int32(basket.BID),
		RestaurantName:  basket.Restaurant,
		RestaurantImage: basket.RestaurantImage,
		Rid:             int32(basket.DeliveryCost),
		Summary:         int32(basket.Summary),
		Address:         &addressProto,
		Dishes:          dishesProto,
	}

	return &basketResponse, nil
}

func (s service) AddBasket(ctx context.Context, info *protoBasket.BasketInfo) (*protoBasket.BasketInfo, error) {
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
}
