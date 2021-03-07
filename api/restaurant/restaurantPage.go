package restaurant

import (
	"backend/api"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func RestaurantPage(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	for _, restaurant := range api.Restaurants {
		if restaurant.ID == id {
			restaurant, err := json.Marshal(restaurant)
			if err != nil {
				return err
			}
			return c.String(http.StatusOK, string(restaurant))
		}
	}

	return c.String(http.StatusNotImplemented, "restaurant not found")
}