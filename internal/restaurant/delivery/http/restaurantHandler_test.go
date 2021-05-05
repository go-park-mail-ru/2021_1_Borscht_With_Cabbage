package http

import (
	"github.com/borscht/backend/internal/models"
	restModel "github.com/borscht/backend/internal/restaurant"
	"github.com/borscht/backend/internal/restaurant/mocks"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRestaurantHandler_GetVendor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	RestaurantUsecaseMock := mocks.NewMockRestaurantUsecase(ctrl)
	restaurantHandler := NewRestaurantHandler(RestaurantUsecaseMock)

	response := []models.RestaurantInfo{
		{ID: 1, AvgCheck: 1000, Title: "Rest1", Description: "new", DeliveryCost: 200, Rating: 5, Avatar: "img.jpg"},
		{ID: 2, AvgCheck: 1200, Title: "Rest2", Description: "new2", DeliveryCost: 160, Rating: 4, Avatar: "img.jpg"},
	}

	params := restModel.GetVendorParams{
		Limit:     1,
		Offset:    2,
		Address:   true,
		Latitude:  "55.768096",
		Longitude: "37.646839",
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/restaurants/?limit=1&offset=2&longitude=37.646839&latitude=55.768096", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	ctx := models.GetContext(c)

	RestaurantUsecaseMock.EXPECT().GetVendor(ctx, params).Return(response, nil)

	err := restaurantHandler.GetVendor(c)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestRestaurantHandler_GetRestaurantPage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	RestaurantUsecaseMock := mocks.NewMockRestaurantUsecase(ctrl)
	restaurantHandler := NewRestaurantHandler(RestaurantUsecaseMock)

	dishes := []models.Dish{
		{Name: "dish", Price: 200},
		{Name: "dish2", Price: 300},
	}

	restaurant := models.RestaurantInfo{
		ID: 1, AvgCheck: 1000, Title: "Rest1", Description: "new", DeliveryCost: 200, Rating: 5, Avatar: "img.jpg",
	}

	response := models.RestaurantWithDishes{
		RestaurantInfo: restaurant, Dishes: dishes,
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	ctx := models.GetContext(c)

	RestaurantUsecaseMock.EXPECT().GetById(ctx, models.Coordinates{}, 1).Return(&response, nil)

	err := restaurantHandler.GetRestaurantPage(c)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}
