package middleware

import (
	"github.com/borscht/backend/internal/models"
	errors "github.com/borscht/backend/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var PanicConfig = middleware.RecoverWithConfig(middleware.RecoverConfig{
	StackSize: 1 << 1,
})

func ErrorHandler(err error, c echo.Context) {

	ctx := models.GetContext(c)

	custErr := errors.FailServer(ctx, "PANIC: "+err.Error())
	err = models.SendResponseWithError(c, custErr)
}
