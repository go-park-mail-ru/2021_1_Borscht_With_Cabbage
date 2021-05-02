package http

import (
	"encoding/json"
	authServiceMock "github.com/borscht/backend/internal/services/mocks"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/models"
	userMock "github.com/borscht/backend/internal/user/mocks"
	"github.com/borscht/backend/utils/errors"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
)

func TestHandler_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	UserUsecaseMock := userMock.NewMockUserUsecase(ctrl)
	AuthServiceMock := authServiceMock.NewMockServiceAuth(ctrl)
	userHandler := NewUserHandler(UserUsecaseMock, AuthServiceMock)

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
	response := models.SuccessUserResponse{
		User: output,
		Role: config.RoleUser,
	}
	address := models.Address{}
	response.Address = address

	sessionInfo := models.SessionInfo{
		Id:   output.Uid,
		Role: config.RoleUser,
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(inputJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	ctx := models.GetContext(c)

	AuthServiceMock.EXPECT().Create(ctx, input).Return(&response, nil)
	UserUsecaseMock.EXPECT().AddAddress(ctx, output.Uid, address).Return(nil)
	AuthServiceMock.EXPECT().CreateSession(ctx, sessionInfo).Return("session", nil)

	err := userHandler.Create(c)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestHandler_CreateUser_BindError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	UserUsecaseMock := userMock.NewMockUserUsecase(ctrl)
	AuthServiceMock := authServiceMock.NewMockServiceAuth(ctrl)
	userHandler := NewUserHandler(UserUsecaseMock, AuthServiceMock)

	inputJSON := `{"emailaria@mail.ru","number":"89161166000","name":"Daria","password":"111111"}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(inputJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := userHandler.Create(c)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestSignup_BindError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	UserUsecaseMock := userMock.NewMockUserUsecase(ctrl)
	AuthServiceMock := authServiceMock.NewMockServiceAuth(ctrl)
	userHandler := NewUserHandler(UserUsecaseMock, AuthServiceMock)

	inputJSON := `{email:daria@mail.ru}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(inputJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := userHandler.Create(c)
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
	if b.Code == 200 {
		t.Errorf("incorrect result")
		return
	}
}

func TestHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	UserUsecaseMock := userMock.NewMockUserUsecase(ctrl)
	AuthServiceMock := authServiceMock.NewMockServiceAuth(ctrl)
	userHandler := NewUserHandler(UserUsecaseMock, AuthServiceMock)

	input := models.UserAuth{
		Login:    "daria@mail.ru",
		Password: "111111",
	}
	inputJSON := `{"login":"daria@mail.ru","password":"111111"}`

	output := models.User{
		Email:    "daria@mail.ru",
		Phone:    "89161166000",
		Name:     "Daria",
		Password: "111111",
		Uid:      1,
	}
	response := models.SuccessUserResponse{
		output, config.RoleUser,
	}

	sessionInfo := models.SessionInfo{
		Id:   output.Uid,
		Role: config.RoleUser,
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/signin", strings.NewReader(inputJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	ctx := models.GetContext(c)

	AuthServiceMock.EXPECT().CheckUserExists(ctx, input).Return(&response, nil)
	AuthServiceMock.EXPECT().CreateSession(ctx, sessionInfo)

	err := userHandler.Login(c)

	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestHandler_Login_BindError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	UserUsecaseMock := userMock.NewMockUserUsecase(ctrl)
	AuthServiceMock := authServiceMock.NewMockServiceAuth(ctrl)
	userHandler := NewUserHandler(UserUsecaseMock, AuthServiceMock)

	inputJSON := `{"logindaria@mail.ru","password":"111111"}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/signin", strings.NewReader(inputJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := userHandler.Login(c)

	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestHandler_LoginValidationFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	UserUsecaseMock := userMock.NewMockUserUsecase(ctrl)
	AuthServiceMock := authServiceMock.NewMockServiceAuth(ctrl)
	userHandler := NewUserHandler(UserUsecaseMock, AuthServiceMock)

	inputJSON := `{"login":"dariamail.ru","password":"111111"}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/signin", strings.NewReader(inputJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := userHandler.Login(c)
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
	if b.Code == 200 {
		t.Errorf("incorrect result")
		return
	}
}

func TestHandler_GetUserData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	UserUsecaseMock := userMock.NewMockUserUsecase(ctrl)
	AuthServiceMock := authServiceMock.NewMockServiceAuth(ctrl)
	userHandler := NewUserHandler(UserUsecaseMock, AuthServiceMock)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/user", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	ctx := models.GetContext(c)

	user := models.User{
		Uid:  1,
		Name: "Daria",
	}
	response := models.SuccessUserResponse{
		user, config.RoleUser,
	}

	UserUsecaseMock.EXPECT().GetUserData(ctx).Return(&response, nil)

	err := userHandler.GetUserData(c)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestHandler_UpdateData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	UserUsecaseMock := userMock.NewMockUserUsecase(ctrl)
	AuthServiceMock := authServiceMock.NewMockServiceAuth(ctrl)
	userHandler := NewUserHandler(UserUsecaseMock, AuthServiceMock)

	input := models.UserData{
		Email:       "daria@mail.ru",
		Phone:       "89161166000",
		Name:        "Daria",
		Password:    "111111",
		PasswordOld: "111111",
	}
	inputJSON := `{"email":"daria@mail.ru","number":"89161166000","name":"Daria","password":"111111","password_current":"111111"}`

	output := models.User{
		Email:    "daria@mail.ru",
		Phone:    "89161166000",
		Name:     "Daria",
		Password: "111111",
		Uid:      1,
	}
	response := models.SuccessUserResponse{
		User: output,
		Role: config.RoleUser,
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/user", strings.NewReader(inputJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	user := models.User{
		Uid:  1,
		Name: "Daria",
	}
	c.Set("User", user)

	ctx := models.GetContext(c)
	UserUsecaseMock.EXPECT().UpdateData(ctx, input).Return(&response, nil)

	err := userHandler.UpdateData(c)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}

}

