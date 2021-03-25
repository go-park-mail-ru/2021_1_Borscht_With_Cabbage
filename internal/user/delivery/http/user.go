package http

import (
	"fmt"
	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/models"
	_sessionModel "github.com/borscht/backend/internal/session"
	_userModel "github.com/borscht/backend/internal/user"
	_errors "github.com/borscht/backend/utils"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type Handler struct {
	UserUcase    _userModel.UserUsecase
	SessionUcase _sessionModel.SessionUsecase
}

func NewUserHandler(userUcase _userModel.UserUsecase, sessionUcase _sessionModel.SessionUsecase) _userModel.UserHandler {
	handler := &Handler{
		UserUcase:    userUcase,
		SessionUcase: sessionUcase,
	}

	return handler
}

func setResponseCookie(c echo.Context, session string) {
	sessionCookie := http.Cookie{
		Expires:  time.Now().Add(24 * time.Hour),
		Name:     config.SessionCookie,
		Value:    session,
		HttpOnly: true,
	}
	c.SetCookie(&sessionCookie)
}

func deleteResponseCookie(c echo.Context) {
	sessionCookie := http.Cookie{
		Expires:  time.Now().Add(-24 * time.Hour),
		Name:     config.SessionCookie,
		Value:    "session",
		HttpOnly: true,
	}
	c.SetCookie(&sessionCookie)
}

func (h *Handler) Create(c echo.Context) error {
	cc := c.(*models.CustomContext)

	newUser := new(models.User)
	if err := c.Bind(newUser); err != nil {
		sendErr := _errors.NewCustomError(http.StatusUnauthorized, "error with request data")
		return cc.SendResponseWithError(sendErr)
	}

	uid, err := h.UserUcase.Create(*newUser)
	if err != nil {
		return cc.SendResponseWithError(err)
	}

	session, err := h.SessionUcase.Create(uid)
	if err != nil {
		return cc.SendResponseWithError(err)
	}

	fmt.Println("SESSION:", session)
	setResponseCookie(c, session)

	response := models.SuccessResponse{Name: newUser.Name, Avatar: config.DefaultAvatar} // TODO убрать config отсюда
	return cc.SendResponse(response)
}

func (h *Handler) Login(c echo.Context) error {
	cc := c.(*models.CustomContext)
	newUser := new(models.UserAuth)

	if err := c.Bind(newUser); err != nil {
		sendErr := _errors.NewCustomError(http.StatusUnauthorized, "error with request data")
		return cc.SendResponseWithError(sendErr)
	}

	oldUser, err := h.UserUcase.CheckUserExists(*newUser)
	if err != nil {
		return cc.SendResponseWithError(err)
	}

	session, err := h.SessionUcase.Create(oldUser.Uid)
	if err != nil {
		return cc.SendResponseWithError(err)
	}
	fmt.Println("hey")
	setResponseCookie(c, session)
	fmt.Println("hey")

	response := models.SuccessResponse{Name: oldUser.Name, Avatar: oldUser.Avatar}
	fmt.Println("hey")

	return cc.SendResponse(response)
}

func (h *Handler) GetUserData(c echo.Context) error {
	cc := c.(*models.CustomContext)

	if cc.User == nil {
		userError := _errors.Authorization("not authorization")
		return cc.SendResponseWithError(userError)
	}

	return cc.SendResponse(*cc.User)
}

func (h *Handler) EditProfile(c echo.Context) error {
	cc := c.(*models.CustomContext)

	// TODO убрать часть логики в usecase
	formParams, err := c.FormParams()
	if err != nil {
		return _errors.NewCustomError(http.StatusBadRequest, "invalid data form")
	}

	profileEdits := models.UserData{
		Name:        formParams.Get("name"),
		Phone:       formParams.Get("number"),
		Email:       formParams.Get("email"),
		Password:    formParams.Get("password"),
		PasswordOld: formParams.Get("password_current"),
	}
	fmt.Println(profileEdits)

	file, err := c.FormFile("avatar")
	if err != nil {
		return cc.SendResponseWithError(_errors.BadRequest(err.Error()))
	}
	filename, err := h.UserUcase.UploadAvatar(file)
	if err != nil {
		return cc.SendResponseWithError(err)
	}

	profileEdits.Avatar = filename

	if cc.User == nil {
		userError := _errors.Authorization("not authorization")
		return cc.SendResponseWithError(userError)
	}

	err = h.UserUcase.Update(profileEdits, cc.User.Uid)
	if err != nil {
		return cc.SendResponseWithError(err)
	}

	return cc.SendResponse(profileEdits)
}

func (h *Handler) CheckAuth(c echo.Context) error {
	cc := c.(*models.CustomContext)

	if cc.User == nil {
		sendErr := _errors.NewCustomError(http.StatusUnauthorized, "error with request data")
		return cc.SendResponseWithError(sendErr)
	}

	return cc.SendResponse(cc.User)
}

func (h *Handler) Logout(c echo.Context) error {
	cc := c.(*models.CustomContext)

	cook, err := cc.Cookie(config.SessionCookie)
	if err != nil {
		sendErr := _errors.NewCustomError(http.StatusUnauthorized, "error with request data")
		return cc.SendResponseWithError(sendErr)
	}

	err = h.SessionUcase.Delete(cook.Value)
	if err != nil {
		return cc.SendResponseWithError(err)
	}

	deleteResponseCookie(c)

	return cc.SendResponse("")
}
