package restaurant

import (
	"backend/api"
	"fmt"
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

	fmt.Println(restaurant)

	return c.JSON(http.StatusOK, restaurant)
}

