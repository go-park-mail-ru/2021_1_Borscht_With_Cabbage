package order

import "github.com/labstack/echo/v4"

type OrderHandler interface {
	Create(c echo.Context) error
	GetUserOrder(c echo.Context) error
}

type OrderUsecase interface {
	Create(uid int) (string, error)
	GetUserOrder(session string)
}

type OrderRepo interface {
	Create(session string, uid int) error
	GetUserOrder(uid int)
}
