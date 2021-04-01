package middleware

import (
	"github.com/labstack/echo/v4/middleware"
)

var PanicConfig = middleware.RecoverWithConfig(middleware.RecoverConfig{
	StackSize: 1 << 1,
})
