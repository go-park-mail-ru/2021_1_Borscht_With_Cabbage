package http

import (
	"encoding/json"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/order/mocks"
	basketServiceMock "github.com/borscht/backend/internal/services/mocks"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/notifications"
	"github.com/borscht/backend/utils/websocketPool"
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
	BasketService := basketServiceMock.NewMockServiceBasket(ctrl)
	websocketConnectionsUsers := websocketPool.NewConnectionPool()
	websocketConnectionsRestaurants := websocketPool.NewConnectionPool()
	notificator := notifications.NewOrderNotificator(&websocketConnectionsUsers, &websocketConnectionsRestaurants)
	orderHandler := NewOrderHandler(OrderUsecaseMock, BasketService, notificator)
	if orderHandler == nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestHandler_AddBasket(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	OrderUsecaseMock := mocks.NewMockOrderUsecase(ctrl)
	BasketServiceMock := basketServiceMock.NewMockServiceBasket(ctrl)
	websocketConnectionsUsers := websocketPool.NewConnectionPool()
	websocketConnectionsRestaurants := websocketPool.NewConnectionPool()
	notificator := notifications.NewOrderNotificator(&websocketConnectionsUsers, &websocketConnectionsRestaurants)
	orderHandler := NewOrderHandler(OrderUsecaseMock, BasketServiceMock, notificator)

	basket := models.BasketForUser{
		Restaurant: "rest1",
		RID:        1,
		UID:        1,
	}
	inputJSON := `{"restaurantName":"rest1","restaurantID":1}`

	basketResult := models.BasketForUser{
		BID:        1,
		Restaurant: "first",
		RID:        1,
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/user/basket", strings.NewReader(inputJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	user := models.User{
		Uid:  1,
		Name: "Daria",
	}
	c.Set("User", user)
	ctx := models.GetContext(c)

	BasketServiceMock.EXPECT().AddBasket(ctx, basket).Return(&basketResult, nil)

	err := orderHandler.AddBasket(c)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestHandler_AddBasket_BindError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	OrderUsecaseMock := mocks.NewMockOrderUsecase(ctrl)
	BasketServiceMock := basketServiceMock.NewMockServiceBasket(ctrl)
	websocketConnectionsUsers := websocketPool.NewConnectionPool()
	websocketConnectionsRestaurants := websocketPool.NewConnectionPool()
	notificator := notifications.NewOrderNotificator(&websocketConnectionsUsers, &websocketConnectionsRestaurants)
	orderHandler := NewOrderHandler(OrderUsecaseMock, BasketServiceMock, notificator)

	inputJSON := `{"restaurantNamerest1","restaurantID":1}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/user/basket", strings.NewReader(inputJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	user := models.User{
		Uid:  1,
		Name: "Daria",
	}
	c.Set("User", user)

	err := orderHandler.AddBasket(c)
	b := errors.SendError{}
	respCode := rec.Body.Bytes()
	err = json.Unmarshal(respCode, &b)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
	if b.Code == 200 {
		t.Errorf("incorrect result")
		return
	}
}

func TestHandler_AddToBasket_SameRestaurant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	OrderUsecaseMock := mocks.NewMockOrderUsecase(ctrl)
	BasketServiceMock := basketServiceMock.NewMockServiceBasket(ctrl)
	websocketConnectionsUsers := websocketPool.NewConnectionPool()
	websocketConnectionsRestaurants := websocketPool.NewConnectionPool()
	notificator := notifications.NewOrderNotificator(&websocketConnectionsUsers, &websocketConnectionsRestaurants)
	orderHandler := NewOrderHandler(OrderUsecaseMock, BasketServiceMock, notificator)

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

	BasketServiceMock.EXPECT().AddToBasket(ctx, dish, user.Uid).Return(nil)
	BasketServiceMock.EXPECT().GetBasket(ctx, user.Uid).Return(&basket, nil)

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
	BasketServiceMock := basketServiceMock.NewMockServiceBasket(ctrl)
	websocketConnectionsUsers := websocketPool.NewConnectionPool()
	websocketConnectionsRestaurants := websocketPool.NewConnectionPool()
	notificator := notifications.NewOrderNotificator(&websocketConnectionsUsers, &websocketConnectionsRestaurants)
	orderHandler := NewOrderHandler(OrderUsecaseMock, BasketServiceMock, notificator)

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

	BasketServiceMock.EXPECT().AddToBasket(ctx, dish, user.Uid).Return(nil)
	BasketServiceMock.EXPECT().GetBasket(ctx, user.Uid).Return(&basket, nil)

	err := orderHandler.AddToBasket(c)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestHandler_AddToBasket_BindError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	OrderUsecaseMock := mocks.NewMockOrderUsecase(ctrl)
	BasketServiceMock := basketServiceMock.NewMockServiceBasket(ctrl)
	websocketConnectionsUsers := websocketPool.NewConnectionPool()
	websocketConnectionsRestaurants := websocketPool.NewConnectionPool()
	notificator := notifications.NewOrderNotificator(&websocketConnectionsUsers, &websocketConnectionsRestaurants)
	orderHandler := NewOrderHandler(OrderUsecaseMock, BasketServiceMock, notificator)

	inputJSON := `{"dishid:1,"isPlus":true,"same":true}`

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

	err := orderHandler.AddToBasket(c)
	b := errors.SendError{}
	respCode := rec.Body.Bytes()
	err = json.Unmarshal(respCode, &b)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
	if b.Code == 200 {
		t.Errorf("incorrect result")
		return
	}
}

func TestHandler_AddToBasket_UserNilError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	OrderUsecaseMock := mocks.NewMockOrderUsecase(ctrl)
	BasketServiceMock := basketServiceMock.NewMockServiceBasket(ctrl)
	websocketConnectionsUsers := websocketPool.NewConnectionPool()
	websocketConnectionsRestaurants := websocketPool.NewConnectionPool()
	notificator := notifications.NewOrderNotificator(&websocketConnectionsUsers, &websocketConnectionsRestaurants)
	orderHandler := NewOrderHandler(OrderUsecaseMock, BasketServiceMock, notificator)

	inputJSON := `{"dishid":1,"isPlus":true,"same":true}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/user/basket", strings.NewReader(inputJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := orderHandler.AddToBasket(c)
	b := errors.SendError{}
	respCode := rec.Body.Bytes()
	err = json.Unmarshal(respCode, &b)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
	if b.Code == 200 {
		t.Errorf("incorrect result")
		return
	}
}

func TestHandler_DeleteFromBasket_SameRestaurant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	OrderUsecaseMock := mocks.NewMockOrderUsecase(ctrl)
	BasketServiceMock := basketServiceMock.NewMockServiceBasket(ctrl)
	websocketConnectionsUsers := websocketPool.NewConnectionPool()
	websocketConnectionsRestaurants := websocketPool.NewConnectionPool()
	notificator := notifications.NewOrderNotificator(&websocketConnectionsUsers, &websocketConnectionsRestaurants)
	orderHandler := NewOrderHandler(OrderUsecaseMock, BasketServiceMock, notificator)

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

	BasketServiceMock.EXPECT().DeleteFromBasket(ctx, dish, user.Uid).Return(nil)
	BasketServiceMock.EXPECT().GetBasket(ctx, user.Uid).Return(&basket, nil)

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
	BasketServiceMock := basketServiceMock.NewMockServiceBasket(ctrl)
	websocketConnectionsUsers := websocketPool.NewConnectionPool()
	websocketConnectionsRestaurants := websocketPool.NewConnectionPool()
	notificator := notifications.NewOrderNotificator(&websocketConnectionsUsers, &websocketConnectionsRestaurants)
	orderHandler := NewOrderHandler(OrderUsecaseMock, BasketServiceMock, notificator)

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

func TestHandler_Create_UserNilError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	OrderUsecaseMock := mocks.NewMockOrderUsecase(ctrl)
	BasketServiceMock := basketServiceMock.NewMockServiceBasket(ctrl)
	websocketConnectionsUsers := websocketPool.NewConnectionPool()
	websocketConnectionsRestaurants := websocketPool.NewConnectionPool()
	notificator := notifications.NewOrderNotificator(&websocketConnectionsUsers, &websocketConnectionsRestaurants)
	orderHandler := NewOrderHandler(OrderUsecaseMock, BasketServiceMock, notificator)

	inputJSON := `{"address":"prospekt mira 134"}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/user/order", strings.NewReader(inputJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := orderHandler.Create(c)
	b := errors.SendError{}
	respCode := rec.Body.Bytes()
	err = json.Unmarshal(respCode, &b)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
	if b.Code == 200 {
		t.Errorf("incorrect result")
		return
	}
}

func TestHandler_Create_BindError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	OrderUsecaseMock := mocks.NewMockOrderUsecase(ctrl)
	BasketServiceMock := basketServiceMock.NewMockServiceBasket(ctrl)
	websocketConnectionsUsers := websocketPool.NewConnectionPool()
	websocketConnectionsRestaurants := websocketPool.NewConnectionPool()
	notificator := notifications.NewOrderNotificator(&websocketConnectionsUsers, &websocketConnectionsRestaurants)
	orderHandler := NewOrderHandler(OrderUsecaseMock, BasketServiceMock, notificator)

	inputJSON := `{"address"prospekt mira 134"}`

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

	err := orderHandler.Create(c)
	b := errors.SendError{}
	respCode := rec.Body.Bytes()
	err = json.Unmarshal(respCode, &b)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
	if b.Code == 200 {
		t.Errorf("incorrect result")
		return
	}
}

