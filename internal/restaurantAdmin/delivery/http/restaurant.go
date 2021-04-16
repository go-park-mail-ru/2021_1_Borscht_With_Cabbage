package http

import (
	"net/http"
	"time"

	"github.com/borscht/backend/utils/validation"

	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/models"
	adminModel "github.com/borscht/backend/internal/restaurantAdmin"
	sessionModel "github.com/borscht/backend/internal/session"
	errors "github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
	"github.com/labstack/echo/v4"
)

type RestaurantHandler struct {
	RestaurantUsecase adminModel.AdminRestaurantUsecase
	SessionUcase      sessionModel.SessionUsecase
}

func NewRestaurantHandler(adminUCase adminModel.AdminRestaurantUsecase,
	sessionUcase sessionModel.SessionUsecase) adminModel.AdminRestaurantHandler {

	return &RestaurantHandler{
		RestaurantUsecase: adminUCase,
		SessionUcase:      sessionUcase,
	}
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
		Name:     config.SessionCookie,
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
		sendErr := errors.BadRequestError("error with request data")
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}

	if err := validation.ValidateRestRegistration(*newRestaurant); err != nil {
		return models.SendResponseWithError(c, err)
	}

	responseRestaurant, err := a.RestaurantUsecase.CreateRestaurant(ctx, *newRestaurant)
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
		return models.SendResponseWithError(c, err)
	}

	existingRest, err := a.RestaurantUsecase.CheckRestaurantExists(ctx, *newRest)
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
