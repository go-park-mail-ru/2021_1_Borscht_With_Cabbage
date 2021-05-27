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

func createBasketResponse(basket models.BasketForUser, address models.Address) protoBasket.BasketInfo {
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
		Longitude:   float32(basket.Address.Longitude),
		Latitude:    float32(basket.Address.Latitude),
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

	return basketResponse
}

func (s service) AddToBasket(ctx context.Context, dish *protoBasket.DishToBasket) (*protoBasket.Nothing, error) {
	dishToBasket := models.DishToBasket{
		SameBasket: dish.SameBasket,
		IsPlus:     dish.IsPlus,
		DishID:     int(dish.Did),
	}

	err := s.basketRepo.AddToBasket(ctx, dishToBasket, int(dish.Uid))
	if err != nil {
		return &protoBasket.Nothing{}, err
	}

	return &protoBasket.Nothing{}, nil
}

func (s service) DeleteFromBasket(ctx context.Context, dish *protoBasket.DishToBasket) (*protoBasket.Nothing, error) {
	dishFromBasket := models.DishToBasket{
		SameBasket: dish.SameBasket,
		IsPlus:     dish.IsPlus,
		DishID:     int(dish.Did),
	}

	err := s.basketRepo.DeleteFromBasket(ctx, dishFromBasket, int(dish.Uid))
	if err != nil {
		return &protoBasket.Nothing{}, err
	}

	return &protoBasket.Nothing{}, nil
}

func (s service) GetBasket(ctx context.Context, uid *protoBasket.UID) (*protoBasket.BasketInfo, error) {
	basket, err := s.basketRepo.GetBasket(ctx, int(uid.Uid))
	if err != nil {
		return &protoBasket.BasketInfo{}, err
	}
	if basket == nil {
		return &protoBasket.BasketInfo{}, nil
	}

	address, errr := s.basketRepo.GetAddress(ctx, basket.RID)
	if errr != nil {
		return &protoBasket.BasketInfo{}, err
	}
	if address != nil {
		logger.UsecaseLevel().DebugLog(ctx, logger.Fields{"address": *address})
		basket.Address = *address
	}

	basketResponse := createBasketResponse(*basket, *address)

	return &basketResponse, nil
}

func (s service) AddBasket(ctx context.Context, info *protoBasket.BasketInfo) (*protoBasket.BasketInfo, error) {
	uid := int(info.Uid)
	// пока что удаляем предыдущую корзину, в будущем надо будет изменить логику
	basketOld, err := s.basketRepo.GetBasket(ctx, uid)
	if err != nil {
		return &protoBasket.BasketInfo{}, err
	}

	// у пользователя уже есть корзина, удаляем ее
	if basketOld != nil {
		err = s.basketRepo.DeleteBasket(ctx, uid, basketOld.BID)
		if err != nil {
			return &protoBasket.BasketInfo{}, err
		}
	}

	newBasketId, err := s.basketRepo.AddBasket(ctx, uid, int(info.Rid))
	if err != nil {
		return &protoBasket.BasketInfo{}, err
	}

	// TODO: попробовать сделать одной транзакцией это
	// а то может записать только половину корзины
	for _, value := range info.Dishes {
		err = s.basketRepo.AddDishToBasket(ctx, newBasketId, models.DishInBasket{
			ID:     int(value.ID),
			Name:   value.Name,
			Number: int(value.Number),
			Image:  value.Image,
			Price:  int(value.Price),
		})
		if err != nil {
			return &protoBasket.BasketInfo{}, err
		}
	}

	basketResult, err := s.basketRepo.GetBasket(ctx, uid)
	if err != nil {
		return &protoBasket.BasketInfo{}, err
	}

	basketResponse := createBasketResponse(*basketResult, models.Address{})

	return &basketResponse, nil
}