func TestHandler_GetUserOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	OrderUsecaseMock := mocks.NewMockOrderUsecase(ctrl)
	BasketServiceMock := basketServiceMock.NewMockServiceBasket(ctrl)
	websocketConnectionsUsers := websocketPool.NewConnectionPool()
	websocketConnectionsRestaurants := websocketPool.NewConnectionPool()
	notificator := notifications.NewOrderNotificator(&websocketConnectionsUsers, &websocketConnectionsRestaurants)
	orderHandler := NewOrderHandler(OrderUsecaseMock, BasketServiceMock, notificator)

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

func TestHandler_GetUserOrders_UserNilError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	OrderUsecaseMock := mocks.NewMockOrderUsecase(ctrl)
	BasketServiceMock := basketServiceMock.NewMockServiceBasket(ctrl)
	websocketConnectionsUsers := websocketPool.NewConnectionPool()
	websocketConnectionsRestaurants := websocketPool.NewConnectionPool()
	notificator := notifications.NewOrderNotificator(&websocketConnectionsUsers, &websocketConnectionsRestaurants)
	orderHandler := NewOrderHandler(OrderUsecaseMock, BasketServiceMock, notificator)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/user/orders", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := orderHandler.GetUserOrders(c)
	b := errors.SendError{}
	respCode := rec.Body.Bytes()
	err = json.Unmarshal(respCode, &b)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
	if b.Code == 200 {
		t.Errorf("incorrect result")
		return
	}
}