func TestHandler_UpdateData_BindError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	UserUsecaseMock := userMock.NewMockUserUsecase(ctrl)
	AuthServiceMock := authServiceMock.NewMockServiceAuth(ctrl)
	userHandler := NewUserHandler(UserUsecaseMock, AuthServiceMock)

	inputJSON := `{"emaildaria@mail.ru","number":"89161166000","name":"Daria","password":"111111","password_current":"111111"}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/user", strings.NewReader(inputJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	err := userHandler.UpdateData(c)
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

func TestHandler_UpdateDataWithAvatar(t *testing.T) {
	// TODO
}

func TestHandler_UploadAvatar(t *testing.T) {
	// todo
}

func TestHandler_UploadAvatar_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	UserUsecaseMock := userMock.NewMockUserUsecase(ctrl)
	AuthServiceMock := authServiceMock.NewMockServiceAuth(ctrl)
	userHandler := NewUserHandler(UserUsecaseMock, AuthServiceMock)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/auth", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := userHandler.UploadAvatar(c)
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

func TestHandler_CheckAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	UserUsecaseMock := userMock.NewMockUserUsecase(ctrl)
	AuthServiceMock := authServiceMock.NewMockServiceAuth(ctrl)
	userHandler := NewUserHandler(UserUsecaseMock, AuthServiceMock)

	output := models.User{
		Email:    "daria@mail.ru",
		Phone:    "89161166000",
		Name:     "Daria",
		Password: "111111",
		Avatar:   "img.jpg",
		Uid:      1,
	}
	responseUser := models.SuccessUserResponse{
		output, config.RoleUser,
	}

	sessionInfo := models.SessionInfo{
		Id:   output.Uid,
		Role: config.RoleUser,
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/auth", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	sessionCookie := http.Cookie{
		Expires: time.Now().Add(24 * time.Hour),
		Name:    config.SessionCookie,
		Value:   "session1",
	}
	req.AddCookie(&sessionCookie)
	c := e.NewContext(req, rec)

	ctx := models.GetContext(c)
	AuthServiceMock.EXPECT().CheckSession(ctx, "session1").Return(sessionInfo, true, nil)
	AuthServiceMock.EXPECT().GetByUid(ctx, sessionInfo.Id).Return(&responseUser, nil)

	err := userHandler.CheckAuth(c)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestHandler_CheckAuth_GetUserError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	UserUsecaseMock := userMock.NewMockUserUsecase(ctrl)
	AuthServiceMock := authServiceMock.NewMockServiceAuth(ctrl)
	userHandler := NewUserHandler(UserUsecaseMock, AuthServiceMock)

	output := models.User{
		Email:    "daria@mail.ru",
		Phone:    "89161166000",
		Name:     "Daria",
		Password: "111111",
		Avatar:   "img.jpg",
		Uid:      1,
	}
	responseUser := models.SuccessUserResponse{
		output, config.RoleUser,
	}

	sessionInfo := models.SessionInfo{
		Id:   output.Uid,
		Role: config.RoleUser,
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/auth", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	sessionCookie := http.Cookie{
		Expires: time.Now().Add(24 * time.Hour),
		Name:    config.SessionCookie,
		Value:   "session1",
	}
	req.AddCookie(&sessionCookie)
	c := e.NewContext(req, rec)

	ctx := models.GetContext(c)
	AuthServiceMock.EXPECT().CheckSession(ctx, "session1").Return(sessionInfo, true, nil)
	AuthServiceMock.EXPECT().GetByUid(ctx, sessionInfo.Id).Return(&responseUser, errors.BadRequestError("err"))

	err := userHandler.CheckAuth(c)
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

func TestHandler_CheckAuth_GetRestaurantError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	UserUsecaseMock := userMock.NewMockUserUsecase(ctrl)
	AuthServiceMock := authServiceMock.NewMockServiceAuth(ctrl)
	userHandler := NewUserHandler(UserUsecaseMock, AuthServiceMock)

	responseRest := models.SuccessRestaurantResponse{
		Role: config.RoleAdmin,
	}

	sessionInfo := models.SessionInfo{
		Id:   1,
		Role: config.RoleAdmin,
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/auth", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	sessionCookie := http.Cookie{
		Expires: time.Now().Add(24 * time.Hour),
		Name:    config.SessionCookie,
		Value:   "session1",
	}
	req.AddCookie(&sessionCookie)
	c := e.NewContext(req, rec)

	ctx := models.GetContext(c)
	AuthServiceMock.EXPECT().CheckSession(ctx, "session1").Return(sessionInfo, true, nil)
	AuthServiceMock.EXPECT().GetByRid(ctx, sessionInfo.Id).Return(&responseRest, errors.BadRequestError("err"))

	err := userHandler.CheckAuth(c)
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

func TestHandler_CheckAuth_WrongRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	UserUsecaseMock := userMock.NewMockUserUsecase(ctrl)
	AuthServiceMock := authServiceMock.NewMockServiceAuth(ctrl)
	userHandler := NewUserHandler(UserUsecaseMock, AuthServiceMock)

	sessionInfo := models.SessionInfo{
		Id:   1,
		Role: "strange role",
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/auth", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	sessionCookie := http.Cookie{
		Expires: time.Now().Add(24 * time.Hour),
		Name:    config.SessionCookie,
		Value:   "session1",
	}
	req.AddCookie(&sessionCookie)
	c := e.NewContext(req, rec)

	ctx := models.GetContext(c)
	AuthServiceMock.EXPECT().CheckSession(ctx, "session1").Return(sessionInfo, true, nil)

	err := userHandler.CheckAuth(c)
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

func TestHandler_CheckAuth_CookieNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	UserUsecaseMock := userMock.NewMockUserUsecase(ctrl)
	AuthServiceMock := authServiceMock.NewMockServiceAuth(ctrl)
	userHandler := NewUserHandler(UserUsecaseMock, AuthServiceMock)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/auth", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := userHandler.CheckAuth(c)
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

func TestHandler_CheckAuth_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	UserUsecaseMock := userMock.NewMockUserUsecase(ctrl)
	AuthServiceMock := authServiceMock.NewMockServiceAuth(ctrl)
	userHandler := NewUserHandler(UserUsecaseMock, AuthServiceMock)

	output := models.User{
		Email:    "daria@mail.ru",
		Phone:    "89161166000",
		Name:     "Daria",
		Password: "111111",
		Avatar:   "img.jpg",
		Uid:      1,
	}

	sessionInfo := models.SessionInfo{
		Id:   output.Uid,
		Role: config.RoleUser,
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/auth", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	sessionCookie := http.Cookie{
		Expires: time.Now().Add(24 * time.Hour),
		Name:    config.SessionCookie,
		Value:   "session1",
	}
	req.AddCookie(&sessionCookie)
	c := e.NewContext(req, rec)

	ctx := models.GetContext(c)
	AuthServiceMock.EXPECT().CheckSession(ctx, "session1").Return(sessionInfo, false, nil)

	err := userHandler.CheckAuth(c)
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

func TestHandler_Logout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	UserUsecaseMock := userMock.NewMockUserUsecase(ctrl)
	AuthServiceMock := authServiceMock.NewMockServiceAuth(ctrl)
	userHandler := NewUserHandler(UserUsecaseMock, AuthServiceMock)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/auth", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	sessionCookie := http.Cookie{
		Expires: time.Now().Add(24 * time.Hour),
		Name:    config.SessionCookie,
		Value:   "session1",
	}
	req.AddCookie(&sessionCookie)
	c := e.NewContext(req, rec)

	ctx := models.GetContext(c)
	AuthServiceMock.EXPECT().DeleteSession(ctx, "session1").Return(nil)

	err := userHandler.Logout(c)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestHandler_Logout_CookieNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	UserUsecaseMock := userMock.NewMockUserUsecase(ctrl)
	AuthServiceMock := authServiceMock.NewMockServiceAuth(ctrl)
	userHandler := NewUserHandler(UserUsecaseMock, AuthServiceMock)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/auth", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := userHandler.Logout(c)
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
	if b.Code == 200 {
		t.Errorf("incorrect result")
		return
	}
}
