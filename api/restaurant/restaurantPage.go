package restaurant

import (
	"backend/api"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetRestaurantPage(c echo.Context) error {
	cc := c.(*api.CustomContext)

	id := c.Param("id")

	restaurant, isItExists := (*cc.Restaurants)[id]
	if !isItExists {
		return c.String(http.StatusBadRequest, "error with request data")
	}

	return c.JSON(http.StatusOK, restaurant)
}