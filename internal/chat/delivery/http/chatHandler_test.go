package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/borscht/backend/configProject"
	"github.com/borscht/backend/internal/chat/mocks"
	"github.com/borscht/backend/internal/models"
	authServiceMock "github.com/borscht/backend/internal/services/mocks"
	"github.com/borscht/backend/utils/errors"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
)

func TestChatHandler_GetKey_User(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	AuthServiceMock := authServiceMock.NewMockServiceAuth(ctrl)
	ChatUsecaseMock := mocks.NewMockChatUsecase(ctrl)
	ChatHandler := NewChatHandler(ChatUsecaseMock, AuthServiceMock)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/connect/ws", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	user := models.User{
		Uid:  1,
		Name: "Daria",
	}
	c.Set("User", user)
	ctx := models.GetContext(c)

	AuthServiceMock.EXPECT().CreateKey(ctx, models.SessionInfo{
		Id: 1, Role: configProject.RoleUser,
	}).Return("key", nil)

	err := ChatHandler.GetKey(c)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestChatHandler_GetKey_Restaurant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	AuthServiceMock := authServiceMock.NewMockServiceAuth(ctrl)
	ChatUsecaseMock := mocks.NewMockChatUsecase(ctrl)
	ChatHandler := NewChatHandler(ChatUsecaseMock, AuthServiceMock)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/connect/ws", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	user := models.RestaurantInfo{
		ID: 1,
	}
	c.Set("Restaurant", user)
	ctx := models.GetContext(c)

	AuthServiceMock.EXPECT().CreateKey(ctx, models.SessionInfo{
		Id: 1, Role: configProject.RoleAdmin,
	}).Return("key", nil)

	err := ChatHandler.GetKey(c)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestChatHandler_Connect(t *testing.T) {
	//ctrl := gomock.NewController(t)
	//defer ctrl.Finish()
	//AuthServiceMock := authServiceMock.NewMockServiceAuth(ctrl)
	//ChatUsecaseMock := mocks.NewMockChatUsecase(ctrl)
	//ChatHandler := NewChatHandler(ChatUsecaseMock, AuthServiceMock)
	//
	//e := echo.New()
	//req := httptest.NewRequest(http.MethodGet, "/ws/1", nil)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	//rec := httptest.NewRecorder()
	//c := e.NewContext(req, rec)
	//
	//ctx := models.GetContext(c)
	//
	//
	//ChatUsecaseMock.EXPECT().Connect(ctx, ).Return(nil)
	//
	//err := ChatHandler.GetKey(c)
	//if err != nil {
	//	t.Errorf("incorrect result")
	//	return
	//}
}

func TestChatHandler_GetAllChats(t *testing.T) {
	//ctrl := gomock.NewController(t)
	//defer ctrl.Finish()
	//AuthServiceMock := authServiceMock.NewMockServiceAuth(ctrl)
	//ChatUsecaseMock := mocks.NewMockChatUsecase(ctrl)
	//ChatHandler := NewChatHandler(ChatUsecaseMock, AuthServiceMock)
	//
	//e := echo.New()
	//req := httptest.NewRequest(http.MethodGet, "/chats", nil)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	//rec := httptest.NewRecorder()
	//c := e.NewContext(req, rec)
	//
	//ctx := models.GetContext(c)
	//ChatUsecaseMock.EXPECT().GetAllChats(ctx).Return([]models.BriefInfoChat{}, nil)
	//
	//err := ChatHandler.GetKey(c)
	//if err != nil {
	//	t.Errorf("incorrect result")
	//	return
	//}
}

func TestChatHandler_GetAllMessages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	AuthServiceMock := authServiceMock.NewMockServiceAuth(ctrl)
	ChatUsecaseMock := mocks.NewMockChatUsecase(ctrl)
	chatHandler := NewChatHandler(ChatUsecaseMock, AuthServiceMock)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/chats", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	ctx := models.GetContext(c)

	ChatUsecaseMock.EXPECT().GetAllMessages(ctx, 1).Return(&models.InfoChat{}, nil)

	err := chatHandler.GetAllMessages(c)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestChatHandler_GetAllMessages_NoIdError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	AuthServiceMock := authServiceMock.NewMockServiceAuth(ctrl)
	ChatUsecaseMock := mocks.NewMockChatUsecase(ctrl)
	chatHandler := NewChatHandler(ChatUsecaseMock, AuthServiceMock)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/chats", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := chatHandler.GetAllMessages(c)
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

func TestChatHandler_GetAllMessages_BadIdError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	AuthServiceMock := authServiceMock.NewMockServiceAuth(ctrl)
	ChatUsecaseMock := mocks.NewMockChatUsecase(ctrl)
	chatHandler := NewChatHandler(ChatUsecaseMock, AuthServiceMock)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/chats", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("id")

	err := chatHandler.GetAllMessages(c)
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
