package internal

import (
	"context"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/services/basket/mocks"
	protoBasket "github.com/borscht/backend/services/proto/basket"
	"github.com/borscht/backend/utils/logger"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestService_AddToBasket(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	basketRepoMock := mocks.NewMockBasketRepo(ctrl)

	basketService := NewService(basketRepoMock)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)
	logger.InitLogger()

	dishProto := protoBasket.DishToBasket{
		SameBasket: true,
		IsPlus:     true,
		Did:        1,
		Uid:        1,
	}
	dishToBasket := models.DishToBasket{
		SameBasket: dishProto.SameBasket,
		IsPlus:     dishProto.IsPlus,
		DishID:     int(dishProto.Did),
	}

	basketRepoMock.EXPECT().AddToBasket(ctx, dishToBasket, int(dishProto.Uid)).Return(nil)

	_, err := basketService.AddToBasket(ctx, &dishProto)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
}

func TestService_DeleteFromBasket(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	basketRepoMock := mocks.NewMockBasketRepo(ctrl)

	basketService := NewService(basketRepoMock)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)
	logger.InitLogger()

	dishProto := protoBasket.DishToBasket{
		SameBasket: true,
		IsPlus:     true,
		Did:        1,
		Uid:        1,
	}
	dishToBasket := models.DishToBasket{
		SameBasket: dishProto.SameBasket,
		IsPlus:     dishProto.IsPlus,
		DishID:     int(dishProto.Did),
	}

	basketRepoMock.EXPECT().DeleteFromBasket(ctx, dishToBasket, int(dishProto.Uid)).Return(nil)

	_, err := basketService.DeleteFromBasket(ctx, &dishProto)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
}

func TestService_GetBasket(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	basketRepoMock := mocks.NewMockBasketRepo(ctrl)

	basketService := NewService(basketRepoMock)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)
	logger.InitLogger()

	basket := models.BasketForUser{
		RID: 1,
		BID: 1,
	}
	address := models.Address{
		Name:      "address1",
		Latitude:  "12343",
		Longitude: "123442",
		Radius:    1000,
	}

	basketRepoMock.EXPECT().GetBasket(ctx, 1).Return(&basket, nil)
	basketRepoMock.EXPECT().GetAddress(ctx, basket.RID).Return(&address, nil)

	basketResponse, err := basketService.GetBasket(ctx, &protoBasket.UID{Uid: 1})
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	require.EqualValues(t, basketResponse.Address.AddressName, address.Name)
}

func TestService_AddBasket(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	basketRepoMock := mocks.NewMockBasketRepo(ctrl)

	basketService := NewService(basketRepoMock)

	c := context.Background()
	user := models.User{
		Email:    "dasha@mail.ru",
		Phone:    "89111111111",
		Name:     "111111",
		Password: "1111111",
		Uid:      1,
	}
	ctx := context.WithValue(c, "User", user)

	basketProto := protoBasket.BasketInfo{
		Uid: 1,
		Rid: 1,
	}
	basketOld := models.BasketForUser{
		BID: 1,
	}

	basketRepoMock.EXPECT().GetBasket(ctx, user.Uid).Return(&basketOld, nil)
	basketRepoMock.EXPECT().DeleteBasket(ctx, user.Uid, basketOld.BID).Return(nil)
	basketRepoMock.EXPECT().AddBasket(ctx, user.Uid, 1)
	basketRepoMock.EXPECT().GetBasket(ctx, user.Uid).Return(&models.BasketForUser{BID: 2}, nil)

	basketResponse, err := basketService.AddBasket(ctx, &basketProto)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	require.EqualValues(t, basketResponse.Bid, 2)
}
