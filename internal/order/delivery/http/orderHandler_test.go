package http

import (
	"testing"
)

func TestNewOrderHandler(t *testing.T) {

}

func TestHandler_AddToBasket(t *testing.T) {

}

func TestHandler_Create(t *testing.T) {
	//ctrl := gomock.NewController(t)
	//defer ctrl.Finish()
	//OrderUsecaseMock := mocks.NewMockOrderUsecase(ctrl)
	//orderHandler := NewOrderHandler(OrderUsecaseMock)
	//
	//order := models.CreateOrder{
	//	Address: "prospekt mira 134",
	//}
	//inputJSON := `{"address":"prospekt mira 134"}`
	//
	//e := echo.New()
	//req := httptest.NewRequest(http.MethodPost, "/user/order", strings.NewReader(inputJSON))
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	//rec := httptest.NewRecorder()
	//c := e.NewContext(req, rec)
	//user := models.User{
	//	Uid:  1,
	//	Name: "Daria",
	//}
	//c.Set("User", user)
	//ctx := models.GetContext(c)
	//
	//OrderUsecaseMock.EXPECT().Create(ctx, 1, order).Return(nil)
	//
	//err := orderHandler.Create(c)
	//if err != nil {
	//	t.Errorf("incorrect result")
	//	return
	//}
}

func TestHandler_GetUserOrders(t *testing.T) {

}

func TestHandler_GetRestaurantOrders(t *testing.T) {

}
