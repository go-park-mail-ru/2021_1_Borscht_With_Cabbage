package http

import (
	"encoding/json"
	"github.com/borscht/backend/internal/models"
	adminMock "github.com/borscht/backend/internal/restaurantAdmin/mocks"
	"github.com/borscht/backend/utils/errors"
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

func TestDishHandler_UpdateDishData_BindError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	DishUsecaseMock := adminMock.NewMockAdminDishUsecase(ctrl)
	dishHandler := NewDishHandler(DishUsecaseMock)

	requestJSON := `{"name""dish1"}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/restaurant/dish", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := dishHandler.UpdateDishData(c)
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

func TestDishHandler_DeleteDish_BindError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	DishUsecaseMock := adminMock.NewMockAdminDishUsecase(ctrl)
	dishHandler := NewDishHandler(DishUsecaseMock)

	requestJSON := `{"id"1}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/restaurant/dish", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := dishHandler.DeleteDish(c)
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

func TestDishHandler_AddDish_BindError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	DishUsecaseMock := adminMock.NewMockAdminDishUsecase(ctrl)
	dishHandler := NewDishHandler(DishUsecaseMock)

	requestJSON := `{"name""dish1","price":200,"weight":120}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/restaurant/dish", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := dishHandler.AddDish(c)
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

func TestDishHandler_UploadDishImage(t *testing.T) {
	// TODO
}

func TestDishHandler_UploadDishImage_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	DishUsecaseMock := adminMock.NewMockAdminDishUsecase(ctrl)
	dishHandler := NewDishHandler(DishUsecaseMock)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/restaurant/dish", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := dishHandler.UploadDishImage(c)
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
