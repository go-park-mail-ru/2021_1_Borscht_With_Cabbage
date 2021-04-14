package middleware

import (
	"encoding/json"
	restaurantMock "github.com/borscht/backend/internal/restaurantAdmin/mocks"
	sessionMock "github.com/borscht/backend/internal/session/mocks"
	httpUser "github.com/borscht/backend/internal/user/delivery/http"
	userMock "github.com/borscht/backend/internal/user/mocks"
	"github.com/borscht/backend/utils/errors"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestInitAuthMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userUcase := userMock.NewMockUserUsecase(ctrl)
	restaurantUcase := restaurantMock.NewMockAdminRestaurantUsecase(ctrl)
	sessionUsecase := sessionMock.NewMockSessionUsecase(ctrl)

	authMiddleware := InitAuthMiddleware(userUcase, restaurantUcase, sessionUsecase)
	if authMiddleware == nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestAuthMiddleware_Auth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userUcase := userMock.NewMockUserUsecase(ctrl)
	restaurantUcase := restaurantMock.NewMockAdminRestaurantUsecase(ctrl)
	sessionUsecase := sessionMock.NewMockSessionUsecase(ctrl)
	userHandler := httpUser.NewUserHandler(userUcase, restaurantUcase, sessionUsecase)

	authMiddleware := InitAuthMiddleware(userUcase, restaurantUcase, sessionUsecase)
	e := echo.New()
	auth := e.Group("", authMiddleware.Auth)
	auth.GET("/auth", userHandler.CheckAuth)
	go e.Start(":3000")
	res, err := http.Get("http://127.0.0.1:5000/auth")
	if err != nil {
		t.Fatal(err)
	}

	b := errors.SendError{}
	bodyBytes, errr := ioutil.ReadAll(res.Body)
	if errr != nil {
		t.Errorf("incorrect result")
		return
	}

	err = json.Unmarshal(bodyBytes, &b)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}

	if b.Code == 200 {
		t.Errorf("incorrect result")
		return
	}
}

func TestInitUserMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userUcase := userMock.NewMockUserUsecase(ctrl)
	sessionUsecase := sessionMock.NewMockSessionUsecase(ctrl)

	userMiddleware := InitUserMiddleware(userUcase, sessionUsecase)
	if userMiddleware == nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestUserAuthMiddleware_Auth(t *testing.T) {

}

func TestInitAdminMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	restaurantUcase := restaurantMock.NewMockAdminRestaurantUsecase(ctrl)
	sessionUsecase := sessionMock.NewMockSessionUsecase(ctrl)

	adminMiddleware := InitAdminMiddleware(restaurantUcase, sessionUsecase)
	if adminMiddleware == nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestAdminAuthMiddleware_Auth(t *testing.T) {

}
