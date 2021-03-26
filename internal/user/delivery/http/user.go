package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/models"
	sessionModel "github.com/borscht/backend/internal/session"
	userModel "github.com/borscht/backend/internal/user"
	errors "github.com/borscht/backend/utils"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	UserUcase    userModel.UserUsecase
	SessionUcase sessionModel.SessionUsecase
}

func NewUserHandler(userUcase userModel.UserUsecase, sessionUcase sessionModel.SessionUsecase) userModel.UserHandler {
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
	ctx := models.GetContext(c)

	newUser := new(models.User)
	if err := c.Bind(newUser); err != nil {
		sendErr := errors.Authorization("error with request data")
		return models.SendResponseWithError(c, sendErr)
	}

	if err := h.UserUcase.Create(ctx, *newUser); err != nil {
		return models.SendResponseWithError(c, err)
	}

	session, err := h.SessionUcase.Create(ctx, newUser.Phone)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	fmt.Println("SESSION:", session)
	setResponseCookie(c, session)

	response := models.SuccessResponse{Name: newUser.Name, Avatar: config.DefaultAvatar} // TODO убрать config отсюда
	return models.SendResponse(c, response)
}

func (h *Handler) Login(c echo.Context) error {
	ctx := models.GetContext(c)

	newUser := new(models.UserAuth)
	if err := c.Bind(newUser); err != nil {
		sendErr := errors.Authorization("error with request data")
		return models.SendResponseWithError(c, sendErr)
	}

	oldUser, err := h.UserUcase.GetByLogin(ctx, *newUser)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	session, err := h.SessionUcase.Create(ctx, oldUser.Phone)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	setResponseCookie(c, session)

	response := models.SuccessResponse{Name: oldUser.Name, Avatar: oldUser.Avatar}
	return models.SendResponse(c, response)
}

func (h *Handler) GetUserData(c echo.Context) error {
	user := c.Get("User")

	if user == nil {
		userError := errors.Authorization("not authorization")
		return models.SendResponseWithError(c, userError)
	}

	return models.SendResponse(c, user)
}

func (h *Handler) EditProfile(c echo.Context) error {
	ctx := models.GetContext(c)

	// TODO убрать часть логики в usecase
	formParams, err := c.FormParams()
	if err != nil {
		return errors.BadRequest("invalid data form")
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

	fmt.Println("FILE DELIVERY: ", err)

	if file != nil {
		if err != nil {
			return models.SendResponseWithError(c, errors.BadRequest(err.Error()))
		}

		filename, err := h.UserUcase.UploadAvatar(ctx, file)
		if err != nil {
			return models.SendResponseWithError(c, err)
		}

		fmt.Println("FILENAME: ", filename)

		profileEdits.Avatar = filename
	}

	user := c.Get("User")

	err = h.UserUcase.Update(ctx, profileEdits)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}
	err = h.SessionUcase.Update(ctx, profileEdits.Phone, user.(models.User).Phone)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	return models.SendResponse(c, profileEdits)
}

func (h *Handler) CheckAuth(c echo.Context) error {
	user := c.Get("User")

	if user == nil {
		sendErr := errors.Authorization("error with request data")
		return models.SendResponseWithError(c, sendErr)
	}

	return models.SendResponse(c, user)
}

func (h *Handler) Logout(c echo.Context) error {
	ctx := models.GetContext(c)

	cook, err := c.Cookie(config.SessionCookie)
	if err != nil {
		sendErr := errors.Authorization("error with request data")
		return models.SendResponseWithError(c, sendErr)
	}

	h.SessionUcase.Delete(ctx, cook.Value)

	deleteResponseCookie(c)

	return models.SendResponse(c, "")
}
