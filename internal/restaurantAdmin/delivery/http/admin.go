package http

import (
	"net/http"
	"time"

	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/models"
	adminModel "github.com/borscht/backend/internal/restaurantAdmin"
	sessionModel "github.com/borscht/backend/internal/session"
	errors "github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
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

func (a AdminHandler) Update(c echo.Context) error {
	ctx := models.GetContext(c)

	restaurant := new(models.RestaurantUpdate)
	if err := c.Bind(restaurant); err != nil {
		sendErr := errors.BadRequestError(err.Error())
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}

	responseRestaurant, err := a.AdminUsecase.Update(ctx, *restaurant)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	return models.SendResponse(c, *responseRestaurant)
}

func (a AdminHandler) GetAllDishes(c echo.Context) error {
	ctx := models.GetContext(c)

	response, err := a.AdminUsecase.GetAllDishes(ctx)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	return models.SendResponse(c, response)
}

func (a AdminHandler) UpdateDish(c echo.Context) error {
	ctx := models.GetContext(c)

	updateDish := new(models.Dish)
	if err := c.Bind(updateDish); err != nil {
		sendErr := errors.BadRequestError(err.Error())
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}

	response, err := a.AdminUsecase.UpdateDish(ctx, *updateDish)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	return models.SendResponse(c, *response)
}

func (a AdminHandler) DeleteDish(c echo.Context) error {
	ctx := models.GetContext(c)

	idDish := new(models.DishDelete)
	if err := c.Bind(idDish); err != nil {
		sendErr := errors.BadRequestError(err.Error())
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}

	err := a.AdminUsecase.DeleteDish(ctx, idDish.ID)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	return models.SendResponse(c, nil)
}

func (a AdminHandler) AddDish(c echo.Context) error {
	ctx := models.GetContext(c)

	newDish := new(models.Dish)
	if err := c.Bind(newDish); err != nil {
		sendErr := errors.BadRequestError(err.Error())
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}

	response, err := a.AdminUsecase.AddDish(ctx, *newDish)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	return models.SendResponse(c, *response)
}

func setResponseCookie(c echo.Context, session string) {
	sessionCookie := http.Cookie{
		Expires:  time.Now().Add(24 * time.Hour),
		Name:     config.SessionCookie,
		Value:    session,
		HttpOnly: true,
		Path:     "/",
	}
	c.SetCookie(&sessionCookie)
}

func (a AdminHandler) Create(c echo.Context) error {
	ctx := models.GetContext(c)
	newRestaurant := new(models.Restaurant)
	if err := c.Bind(newRestaurant); err != nil {
		sendErr := errors.NewCustomError(http.StatusUnauthorized, "error with request data")
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}

	responseRestaurant, err := a.AdminUsecase.Create(ctx, *newRestaurant)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	sessionInfo := models.SessionInfo{
		Id:   responseRestaurant.ID,
		Role: config.RoleAdmin,
	}
	session, err := a.SessionUcase.Create(ctx, sessionInfo)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	setResponseCookie(c, session)

	response := models.SuccessRestaurantResponse{Restaurant: *responseRestaurant, Role: config.RoleAdmin} // TODO убрать config отсюда
	return models.SendResponse(c, response)
}

func (a AdminHandler) Login(c echo.Context) error {
	ctx := models.GetContext(c)
	newRest := new(models.RestaurantAuth)

	if err := c.Bind(newRest); err != nil {
		sendErr := errors.NewCustomError(http.StatusUnauthorized, "error with request data")
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}

	existingRest, err := a.AdminUsecase.CheckRestaurantExists(ctx, *newRest)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	sessionInfo := models.SessionInfo{
		Id:   existingRest.ID,
		Role: config.RoleAdmin,
	}
	session, err := a.SessionUcase.Create(ctx, sessionInfo)

	if err != nil {
		return models.SendResponseWithError(c, err)
	}
	setResponseCookie(c, session)

	response := models.SuccessRestaurantResponse{Restaurant: *existingRest, Role: config.RoleAdmin}
	return models.SendResponse(c, response)
}

func (a AdminHandler) GetUserData(c echo.Context) error {
	ctx := models.GetContext(c)
	sendErr := errors.NewCustomError(500, "не реализовано")
	logger.DeliveryLevel().ErrorLog(ctx, sendErr)
	return models.SendResponseWithError(c, sendErr) // TODO
}

func (a AdminHandler) EditProfile(c echo.Context) error {
	ctx := models.GetContext(c)
	sendErr := errors.NewCustomError(500, "не реализовано")
	logger.DeliveryLevel().ErrorLog(ctx, sendErr)
	return models.SendResponseWithError(c, sendErr) // TODO
}
