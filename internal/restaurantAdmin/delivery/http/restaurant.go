package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/borscht/backend/internal/services/auth"

	"github.com/borscht/backend/utils/validation"

	"github.com/borscht/backend/configProject"
	"github.com/borscht/backend/internal/models"
	adminModel "github.com/borscht/backend/internal/restaurantAdmin"
	errors "github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
	"github.com/labstack/echo/v4"
)

type RestaurantHandler struct {
	RestaurantUsecase adminModel.AdminRestaurantUsecase
	AuthService       auth.ServiceAuth
}

func NewRestaurantHandler(adminUCase adminModel.AdminRestaurantUsecase, authService auth.ServiceAuth) adminModel.AdminRestaurantHandler {
	return &RestaurantHandler{
		RestaurantUsecase: adminUCase,
		AuthService:       authService,
	}
}

func (a RestaurantHandler) AddCategories(c echo.Context) error {
	ctx := models.GetContext(c)

	nameCategories := new(models.Categories)
	if err := c.Bind(nameCategories); err != nil {
		sendErr := errors.BadRequestError(err.Error())
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}

	err := a.RestaurantUsecase.AddCategories(ctx, *nameCategories)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	// TODO: подумать что должен вернуть бэк
	return models.SendResponse(c, nil)
}

func (a RestaurantHandler) UpdateRestaurantData(c echo.Context) error {
	ctx := models.GetContext(c)

	restaurant := new(models.RestaurantUpdateData)
	if err := c.Bind(restaurant); err != nil {
		sendErr := errors.BadRequestError(err.Error())
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}

	responseRestaurant, err := a.RestaurantUsecase.UpdateRestaurantData(ctx, *restaurant)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	return models.SendResponse(c, responseRestaurant)
}

func setResponseCookie(c echo.Context, session string) {
	sessionCookie := http.Cookie{
		Expires:  time.Now().Add(24 * time.Hour),
		Name:     configProject.SessionCookie,
		Value:    session,
		HttpOnly: true,
		Path:     "/",
	}
	c.SetCookie(&sessionCookie)
}

func (a RestaurantHandler) CreateRestaurant(c echo.Context) error {
	ctx := models.GetContext(c)
	newRestaurant := new(models.RestaurantInfo)
	if err := c.Bind(newRestaurant); err != nil {
		fmt.Println(err)
		sendErr := errors.BadRequestError("error with request data")
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}

	if err := validation.ValidateRestRegistration(*newRestaurant); err != nil {
		return models.SendResponseWithError(c, err)
	}

	responseRestaurant, err := a.AuthService.CreateRestaurant(ctx, *newRestaurant)
	logger.DeliveryLevel().DebugLog(ctx, logger.Fields{"restaurant auth": responseRestaurant})
	if err != nil {
		return models.SendResponseWithError(c, err)
	}
	err = a.RestaurantUsecase.AddAddress(ctx, responseRestaurant.ID, newRestaurant.Address)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	sessionInfo := models.SessionInfo{
		Id:   responseRestaurant.ID,
		Role: configProject.RoleAdmin,
	}

	session, err := a.AuthService.CreateSession(ctx, sessionInfo)

	if err != nil {
		return models.SendResponseWithError(c, err)
	}
	setResponseCookie(c, session)

	return models.SendResponse(c, responseRestaurant)
}

func (a RestaurantHandler) Login(c echo.Context) error {
	ctx := models.GetContext(c)
	newRest := new(models.RestaurantAuth)

	if err := c.Bind(newRest); err != nil {
		sendErr := errors.BadRequestError("error with request data")
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}

	if err := validation.ValidateSignIn(newRest.Login, newRest.Password); err != nil {
		logger.DeliveryLevel().ErrorLog(ctx, err)
		return models.SendResponseWithError(c, err)
	}

	existingRest, err := a.AuthService.CheckRestaurantExists(ctx, *newRest)
	if err != nil {
		logger.DeliveryLevel().ErrorLog(ctx, err)
		return models.SendResponseWithError(c, err)
	}

	address, err := a.RestaurantUsecase.GetAddress(ctx, existingRest.ID)
	if err != nil {
		logger.DeliveryLevel().ErrorLog(ctx, err)
		return models.SendResponseWithError(c, err)
	}
	existingRest.Address = *address

	sessionInfo := models.SessionInfo{
		Id:   existingRest.ID,
		Role: configProject.RoleAdmin,
	}
	session, err := a.AuthService.CreateSession(ctx, sessionInfo)

	if err != nil {
		return models.SendResponseWithError(c, err)
	}
	setResponseCookie(c, session)

	return models.SendResponse(c, existingRest)
}

func (a RestaurantHandler) GetUserData(c echo.Context) error {
	ctx := models.GetContext(c)
	sendErr := errors.FailServerError("не реализовано")
	logger.DeliveryLevel().ErrorLog(ctx, sendErr)
	return models.SendResponseWithError(c, sendErr) // TODO
}

func (a RestaurantHandler) UploadRestaurantImage(c echo.Context) error {
	ctx := models.GetContext(c)

	file, err := c.FormFile("avatar")
	if err != nil {
		requestError := errors.BadRequestError(err.Error())
		logger.DeliveryLevel().ErrorLog(ctx, requestError)
		return models.SendResponseWithError(c, requestError)
	}

	response, err := a.RestaurantUsecase.UploadRestaurantImage(ctx, file)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	return models.SendResponse(c, response)
}