func TestHandler_GetRestaurantOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	OrderUsecaseMock := mocks.NewMockOrderUsecase(ctrl)
	BasketServiceMock := basketServiceMock.NewMockServiceBasket(ctrl)
	websocketConnectionsUsers := websocketPool.NewConnectionPool()
	websocketConnectionsRestaurants := websocketPool.NewConnectionPool()
	notificator := notifications.NewOrderNotificator(&websocketConnectionsUsers, &websocketConnectionsRestaurants)
	orderHandler := NewOrderHandler(OrderUsecaseMock, BasketServiceMock, notificator)

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

func TestHandler_GetRestaurantOrders_RestaurantNilError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	OrderUsecaseMock := mocks.NewMockOrderUsecase(ctrl)
	BasketServiceMock := basketServiceMock.NewMockServiceBasket(ctrl)
	websocketConnectionsUsers := websocketPool.NewConnectionPool()
	websocketConnectionsRestaurants := websocketPool.NewConnectionPool()
	notificator := notifications.NewOrderNotificator(&websocketConnectionsUsers, &websocketConnectionsRestaurants)
	orderHandler := NewOrderHandler(OrderUsecaseMock, BasketServiceMock, notificator)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/restaurant/orders", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := orderHandler.GetRestaurantOrders(c)
	b := errors.SendError{}
	respCode := rec.Body.Bytes()
	err = json.Unmarshal(respCode, &b)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
	if b.Code == 200 {
		t.Errorf("incorrect result")
		return
	}
}

