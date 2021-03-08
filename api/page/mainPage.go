package page

import (
	"backend/api"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

// загрузка списка рестаранов
func GetVendor(c echo.Context) error {
	Limit, errLimit := strconv.Atoi(c.QueryParam("limit"))
	Offset, errOffset := strconv.Atoi(c.QueryParam("offset"))

	if errLimit != nil || errOffset != nil {
		return c.JSON(http.StatusOK, struct {
			Error string `json:"error"`
		}{
			Error: "400",
		})
	}

	cc := c.(*api.CustomContext)

	result := getRestaurant(Limit, Offset, cc)
	return c.JSON(http.StatusOK, result)
}

// в будущем здесь будет поход в базу данных
func getRestaurant(limit, offset int, context *api.CustomContext) []api.Restaurant {
	var result []api.Restaurant

	for _, restaurant := range *context.Restaurants {
		if restaurant.ID >= offset && restaurant.ID < offset+limit {

			result = append(result, restaurant)
		}
	}

	return result
}
