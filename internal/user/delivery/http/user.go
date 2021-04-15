package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/borscht/backend/utils/validation"

	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/models"
	adminModel "github.com/borscht/backend/internal/restaurantAdmin"
	sessionModel "github.com/borscht/backend/internal/session"
	userModel "github.com/borscht/backend/internal/user"
	errors "github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	UserUcase    userModel.UserUsecase
	AdminUcase   adminModel.AdminRestaurantUsecase
	SessionUcase sessionModel.SessionUsecase
}

func NewUserHandler(userUcase userModel.UserUsecase, adminUcase adminModel.AdminRestaurantUsecase, sessionUcase sessionModel.SessionUsecase) userModel.UserHandler {
	handler := &Handler{
		UserUcase:    userUcase,
		AdminUcase:   adminUcase,
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

func (h Handler) Create(c echo.Context) error {
	ctx := models.GetContext(c)

	newUser := new(models.User)

	if err := c.Bind(newUser); err != nil {
		sendErr := errors.AuthorizationError("error with request data")
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}

	if err := validation.ValidateUserRegistration(*newUser); err != nil {
		return models.SendResponseWithError(c, err)
	}

	responseUser, err := h.UserUcase.Create(ctx, *newUser)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	sessionInfo := models.SessionInfo{
		Id:   responseUser.Uid,
		Role: config.RoleUser,
	}

	session, err := h.SessionUcase.Create(ctx, sessionInfo)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	setResponseCookie(c, session)

	return models.SendResponse(c, responseUser)
}

// TODO: убрать эту логику отсюда
func (h Handler) Login(c echo.Context) error {
	ctx := models.GetContext(c)

	newUser := new(models.UserAuth)

	if err := c.Bind(newUser); err != nil {
		sendErr := errors.AuthorizationError("error with request data")
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}

	if err := validation.ValidateSignIn(newUser.Login, newUser.Password); err != nil {
		fmt.Println(err)
		return models.SendResponseWithError(c, err)
	}

	oldUser, err := h.UserUcase.CheckUserExists(ctx, *newUser)

	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	sessionInfo := models.SessionInfo{
		Id:   oldUser.Uid,
		Role: config.RoleUser,
	}
	session, err := h.SessionUcase.Create(ctx, sessionInfo)

	if err != nil {
		return models.SendResponseWithError(c, err)
	}
	setResponseCookie(c, session)

	return models.SendResponse(c, oldUser)
}

func (h Handler) GetUserData(c echo.Context) error {
	ctx := models.GetContext(c)

	user, err := h.UserUcase.GetUserData(ctx)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	return models.SendResponse(c, user)
}

func (h Handler) UpdateData(c echo.Context) error {
	ctx := models.GetContext(c)

	user := new(models.UserData)
	if err := c.Bind(user); err != nil {
		sendErr := errors.BadRequestError(err.Error())
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}

	responseUser, err := h.UserUcase.UpdateData(ctx, *user)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	return models.SendResponse(c, responseUser)
}

func (h Handler) UploadAvatar(c echo.Context) error {
	ctx := models.GetContext(c)

	file, err := c.FormFile("avatar")
	if err != nil {
		requestError := errors.BadRequestError(err.Error())
		logger.DeliveryLevel().ErrorLog(ctx, requestError)
		return models.SendResponseWithError(c, requestError)
	}

	response, err := h.UserUcase.UploadAvatar(ctx, file)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	return models.SendResponse(c, response)
}

// TODO: подумать как это можно изменить
func (h Handler) CheckAuth(c echo.Context) error {
	ctx := models.GetContext(c)
	cookie, err := c.Cookie(config.SessionCookie)
	if err != nil {
		sendErr := errors.NewCustomError(http.StatusUnauthorized, "error with request data")
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}

	sessionData := new(models.SessionInfo)
	var exist bool
	*sessionData, exist, err = h.SessionUcase.Check(ctx, cookie.Value)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}
	if !exist {
		sendErr := errors.NewCustomError(http.StatusUnauthorized, "error with request data")
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}

	switch sessionData.Role {
	case config.RoleAdmin:
		restaurant, err := h.AdminUcase.GetByRid(ctx, sessionData.Id)
		if err != nil {
			sendErr := errors.NewCustomError(http.StatusUnauthorized, err.Error())
			logger.DeliveryLevel().ErrorLog(ctx, sendErr)
			return models.SendResponseWithError(c, sendErr)
		}
		return models.SendResponse(c, restaurant)

	case config.RoleUser:
		user, err := h.UserUcase.GetByUid(ctx, sessionData.Id)
		if err != nil {
			sendErr := errors.NewCustomError(http.StatusUnauthorized, err.Error())
			logger.DeliveryLevel().ErrorLog(ctx, sendErr)
			return models.SendResponseWithError(c, sendErr)
		}
		return models.SendResponse(c, user)
	default:
		sendErr := errors.NewCustomError(http.StatusUnauthorized, "error with roles")
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}
}

func (h Handler) Logout(c echo.Context) error {
	ctx := models.GetContext(c)

	cook, err := c.Cookie(config.SessionCookie)
	if err != nil {
		sendErr := errors.AuthorizationError("error with request data")
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}

	err = h.SessionUcase.Delete(ctx, cook.Value)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}
	deleteResponseCookie(c)

	return models.SendResponse(c, nil)
}
