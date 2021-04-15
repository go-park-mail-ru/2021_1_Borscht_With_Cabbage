package http

import (
	"encoding/json"
	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/models"
	adminMock "github.com/borscht/backend/internal/restaurantAdmin/mocks"
	"github.com/borscht/backend/internal/session/mocks"
	"github.com/borscht/backend/utils/errors"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewRestaurantHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	RestaurantUsecaseMock := adminMock.NewMockAdminRestaurantUsecase(ctrl)
	SessionMock := mocks.NewMockSessionUsecase(ctrl)
	restaurantHandler := NewRestaurantHandler(RestaurantUsecaseMock, SessionMock)
	if restaurantHandler == nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestRestaurantHandler_UpdateRestaurantData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	RestaurantUsecaseMock := adminMock.NewMockAdminRestaurantUsecase(ctrl)
	SessionMock := mocks.NewMockSessionUsecase(ctrl)
	restaurantHandler := NewRestaurantHandler(RestaurantUsecaseMock, SessionMock)

	restaurant := models.RestaurantUpdateData{
		Title: "newName", Description: "yo", DeliveryCost: 200,
	}

	restaurantInfo := models.RestaurantInfo{
		ID:          1,
		Title:       "newName",
		Description: "yo", // итд
	}
	response := models.SuccessRestaurantResponse{
		restaurantInfo, config.RoleAdmin,
	}
	requestJSON := `{"Title":"newName","Description":"yo","deliveryCost":200}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/restaurant", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	ctx := models.GetContext(c)

	RestaurantUsecaseMock.EXPECT().UpdateRestaurantData(ctx, restaurant).Return(&response, nil)

	err := restaurantHandler.UpdateRestaurantData(c)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestRestaurantHandler_UpdateRestaurantData_BindError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	RestaurantUsecaseMock := adminMock.NewMockAdminRestaurantUsecase(ctrl)
	SessionMock := mocks.NewMockSessionUsecase(ctrl)
	restaurantHandler := NewRestaurantHandler(RestaurantUsecaseMock, SessionMock)

	requestJSON := `{"Title""newName","Description":"yo","deliveryCost":200}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/restaurant", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := restaurantHandler.UpdateRestaurantData(c)
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

func TestRestaurantHandler_CreateRestaurant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	RestaurantUsecaseMock := adminMock.NewMockAdminRestaurantUsecase(ctrl)
	SessionMock := mocks.NewMockSessionUsecase(ctrl)
	restaurantHandler := NewRestaurantHandler(RestaurantUsecaseMock, SessionMock)

	newRestaurant := models.RestaurantInfo{
		Title:         "newName",
		AdminPassword: "111111",
		AdminPhone:    "89111111111",
		AdminEmail:    "dasha@mail.ru",
	}
	createdRestaurant := newRestaurant
	createdRestaurant.ID = 1
	response := models.SuccessRestaurantResponse{
		createdRestaurant, config.RoleAdmin,
	}
	requestJSON := `{"Title":"newName","password":"111111","number":"89111111111","email":"dasha@mail.ru"}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/restaurant/signup", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	ctx := models.GetContext(c)

	sessionInfo := models.SessionInfo{
		Id:   createdRestaurant.ID,
		Role: config.RoleAdmin,
	}

	RestaurantUsecaseMock.EXPECT().CreateRestaurant(ctx, newRestaurant).Return(&response, nil)
	SessionMock.EXPECT().Create(ctx, sessionInfo)

	err := restaurantHandler.CreateRestaurant(c)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestRestaurantHandler_CreateRestaurant_BindError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	RestaurantUsecaseMock := adminMock.NewMockAdminRestaurantUsecase(ctrl)
	SessionMock := mocks.NewMockSessionUsecase(ctrl)
	restaurantHandler := NewRestaurantHandler(RestaurantUsecaseMock, SessionMock)

	requestJSON := `{"Title""newName","password":"111111","number":"89111111111","email":"dasha@mail.ru"}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/restaurant/signup", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := restaurantHandler.CreateRestaurant(c)
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

func TestRestaurantHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	RestaurantUsecaseMock := adminMock.NewMockAdminRestaurantUsecase(ctrl)
	SessionMock := mocks.NewMockSessionUsecase(ctrl)
	restaurantHandler := NewRestaurantHandler(RestaurantUsecaseMock, SessionMock)

	newRestaurant := models.RestaurantAuth{
		Login:    "dasha@mail.ru",
		Password: "111111",
	}
	requestJSON := `{"password":"111111","login":"dasha@mail.ru"}`

	existingRestaurant := models.RestaurantInfo{
		ID:    1,
		Title: "rest1",
	}
	response := models.SuccessRestaurantResponse{
		existingRestaurant, config.RoleAdmin,
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/restaurant/signin", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	ctx := models.GetContext(c)

	sessionInfo := models.SessionInfo{
		Id:   existingRestaurant.ID,
		Role: config.RoleAdmin,
	}

	RestaurantUsecaseMock.EXPECT().CheckRestaurantExists(ctx, newRestaurant).Return(&response, nil)
	SessionMock.EXPECT().Create(ctx, sessionInfo)

	err := restaurantHandler.Login(c)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestRestaurantHandler_Login_BindError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	RestaurantUsecaseMock := adminMock.NewMockAdminRestaurantUsecase(ctrl)
	SessionMock := mocks.NewMockSessionUsecase(ctrl)
	restaurantHandler := NewRestaurantHandler(RestaurantUsecaseMock, SessionMock)

	requestJSON := `{"password""111111","login":"dasha@mail.ru"}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/restaurant/signin", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := restaurantHandler.Login(c)
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

func TestRestaurantHandler_GetUserData(t *testing.T) {
	// TODO
}

func TestRestaurantHandler_UploadRestaurantImage(t *testing.T) {
	// TODO
}

func TestRestaurantHandler_UploadRestaurantImage_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	RestaurantUsecaseMock := adminMock.NewMockAdminRestaurantUsecase(ctrl)
	SessionMock := mocks.NewMockSessionUsecase(ctrl)
	restaurantHandler := NewRestaurantHandler(RestaurantUsecaseMock, SessionMock)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/restaurant/signin", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := restaurantHandler.UploadRestaurantImage(c)
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
