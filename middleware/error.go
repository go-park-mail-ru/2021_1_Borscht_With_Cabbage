package middleware

import (
	"github.com/borscht/backend/internal/models"
	errors "github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"

	"github.com/labstack/echo/v4"
)

func ErrorHandler(err error, c echo.Context) {

	ctx := models.GetContext(c)

	custErr := errors.FailServerError("PANIC" + err.Error())
	logger.ResponseLevel().ErrorLog(ctx, custErr)
	err = models.SendResponseWithError(c, custErr)
}
