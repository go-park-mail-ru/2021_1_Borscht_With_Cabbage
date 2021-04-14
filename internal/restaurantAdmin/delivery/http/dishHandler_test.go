package http

import (
	"github.com/borscht/backend/internal/models"
	adminMock "github.com/borscht/backend/internal/restaurantAdmin/mocks"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewDishHandlerew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	DishUsecaseMock := adminMock.NewMockAdminDishUsecase(ctrl)
	dishHandler := NewDishHandler(DishUsecaseMock)
	if dishHandler == nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestDishHandler_GetAllDishes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	DishUsecaseMock := adminMock.NewMockAdminDishUsecase(ctrl)
	dishHandler := NewDishHandler(DishUsecaseMock)

	dishes := []models.SectionWithDishes{
		{SectionId: 1},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/restaurant/dishes", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	ctx := models.GetContext(c)

	DishUsecaseMock.EXPECT().GetAllDishes(ctx).Return(dishes, nil)

	err := dishHandler.GetAllDishes(c)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestDishHandler_UpdateDishData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	DishUsecaseMock := adminMock.NewMockAdminDishUsecase(ctrl)
	dishHandler := NewDishHandler(DishUsecaseMock)

	dishes := models.Dish{
		Name: "dish1",
	}
	updateDish := models.Dish{
		Name: "dish1",
	}
	requestJSON := `{"name":"dish1"}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/restaurant/dish", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	ctx := models.GetContext(c)

	DishUsecaseMock.EXPECT().UpdateDishData(ctx, updateDish).Return(&dishes, nil)

	err := dishHandler.UpdateDishData(c)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestDishHandler_DeleteDish(t *testing.T) {
	//ctrl := gomock.NewController(t)
	//defer ctrl.Finish()
	//DishUsecaseMock := adminMock.NewMockAdminDishUsecase(ctrl)
	//dishHandler := NewDishHandler(DishUsecaseMock)
	//
	//deleteSuccess := models.DeleteSuccess{
	//	ID: 1,
	//}
	//updateDish := models.DishDelete{
	//	ID: 1,
	//}
	//requestJSON := `{"id":1}`
	//
	//e := echo.New()
	//req := httptest.NewRequest(http.MethodDelete, "/restaurant/dish", strings.NewReader(requestJSON))
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	//rec := httptest.NewRecorder()
	//c := e.NewContext(req, rec)
	//ctx := models.GetContext(c)
	//
	//DishUsecaseMock.EXPECT().DeleteDish(ctx, updateDish).Return(&deleteSuccess, nil)
	//
	//err := dishHandler.DeleteDish(c)
	//if err != nil {
	//	t.Errorf("incorrect result")
	//	return
	//}
}

func TestDishHandler_AddDish(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	DishUsecaseMock := adminMock.NewMockAdminDishUsecase(ctrl)
	dishHandler := NewDishHandler(DishUsecaseMock)

	response := models.Dish{
		ID:     1,
		Name:   "dish1",
		Price:  200,
		Weight: 120,
	}
	newDish := models.Dish{
		Name:   "dish1",
		Price:  200,
		Weight: 120,
	}
	requestJSON := `{"name":"dish1","price":200,"weight":120}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/restaurant/dish", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	ctx := models.GetContext(c)

	DishUsecaseMock.EXPECT().AddDish(ctx, newDish).Return(&response, nil)

	err := dishHandler.AddDish(c)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestDishHandler_UploadDishImage(t *testing.T) {
	// TODO
}
