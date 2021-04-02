package http

import (
	"net/http"
	"time"

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
	AdminUcase   adminModel.AdminUsecase
	SessionUcase sessionModel.SessionUsecase
}

func NewUserHandler(userUcase userModel.UserUsecase, adminUcase adminModel.AdminUsecase, sessionUcase sessionModel.SessionUsecase) userModel.UserHandler {
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

func (h *Handler) Create(c echo.Context) error {
	ctx := models.GetContext(c)

	newUser := new(models.User)

	if err := c.Bind(newUser); err != nil {
		sendErr := errors.AuthorizationError("error with request data")
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}

	uid, err := h.UserUcase.Create(ctx, *newUser)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	sessionInfo := models.SessionInfo{
		Id:   uid,
		Role: config.RoleUser,
	}
	session, err := h.SessionUcase.Create(ctx, sessionInfo)

	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	setResponseCookie(c, session)

	response := models.SuccessResponse{Name: newUser.Name, Avatar: config.DefaultAvatar, Role: config.RoleUser} // TODO убрать config отсюда
	return models.SendResponse(c, response)
}

func (h *Handler) Login(c echo.Context) error {
	ctx := models.GetContext(c)

	newUser := new(models.UserAuth)

	if err := c.Bind(newUser); err != nil {
		sendErr := errors.AuthorizationError("error with request data")
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
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

	response := models.SuccessResponse{Name: oldUser.Name, Avatar: oldUser.Avatar, Role: config.RoleUser}

	return models.SendResponse(c, response)
}

func (h *Handler) GetUserData(c echo.Context) error {
	ctx := models.GetContext(c)

	user := c.Get("User")

	if user == nil {
		userError := errors.AuthorizationError("not authorization")
		logger.DeliveryLevel().ErrorLog(ctx, userError)
		return models.SendResponseWithError(c, userError)
	}

	return models.SendResponse(c, user)
}

func (h *Handler) EditProfile(c echo.Context) error {
	ctx := models.GetContext(c)

	formParams, err := c.FormParams()
	if err != nil {
		requestError := errors.BadRequestError("invalid data form")
		logger.DeliveryLevel().ErrorLog(ctx, requestError)
		return models.SendResponseWithError(c, requestError)
	}

	profileEdits := models.UserData{
		Name:        formParams.Get("name"),
		Phone:       formParams.Get("number"),
		Email:       formParams.Get("email"),
		Password:    formParams.Get("password"),
		PasswordOld: formParams.Get("password_current"),
	}

	file, err := c.FormFile("avatar")
	var filename string
	if err == nil { // если аватарка прикреплена
		filename, err = h.UserUcase.UploadAvatar(ctx, file)
		if err != nil {
			return models.SendResponseWithError(c, err)
		}
	}

	profileEdits.Avatar = filename
	user := c.Get("User")

	if file != nil {
		if err != nil {
			requestError := errors.BadRequestError(err.Error())
			logger.DeliveryLevel().ErrorLog(ctx, requestError)
			return models.SendResponseWithError(c, requestError)
		}

		filename, err := h.UserUcase.UploadAvatar(ctx, file)
		if err != nil {
			return models.SendResponseWithError(c, err)
		}

		profileEdits.Avatar = filename
	}

	err = h.UserUcase.Update(ctx, profileEdits, user.(models.User).Uid)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	return models.SendResponse(c, profileEdits)
}

func (h *Handler) CheckAuth(c echo.Context) error {
	ctx := models.GetContext(c)

	cookie, err := c.Cookie(config.SessionCookie)
	if err != nil {
		sendErr := errors.NewCustomError(http.StatusUnauthorized, "error with request data")
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}

	authResponse := new(models.Auth)

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
		authResponse.Name = restaurant.Title
		authResponse.Avatar = restaurant.Avatar
		authResponse.Role = config.RoleAdmin
		return models.SendResponse(c, authResponse)

	case config.RoleUser:
		user, err := h.UserUcase.GetByUid(ctx, sessionData.Id)
		if err != nil {
			sendErr := errors.NewCustomError(http.StatusUnauthorized, err.Error())
			logger.DeliveryLevel().ErrorLog(ctx, sendErr)
			return models.SendResponseWithError(c, sendErr)
		}
		authResponse.Name = user.Name
		authResponse.Avatar = user.Avatar
		authResponse.Role = config.RoleUser
		return models.SendResponse(c, authResponse)
	default:
		sendErr := errors.NewCustomError(http.StatusUnauthorized, "error with roles")
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}
}

func (h *Handler) Logout(c echo.Context) error {
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

	return models.SendResponse(c, "")
}
