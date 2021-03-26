package http

import (
	"fmt"
	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/models"
	adminModel "github.com/borscht/backend/internal/restaurantAdmin"
	sessionModel "github.com/borscht/backend/internal/session"
	errors "github.com/borscht/backend/utils"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
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
	newRestaurant := new(models.Restaurant)
	if err := c.Bind(newRestaurant); err != nil {
		sendErr := errors.NewCustomError(http.StatusUnauthorized, "error with request data")
		return models.SendResponseWithError(c, sendErr)
	}

	rid, err := a.AdminUsecase.Create(*newRestaurant)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	session, err := a.SessionUcase.Create(rid, config.RoleAdmin)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	setResponseCookie(c, session)

	response := models.SuccessResponse{Name: newRestaurant.Name, Avatar: config.DefaultAvatar} // TODO убрать config отсюда
	return models.SendResponse(c, response)
}

func (a AdminHandler) Login(c echo.Context) error {
	newRest := new(models.RestaurantAuth)

	if err := c.Bind(newRest); err != nil {
		sendErr := errors.NewCustomError(http.StatusUnauthorized, "error with request data")
		return models.SendResponseWithError(c, sendErr)
	}

	fmt.Println(newRest)

	existingRest, err := a.AdminUsecase.CheckRestaurantExists(*newRest)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	fmt.Println(existingRest)

	session, err := a.SessionUcase.Create(existingRest.ID, config.RoleAdmin)

	if err != nil {
		return models.SendResponseWithError(c, err)
	}
	setResponseCookie(c, session)

	response := models.SuccessResponse{Name: existingRest.Name, Avatar: existingRest.Avatar}
	return models.SendResponse(c, response)
}

func (a AdminHandler) GetUserData(c echo.Context) error {
	panic("implement me") // TODO
}

func (a AdminHandler) EditProfile(c echo.Context) error {
	panic("implement me") // TODO
}
