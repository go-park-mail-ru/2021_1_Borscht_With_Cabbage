package restaurant

import (
	"backend/api"
	"backend/api/auth"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

var Restaurants []api.Restaurant

func RestaurantPage(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	session, err := c.Cookie("session")
	if err != nil {
		return err
	}
	//auth.CheckSession()

	for i, restaurant := range Restaurants {
		if i == id {
			restaurant, err := json.Marshal(restaurant)
			if err != nil {
				return err
			}
			return c.String(http.StatusOK, string(restaurant))
		}
	}

	return c.String(http.StatusNotImplemented, "restaurant not found")
}