package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/models"
	adminModel "github.com/borscht/backend/internal/restaurantAdmin"
	sessionModel "github.com/borscht/backend/internal/session"
	errors "github.com/borscht/backend/utils"
	"github.com/labstack/echo/v4"
)

type AdminHandler struct {
	AdminUsecase adminModel.AdminUsecase
	SessionUcase sessionModel.SessionUsecase
}

func NewAdminHandler(adminUCase adminModel.AdminUsecase, sessionUcase sessionModel.SessionUsecase) adminModel.AdminHandler {
	return &AdminHandler{
		AdminUsecase: adminUCase,
		SessionUcase: sessionUcase,
	}
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

func (a AdminHandler) Create(c echo.Context) error {
	ctx := models.GetContext(c)
	newRestaurant := new(models.Restaurant)
	if err := c.Bind(newRestaurant); err != nil {
		sendErr := errors.NewCustomError(ctx, http.StatusUnauthorized, "error with request data")
		return models.SendResponseWithError(c, sendErr)
	}

	rid, err := a.AdminUsecase.Create(ctx, *newRestaurant)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	session, err := a.SessionUcase.Create(ctx, rid, config.RoleAdmin)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	setResponseCookie(c, session)

	response := models.SuccessResponse{Name: newRestaurant.Name, Avatar: config.DefaultAvatar, Role: config.RoleAdmin} // TODO убрать config отсюда
	return models.SendResponse(c, response)
}

func (a AdminHandler) Login(c echo.Context) error {
	ctx := models.GetContext(c)
	newRest := new(models.RestaurantAuth)

	if err := c.Bind(newRest); err != nil {
		sendErr := errors.NewCustomError(ctx, http.StatusUnauthorized, "error with request data")
		return models.SendResponseWithError(c, sendErr)
	}

	fmt.Println(newRest)

	existingRest, err := a.AdminUsecase.CheckRestaurantExists(ctx, *newRest)
	if err != nil {
		fmt.Println(err)
		return models.SendResponseWithError(c, err)
	}

	fmt.Println(existingRest)

	session, err := a.SessionUcase.Create(ctx, existingRest.ID, config.RoleAdmin)

	if err != nil {
		return models.SendResponseWithError(c, err)
	}
	setResponseCookie(c, session)

	response := models.SuccessResponse{Name: existingRest.Name, Avatar: existingRest.Avatar, Role: config.RoleAdmin}
	return models.SendResponse(c, response)
}

func (a AdminHandler) GetUserData(c echo.Context) error {
	panic("implement me") // TODO
}

func (a AdminHandler) EditProfile(c echo.Context) error {
	panic("implement me") // TODO
}
