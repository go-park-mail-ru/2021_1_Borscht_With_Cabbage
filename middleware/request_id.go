package middleware

import (
	"fmt"
	"math/rand"

	"github.com/labstack/echo/v4"
)

func RequestIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestID := fmt.Sprintf("%016x", rand.Int())[:10]

		c.Set("request_id", requestID)
		return next(c)
	}
}
