package address

import (
	"context"

	"github.com/borscht/backend/internal/models"
	"github.com/labstack/echo/v4"
)

type AddressDelivery interface {
	UpdateMainAddress(c echo.Context) error
}

type AddressUsecase interface {
	UpdateMainAddress(ctx context.Context, address models.Address) error
}
