package websocket

import (
	"github.com/labstack/echo/v4"
)

type WebSocketHandler interface {
	Connect(c echo.Context) error
	GetKey(c echo.Context) error
}

type WebSocketUsecase interface {
}
