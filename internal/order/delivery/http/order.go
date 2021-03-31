package http

import (
	"github.com/borscht/backend/internal/order"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	OrderUcase order.OrderUsecase
}

func NewOrderHandler(orderUcase order.OrderUsecase) order.OrderHandler {
	handler := &Handler{
		OrderUcase: orderUcase,
	}

	return handler
}

func (h Handler) Create(c echo.Context) error {
	panic("implement me")
}

func (h Handler) GetUserOrder(c echo.Context) error {
	panic("implement me")
}
