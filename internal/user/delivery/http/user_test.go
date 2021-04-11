package http

import (
	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/models"
	adminMock "github.com/borscht/backend/internal/restaurantAdmin/mocks"
	sessionMock "github.com/borscht/backend/internal/session/mocks"
	userMock "github.com/borscht/backend/internal/user/mocks"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func mockInit() {
}

func TestUserCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	UserUsecaseMock := userMock.NewMockUserUsecase(ctrl)
	AdminUsecaseMock := adminMock.NewMockAdminUsecase(ctrl)
	SessionUseCaseMock := sessionMock.NewMockSessionUsecase(ctrl)

	input := models.User{
		Email:    "daria@mail.ru",
		Phone:    "89161166000",
		Name:     "Daria",
		Password: "111111",
	}
	inputJSON := `{"email":"daria@mail.ru","number":"89161166000","name":"Daria","password":"111111"}`

	output := models.User{
		Email:    "daria@mail.ru",
		Phone:    "89161166000",
		Name:     "Daria",
		Password: "111111",
		Uid:      1,
	}

	sessionInfo := models.SessionInfo{
		Id:   output.Uid,
		Role: config.RoleUser,
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(inputJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	userHandler := NewUserHandler(UserUsecaseMock, AdminUsecaseMock, SessionUseCaseMock)
	ctx := models.GetContext(c)
	UserUsecaseMock.EXPECT().Create(ctx, input).Return(&output, nil)
	SessionUseCaseMock.EXPECT().Create(ctx, sessionInfo)

	err := userHandler.Create(c)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}
