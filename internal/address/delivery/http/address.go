package http

import (
	"github.com/borscht/backend/internal/address"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
	"github.com/labstack/echo/v4"
)

type AddressHandler struct {
	AddressUsecase address.AddressUsecase
}

func NewAddressHandler(addressUsecase address.AddressUsecase) address.AddressDelivery {
	return &AddressHandler{
		AddressUsecase: addressUsecase,
	}
}

func (h AddressHandler) UpdateMainAddress(c echo.Context) error {
	ctx := models.GetContext(c)
	logger.DeliveryLevel().InlineDebugLog(ctx, "address delivery")

	address := new(models.Address)
	if err := c.Bind(address); err != nil {
		sendErr := errors.BadRequestError(err.Error())
		logger.DeliveryLevel().ErrorLog(ctx, sendErr)
		return models.SendResponseWithError(c, sendErr)
	}

	err := h.AddressUsecase.UpdateMainAddress(ctx, *address)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	return models.SendResponse(c, nil)
}
