package basket

import (
	"context"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/services/basket/mocks"
	proto "github.com/borscht/backend/services/proto/basket"
	"github.com/borscht/backend/utils/logger"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestService_AddBasket(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	clientMock := mocks.NewMockBasketClient(ctrl)
	basketService := NewService(clientMock)

	c := context.Background()
	uid := 1
	dish := models.DishToBasket{
		SameBasket: true,
		DishID:     1,
		IsPlus:     true,
	}
	dishRequest := proto.DishToBasket{
		SameBasket: dish.SameBasket,
		Uid:        int32(uid),
		Did:        int32(dish.DishID),
		IsPlus:     dish.IsPlus,
	}

	clientMock.EXPECT().AddToBasket(c, &dishRequest).Return(&proto.Nothing{}, nil)

	err := basketService.AddToBasket(c, dish, uid)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
}

func TestService_DeleteFromBasket(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	clientMock := mocks.NewMockBasketClient(ctrl)
	basketService := NewService(clientMock)

	c := context.Background()
	uid := 1
	dish := models.DishToBasket{
		SameBasket: true,
		DishID:     1,
		IsPlus:     true,
	}
	dishRequest := proto.DishToBasket{
		SameBasket: dish.SameBasket,
		Uid:        int32(uid),
		Did:        int32(dish.DishID),
		IsPlus:     dish.IsPlus,
	}

	clientMock.EXPECT().DeleteFromBasket(c, &dishRequest).Return(&proto.Nothing{}, nil)

	err := basketService.DeleteFromBasket(c, dish, uid)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
}

func TestService_GetBasket(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	clientMock := mocks.NewMockBasketClient(ctrl)
	basketService := NewService(clientMock)

	c := context.Background()
	logger.InitLogger()

	uid := 1
	uidRequest := proto.UID{
		Uid: int32(uid),
	}
	basketInfo := proto.BasketInfo{
		Bid:            1,
		RestaurantName: "title",
		Rid:            1,
		Address: &proto.Address{
			AddressName: "Бауманская 2",
			Radius:      1500,
		},
	}

	clientMock.EXPECT().GetBasket(c, &uidRequest).Return(&basketInfo, nil)

	basket, err := basketService.GetBasket(c, uid)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	require.EqualValues(t, basket.RID, 1)
}

func TestService_AddToBasket(t *testing.T) {
	//ctrl := gomock.NewController(t)
	//defer ctrl.Finish()
	//clientMock := mocks.NewMockBasketClient(ctrl)
	//basketService := NewService(clientMock)
	//
	//c := context.Background()
	//logger.InitLogger()
	//
	//basketToSave := models.BasketForUser{
	//	RID:        1,
	//	UID:        1,
	//	Restaurant: "title",
	//	Address: models.Address{
	//		Name:   "Бауманская 2",
	//		Radius: 1500,
	//	},
	//}
	//basketInfo := proto.BasketInfo{
	//	Uid:            1,
	//	RestaurantName: "title",
	//	Rid:            1,
	//	Address: &proto.Address{
	//		AddressName: "Бауманская 2",
	//		Radius:      1500,
	//	},
	//}
	//result := basketInfo
	//result.Bid = 1
	//
	//clientMock.EXPECT().AddBasket(c, &basketInfo).Return(&result, nil)
	//
	//basket, err := basketService.AddBasket(c, basketToSave)
	//if err != nil {
	//	t.Errorf("unexpected err: %s", err)
	//	return
	//}
	//
	//require.EqualValues(t, basket.RID, 1)
}
