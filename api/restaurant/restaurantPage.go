package restaurant

import (
	"backend/api"
	"encoding/json"
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

	restaurantToOutput, err := json.Marshal(restaurant)
	if err != nil {
		return c.String(http.StatusNotImplemented, "error while marshalling result")
	}
	return c.String(http.StatusOK, string(restaurantToOutput))
}