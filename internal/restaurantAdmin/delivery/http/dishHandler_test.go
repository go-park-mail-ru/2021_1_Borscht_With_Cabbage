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
	"testing"
)

func TestNewDishHandlerew(t *testing.T) {

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
	req := httptest.NewRequest(http.MethodPost, "/restaurant/dishes", nil)
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

	b := errors.SendError{}
	respCode := rec.Body.Bytes()
	err = json.Unmarshal(respCode, &b)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
	if b.Code != 200 {
		t.Errorf("incorrect result")
		return
	}
}

func TestDishHandler_UpdateDishData(t *testing.T) {

}

func TestDishHandler_DeleteDish(t *testing.T) {

}

func TestDishHandler_AddDish(t *testing.T) {

}

func TestDishHandler_UploadDishImage(t *testing.T) {
	// TODO
}