func TestHandler_GetBasket(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	OrderUsecaseMock := mocks.NewMockOrderUsecase(ctrl)
	BasketServiceMock := basketServiceMock.NewMockServiceBasket(ctrl)
	websocketConnectionsUsers := websocketPool.NewConnectionPool()
	websocketConnectionsRestaurants := websocketPool.NewConnectionPool()
	notificator := notifications.NewOrderNotificator(&websocketConnectionsUsers, &websocketConnectionsRestaurants)
	orderHandler := NewOrderHandler(OrderUsecaseMock, BasketServiceMock, notificator)

	basket := models.BasketForUser{
		BID:          1,
		Restaurant:   "rest1",
		RID:          1,
		DeliveryCost: 200,
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/restaurant/orders", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	user := models.User{
		Uid:  1,
		Name: "Daria",
	}
	c.Set("User", user)
	ctx := models.GetContext(c)

	BasketServiceMock.EXPECT().GetBasket(ctx, user.Uid).Return(&basket, nil)

	err := orderHandler.GetBasket(c)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestHandler_SetNewStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	OrderUsecaseMock := mocks.NewMockOrderUsecase(ctrl)
	BasketServiceMock := basketServiceMock.NewMockServiceBasket(ctrl)
	websocketConnectionsUsers := websocketPool.NewConnectionPool()
	websocketConnectionsRestaurants := websocketPool.NewConnectionPool()
	notificator := notifications.NewOrderNotificator(&websocketConnectionsUsers, &websocketConnectionsRestaurants)
	orderHandler := NewOrderHandler(OrderUsecaseMock, BasketServiceMock, notificator)

	inputJSON := `{"order":1,"status":"delivering","deliveryTime":"12.01.2001"}`
	newStatus := models.SetNewStatus{
		OID:          1,
		Status:       models.StatusOrderDelivering,
		DeliveryTime: "12.01.2001",
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/user/order/status", strings.NewReader(inputJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	ctx := models.GetContext(c)

	OrderUsecaseMock.EXPECT().SetNewStatus(ctx, newStatus).Return(0, nil)

	err := orderHandler.SetNewStatus(c)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestHandler_SetNewStatus_BindError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	OrderUsecaseMock := mocks.NewMockOrderUsecase(ctrl)
	BasketServiceMock := basketServiceMock.NewMockServiceBasket(ctrl)
	websocketConnectionsUsers := websocketPool.NewConnectionPool()
	websocketConnectionsRestaurants := websocketPool.NewConnectionPool()
	notificator := notifications.NewOrderNotificator(&websocketConnectionsUsers, &websocketConnectionsRestaurants)
	orderHandler := NewOrderHandler(OrderUsecaseMock, BasketServiceMock, notificator)

	inputJSON := `{"order:1,"status":"delivering","deliveryTime":"12.01.2001"}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/user/order/status", strings.NewReader(inputJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := orderHandler.SetNewStatus(c)
	b := errors.SendError{}
	respCode := rec.Body.Bytes()
	err = json.Unmarshal(respCode, &b)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
	if b.Code == 200 {
		t.Errorf("incorrect result")
		return
	}
}

func TestHandler_CreateReview(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	OrderUsecaseMock := mocks.NewMockOrderUsecase(ctrl)
	BasketServiceMock := basketServiceMock.NewMockServiceBasket(ctrl)
	websocketConnectionsUsers := websocketPool.NewConnectionPool()
	websocketConnectionsRestaurants := websocketPool.NewConnectionPool()
	notificator := notifications.NewOrderNotificator(&websocketConnectionsUsers, &websocketConnectionsRestaurants)
	orderHandler := NewOrderHandler(OrderUsecaseMock, BasketServiceMock, notificator)

	inputJSON := `{"oid":1,"review":"cool","stars":5}`
	newReview := models.SetNewReview{
		OID:    1,
		Review: "cool",
		Stars:  5,
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/user/order/review", strings.NewReader(inputJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	ctx := models.GetContext(c)

	OrderUsecaseMock.EXPECT().CreateReview(ctx, newReview).Return(nil)

	err := orderHandler.CreateReview(c)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestHandler_CreateReview_BindError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	OrderUsecaseMock := mocks.NewMockOrderUsecase(ctrl)
	BasketServiceMock := basketServiceMock.NewMockServiceBasket(ctrl)
	websocketConnectionsUsers := websocketPool.NewConnectionPool()
	websocketConnectionsRestaurants := websocketPool.NewConnectionPool()
	notificator := notifications.NewOrderNotificator(&websocketConnectionsUsers, &websocketConnectionsRestaurants)
	orderHandler := NewOrderHandler(OrderUsecaseMock, BasketServiceMock, notificator)

	inputJSON := `{"oid1,"review":"cool","stars":5}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/user/order/review", strings.NewReader(inputJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := orderHandler.CreateReview(c)
	b := errors.SendError{}
	respCode := rec.Body.Bytes()
	err = json.Unmarshal(respCode, &b)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
	if b.Code == 200 {
		t.Errorf("incorrect result")
		return
	}
}
