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
		Rid:             int32(basket.RID),
		RestaurantName:  basket.Restaurant,
		RestaurantImage: basket.RestaurantImage,
		DeliveryCost:    int32(basket.DeliveryCost),
		Summary:         int32(basket.Summary),
		Address:         &addressProto,
		Dishes:          dishesProto,
	}

	return basketResponse
}

func (s service) AddToBasket(ctx context.Context, dish *protoBasket.DishToBasket) (*protoBasket.Nothing, error) {
	dishToBasket := models.DishToBasket{
		RestaurantID: int(dish.Rid),
		IsPlus:       dish.IsPlus,
		DishID:       int(dish.Did),
	}

	err := s.basketRepo.AddToBasket(ctx, dishToBasket, int(dish.Uid))
	if err != nil {
		return &protoBasket.Nothing{}, err
	}

	return &protoBasket.Nothing{}, nil
}

func (s service) DeleteFromBasket(ctx context.Context, dish *protoBasket.DishToBasket) (*protoBasket.Nothing, error) {
	dishFromBasket := models.DishToBasket{
		RestaurantID: int(dish.Rid),
		IsPlus:       dish.IsPlus,
		DishID:       int(dish.Did),
	}

	err := s.basketRepo.DeleteFromBasket(ctx, dishFromBasket, int(dish.Uid))
	if err != nil {
		return &protoBasket.Nothing{}, err
	}

	return &protoBasket.Nothing{}, nil
}

func (s service) GetBasket(ctx context.Context, uid *protoBasket.IDs) (*protoBasket.BasketInfo, error) {
	basket, err := s.basketRepo.GetBasket(ctx, int(uid.Uid), int(uid.Rid))
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

func (s service) GetBaskets(ctx context.Context, params *protoBasket.GetBasketsParams) (*protoBasket.Baskets, error) {
	paramsUser := models.GetBasketParams{
		Uid:       int(params.Uid),
		Latitude:  float64(params.Latitude),
		Longitude: float64(params.Longitude),
	}

	baskets, err := s.basketRepo.GetBaskets(ctx, paramsUser)
	if err != nil {
		return &protoBasket.Baskets{}, err
	}

	basketsResponse := make([]*protoBasket.BasketInfo, 0)
	for _, basket := range baskets.Baskets {
		basketResponse := createBasketResponse(basket, models.Address{})
		basketsResponse = append(basketsResponse, &basketResponse)
	}

	return &protoBasket.Baskets{Baskets: basketsResponse}, nil
}

func (s service) AddBaskets(ctx context.Context, info *protoBasket.Baskets) (*protoBasket.Baskets, error) {
	uid := int(info.Baskets[0].Uid)

	response := protoBasket.Baskets{}
	for _, basket := range info.Baskets {
		// пока что удаляем предыдущую корзину, в будущем надо будет изменить логику
		basketOld, err := s.basketRepo.GetBasket(ctx, uid, int(basket.Rid))
		if err != nil {
			return &protoBasket.Baskets{}, err
		}

		// у пользователя уже есть корзина, удаляем ее
		if basketOld != nil {
			err = s.basketRepo.DeleteBasket(ctx, basketOld.BID)
			if err != nil {
				return &protoBasket.Baskets{}, err
			}
		}

		newBasketId, err := s.basketRepo.AddBasket(ctx, uid, int(basket.Rid))
		if err != nil {
			return &protoBasket.Baskets{}, err
		}

		// TODO: попробовать сделать одной транзакцией это
		// а то может записать только половину корзины
		for _, value := range basket.Dishes {
			err = s.basketRepo.AddDishToBasket(ctx, newBasketId, models.DishInBasket{
				ID:     int(value.ID),
				Name:   value.Name,
				Number: int(value.Number),
				Image:  value.Image,
				Price:  int(value.Price),
			})
			if err != nil {
				return &protoBasket.Baskets{}, err
			}
		}

		basketResult, err := s.basketRepo.GetBasket(ctx, uid, int(basket.Rid))
		if err != nil {
			return &protoBasket.Baskets{}, err
		}

		basketResponse := createBasketResponse(*basketResult, models.Address{})
		response.Baskets = append(response.Baskets, &basketResponse)
	}

	return &response, nil
}

func (s service) DeleteBaskets(ctx context.Context, uid *protoBasket.UID) (*protoBasket.Nothing, error) {
	err := s.basketRepo.DeleteBaskets(ctx, int(uid.Uid))
	if err != nil {
		return &protoBasket.Nothing{}, err
	}

	return &protoBasket.Nothing{}, nil
}

func (s service) DeleteBasket(ctx context.Context, bid *protoBasket.BID) (*protoBasket.Nothing, error) {
	err := s.basketRepo.DeleteBasket(ctx, int(bid.Bid))
	if err != nil {
		return &protoBasket.Nothing{}, err
	}

	return &protoBasket.Nothing{}, nil
}
