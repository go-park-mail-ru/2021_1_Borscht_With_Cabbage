package http

import (
	"strconv"

	"github.com/borscht/backend/internal/models"
	adminModel "github.com/borscht/backend/internal/restaurantAdmin"
	errors "github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
	"github.com/labstack/echo/v4"
)

type DishHandler struct {
	DishUsecase adminModel.AdminDishUsecase
}

func NewDishHandler(dishUCase adminModel.AdminDishUsecase) adminModel.AdminDishHandler {
	return &DishHandler{
		DishUsecase: dishUCase,
	}
}

func (a DishHandler) GetAllDishes(c echo.Context) error {
	ctx := models.GetContext(c)

	result, err := a.DishUsecase.GetAllDishes(ctx)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	response := make([]models.Response, 0)
	for _, val := range result {
		response = append(response, &val)
	}
	return models.SendResponse(c, response...)
}

func (a DishHandler) UpdateDishData(c echo.Context) error {
	ctx := models.GetContext(c)

	updateDish := new(models.Dish)
	if err := c.Bind(updateDish); err != nil {
		sendErr := errors.BadRequestError(err.Error())
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}

	response, err := a.DishUsecase.UpdateDishData(ctx, *updateDish)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	return models.SendResponse(c, response)
}

func (a DishHandler) DeleteDish(c echo.Context) error {
	ctx := models.GetContext(c)

	idDish := new(models.DishDelete)
	if err := c.Bind(idDish); err != nil {
		sendErr := errors.BadRequestError(err.Error())
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}
	response, err := a.DishUsecase.DeleteDish(ctx, idDish.ID)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	return models.SendResponse(c, response)
}

func (a DishHandler) AddDish(c echo.Context) error {
	ctx := models.GetContext(c)

	newDish := new(models.Dish)
	if err := c.Bind(newDish); err != nil {
		sendErr := errors.BadRequestError(err.Error())
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}

	response, err := a.DishUsecase.AddDish(ctx, *newDish)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	return models.SendResponse(c, response)
}

func (a DishHandler) UploadDishImage(c echo.Context) error {
	ctx := models.GetContext(c)

	file, err := c.FormFile("image")
	if err != nil {
		requestError := errors.BadRequestError(err.Error())
		logger.DeliveryLevel().ErrorLog(ctx, requestError)
		return models.SendResponseWithError(c, requestError)
	}

	formParams, err := c.FormParams()
	if err != nil {
		requestError := errors.BadRequestError("invalid data form")
		logger.DeliveryLevel().ErrorLog(ctx, requestError)
		return models.SendResponseWithError(c, requestError)
	}

	logger.DeliveryLevel().InlineDebugLog(ctx, formParams.Get("id"))
	idDish, err := strconv.Atoi(formParams.Get("id"))
	if err != nil {
		requestError := errors.BadRequestError("invalid data in formdata")
		logger.DeliveryLevel().ErrorLog(ctx, requestError)
		return models.SendResponseWithError(c, requestError)
	}

	response, err := a.DishUsecase.UploadDishImage(ctx, file, idDish)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	return models.SendResponse(c, response)
}
