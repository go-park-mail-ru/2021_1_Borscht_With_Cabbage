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

func convertBasketToProto(basket models.BasketForUser) *protoBasket.BasketInfo {
	dishes := make([]*protoBasket.DishFromBasket, 0)
	for _, dish := range basket.Foods {
		dishInfo := protoBasket.DishFromBasket{
			ID:     int32(dish.ID),
			Name:   dish.Name,
			Price:  int32(dish.Price),
			Number: int32(dish.Number),
			Image:  dish.Image,
		}
		dishes = append(dishes, &dishInfo)
	}
	address := protoBasket.Address{
		Latitude:    basket.Address.Latitude,
		Longitude:   basket.Address.Longitude,
		Radius:      int32(basket.Address.Radius),
		AddressName: basket.Address.Name,
	}
	basketInfo := protoBasket.BasketInfo{
		Bid:             int32(basket.BID),
		RestaurantName:  basket.Restaurant,
		RestaurantImage: basket.RestaurantImage,
		Rid:             int32(basket.RID),
		DeliveryCost:    int32(basket.DeliveryCost),
		Summary:         int32(basket.Summary),
		Address:         &address,
		Dishes:          dishes,
	}
	return &basketInfo
}

func convertProtoToBasket(basket *protoBasket.BasketInfo) models.BasketForUser {
	dishes := make([]models.DishInBasket, 0)
	for _, dish := range basket.Dishes {
		dishInfo := models.DishInBasket{
			ID:     int(dish.ID),
			Name:   dish.Name,
			Price:  int(dish.Price),
			Number: int(dish.Number),
			Image:  dish.Image,
		}
		dishes = append(dishes, dishInfo)
	}
	address := models.Address{
		Latitude:  basket.Address.Latitude,
		Longitude: basket.Address.Longitude,
		Radius:    int(basket.Address.Radius),
		Name:      basket.Address.AddressName,
	}
	basketInfo := models.BasketForUser{
		BID:             int(basket.Bid),
		Restaurant:      basket.RestaurantName,
		RestaurantImage: basket.RestaurantImage,
		RID:             int(basket.Rid),
		DeliveryCost:    int(basket.DeliveryCost),
		Summary:         int(basket.Summary),
		Address:         address,
		Foods:           dishes,
	}
	return basketInfo
}

func (s service) AddToBasket(ctx context.Context, dish models.DishToBasket, uid int) error {
	dishToBasket := protoBasket.DishToBasket{
		SameBasket: dish.SameBasket,
		Did:        int32(dish.DishID),
		IsPlus:     dish.IsPlus,
		Uid:        int32(uid),
	}

	_, err := s.basketService.AddToBasket(ctx, &dishToBasket)
	if err != nil {
		return err
	}

	return nil
}

func (s service) DeleteFromBasket(ctx context.Context, dish models.DishToBasket, uid int) error {
	dishToDelete := protoBasket.DishToBasket{
		SameBasket: dish.SameBasket,
		Did:        int32(dish.DishID),
		IsPlus:     dish.IsPlus,
	}

	_, err := s.basketService.DeleteFromBasket(ctx, &dishToDelete)
	if err != nil {
		return err
	}

	return nil
}

func (s service) GetBasket(ctx context.Context, uid int) (*models.BasketForUser, error) {
	UID := protoBasket.UID{
		Uid: int32(uid),
	}

	basket, err := s.basketService.GetBasket(ctx, &UID)
	if err != nil {
		return nil, err
	}

	basketForUser := models.BasketForUser{
		BID:             int(basket.Bid),
		Restaurant:      basket.RestaurantName,
		RestaurantImage: basket.RestaurantImage,
		RID:             int(basket.Rid),
		DeliveryCost:    int(basket.DeliveryCost),
		Summary:         int(basket.Summary),
	}

	dishes := make([]models.DishInBasket, 0)
	for _, dish := range basket.Dishes {
		dishInfo := models.DishInBasket{
			ID:     int(dish.ID),
			Name:   dish.Name,
			Price:  int(dish.Price),
			Number: int(dish.Number),
			Image:  dish.Image,
		}
		dishes = append(dishes, dishInfo)
	}
	basketForUser.Foods = dishes

	return &basketForUser, nil
}

func (s service) AddBasket(ctx context.Context, basket models.BasketForUser) (*models.BasketForUser, error) {
	basketInfo := convertBasketToProto(basket)

	basketResult, err := s.basketService.AddBasket(ctx, basketInfo)
	if err != nil {
		return nil, err
	}

	basketResponse := convertProtoToBasket(basketResult)
	return &basketResponse, nil
}
