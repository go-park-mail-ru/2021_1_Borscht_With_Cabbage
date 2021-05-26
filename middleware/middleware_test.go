package middleware

import (
	"fmt"
	"github.com/borscht/backend/internal/restaurantAdmin/mocks"
	authServiceMock "github.com/borscht/backend/internal/services/mocks"
	userMock "github.com/borscht/backend/internal/user/mocks"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"net/http"
	"testing"
)

func TestInitAuthMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	AuthServiceMock := authServiceMock.NewMockServiceAuth(ctrl)
	userUcase := userMock.NewMockUserUsecase(ctrl)
	restaurantUcase := mocks.NewMockAdminRestaurantUsecase(ctrl)

	authMiddleware := InitAuthMiddleware(AuthServiceMock, userUcase, restaurantUcase)
	if authMiddleware == nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestAuthMiddleware_Auth(t *testing.T) {
	//ctrl := gomock.NewController(t)
	//defer ctrl.Finish()
	//userUcase := userMock.NewMockUserUsecase(ctrl)
	//restaurantUcase := restaurantMock.NewMockAdminRestaurantUsecase(ctrl)
	//sessionUsecase := sessionMock.NewMockSessionUsecase(ctrl)
	////userHandler := httpUser.NewUserHandler(userUcase, restaurantUcase, sessionUsecase)
	//
	//authMiddleware := InitAuthMiddleware(userUcase, restaurantUcase, sessionUsecase)
	//e := echo.New()
	//auth := e.Group("", authMiddleware.Auth)
	////auth.GET("/auth", userHandler.CheckAuth)
	//auth.GET("/auth", testFunc)
	//go e.Start(":3000")
	//time.Sleep(1 * time.Second)
	//res, err := http.Get("http://127.0.0.1:5000/auth")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//b := errors.SendError{}
	//bodyBytes, errr := ioutil.ReadAll(res.Body)
	//if errr != nil {
	//	t.Errorf("incorrect result")
	//	return
	//}
	//
	//err = json.Unmarshal(bodyBytes, &b)
	//if err != nil {
	//	t.Errorf("incorrect result")
	//	return
	//}
	//
	//if b.Code == 200 {
	//	t.Errorf("incorrect result")
	//	return
	//}
}

func TestInitUserMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	AuthServiceMock := authServiceMock.NewMockServiceAuth(ctrl)

	userMiddleware := InitUserMiddleware(AuthServiceMock)
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
	AuthServiceMock := authServiceMock.NewMockServiceAuth(ctrl)
	restaurantUcase := mocks.NewMockAdminRestaurantUsecase(ctrl)

	adminMiddleware := InitAdminMiddleware(AuthServiceMock, restaurantUcase)
	if adminMiddleware == nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestAdminAuthMiddleware_Auth(t *testing.T) {

}

func TestCORS(t *testing.T) {
	//e := echo.New()
	//e.Use(CORS)
	//e.Use(LogMiddleware)
	//
	//e.GET("/", testFunc)
	//go e.Start(":3000")
	//time.Sleep(1 * time.Second)
	//res, err := http.Get("http://127.0.0.1:3000")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println(res)
	//
	//if res.StatusCode != 200 {
	//	t.Errorf("incorrect result")
	//	return
	//}
}

func testFunc(c echo.Context) error {
	fmt.Println(23232)
	return c.String(http.StatusOK, "hey")
}
