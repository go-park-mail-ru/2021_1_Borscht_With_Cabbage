package http

import (
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/order/mocks"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewOrderHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	OrderUsecaseMock := mocks.NewMockOrderUsecase(ctrl)
	orderHandler := NewOrderHandler(OrderUsecaseMock)
	if orderHandler == nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestHandler_AddToBasket_SameRestaurant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	OrderUsecaseMock := mocks.NewMockOrderUsecase(ctrl)
	orderHandler := NewOrderHandler(OrderUsecaseMock)

	dish := models.DishToBasket{
		DishID:     1,
		IsPlus:     true,
		SameBasket: true,
	}
	inputJSON := `{"dishid":1,"isPlus":true,"same":true}`

	foods := []models.DishInBasket{
		{1, "food1", 200, 1, "img.jpg"},
	}
	basket := models.BasketForUser{
		BID:          1,
		Restaurant:   "first",
		RID:          1,
		DeliveryCost: 200,
		Summary:      1200,
		Foods:        foods,
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/user/basket", strings.NewReader(inputJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	user := models.User{
		Uid:  1,
		Name: "Daria",
	}
	c.Set("User", user)
	ctx := models.GetContext(c)

	OrderUsecaseMock.EXPECT().AddToBasket(ctx, dish, user.Uid).Return(nil)
	OrderUsecaseMock.EXPECT().GetBasket(ctx, user.Uid).Return(basket, nil)

	err := orderHandler.AddToBasket(c)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestHandler_AddToBasket_NewRestaurant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	OrderUsecaseMock := mocks.NewMockOrderUsecase(ctrl)
	orderHandler := NewOrderHandler(OrderUsecaseMock)

	dish := models.DishToBasket{
		DishID:     1,
		IsPlus:     true,
		SameBasket: false,
	}
	inputJSON := `{"dishid":1,"isPlus":true,"same":false}`

	foods := []models.DishInBasket{
		{1, "food1", 200, 1, "img.jpg"},
	}
	basket := models.BasketForUser{
		BID:          1,
		Restaurant:   "first",
		RID:          1,
		DeliveryCost: 200,
		Summary:      1200,
		Foods:        foods,
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/user/basket", strings.NewReader(inputJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	user := models.User{
		Uid:  1,
		Name: "Daria",
	}
	c.Set("User", user)
	ctx := models.GetContext(c)

	OrderUsecaseMock.EXPECT().AddToBasket(ctx, dish, user.Uid).Return(nil)
	OrderUsecaseMock.EXPECT().GetBasket(ctx, user.Uid).Return(basket, nil)

	err := orderHandler.AddToBasket(c)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestHandler_DeleteFromBasket_SameRestaurant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	OrderUsecaseMock := mocks.NewMockOrderUsecase(ctrl)
	orderHandler := NewOrderHandler(OrderUsecaseMock)

	dish := models.DishToBasket{
		DishID:     1,
		IsPlus:     false,
		SameBasket: true,
	}
	inputJSON := `{"dishid":1,"isPlus":false,"same":true}`

	foods := []models.DishInBasket{
		{1, "food1", 200, 1, "img.jpg"},
	}
	basket := models.BasketForUser{
		BID:          1,
		Restaurant:   "first",
		RID:          1,
		DeliveryCost: 200,
		Summary:      1200,
		Foods:        foods,
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/user/basket", strings.NewReader(inputJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	user := models.User{
		Uid:  1,
		Name: "Daria",
	}
	c.Set("User", user)
	ctx := models.GetContext(c)

	OrderUsecaseMock.EXPECT().DeleteFromBasket(ctx, dish, user.Uid).Return(nil)
	OrderUsecaseMock.EXPECT().GetBasket(ctx, user.Uid).Return(basket, nil)

	err := orderHandler.AddToBasket(c)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	OrderUsecaseMock := mocks.NewMockOrderUsecase(ctrl)
	orderHandler := NewOrderHandler(OrderUsecaseMock)

	order := models.CreateOrder{
		Address: "prospekt mira 134",
	}
	inputJSON := `{"address":"prospekt mira 134"}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/user/order", strings.NewReader(inputJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	user := models.User{
		Uid:  1,
		Name: "Daria",
	}
	c.Set("User", user)
	ctx := models.GetContext(c)

	OrderUsecaseMock.EXPECT().Create(ctx, 1, order).Return(nil)

	err := orderHandler.Create(c)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestHandler_GetUserOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	OrderUsecaseMock := mocks.NewMockOrderUsecase(ctrl)
	orderHandler := NewOrderHandler(OrderUsecaseMock)

	orders := []models.Order{
		{OID: 1, UID: 1, Restaurant: "rest1", Address: "prospekt 1"},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/user/orders", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	user := models.User{
		Uid:  1,
		Name: "Daria",
	}
	c.Set("User", user)
	ctx := models.GetContext(c)

	OrderUsecaseMock.EXPECT().GetUserOrders(ctx, user.Uid).Return(orders, nil)

	err := orderHandler.GetUserOrders(c)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestHandler_GetRestaurantOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	OrderUsecaseMock := mocks.NewMockOrderUsecase(ctrl)
	orderHandler := NewOrderHandler(OrderUsecaseMock)

	orders := []models.Order{
		{OID: 1, UID: 1, Restaurant: "rest1", Address: "prospekt 1"},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/restaurant/orders", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	restaurant := models.RestaurantInfo{
		Title: "firstRest",
	}
	c.Set("Restaurant", restaurant)
	ctx := models.GetContext(c)

	OrderUsecaseMock.EXPECT().GetRestaurantOrders(ctx, restaurant.Title).Return(orders, nil)

	err := orderHandler.GetRestaurantOrders(c)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestHandler_GetBasket(t *testing.T) {

}
