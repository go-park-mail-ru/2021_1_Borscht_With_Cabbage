package middleware

import (
	"github.com/borscht/backend/internal/models"

	"github.com/labstack/echo/v4"
)

func ErrorHandler(err error, ctx echo.Context) {

	err = models.SendResponseWithError(ctx, err)
	// ctx.Logger().Errorf("error happened while processing request: %s", err)

	// switch err := errors.Cause(err); err.(type) {
	// case *echo.HTTPError:
	// 	ctx.JSON(err.(*echo.HTTPError).Code, struct {
	// 		Body string
	// 	}{Body: "internal"})
	// default:
	// 	ctx.JSON(500, struct {
	// 		Body string
	// 	}{Body: "internal"})
	// }

	// err = ctx.HTML(http.StatusInternalServerError, "internal")
	// if err != nil {
	// 	ctx.Logger().Errorf("failed to write 500 internal after error: %s", err)
	// }
}
