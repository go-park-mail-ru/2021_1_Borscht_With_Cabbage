package page

import (
	"backend/api"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type RestaurantResponse struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	DeliveryCost int    `json:"deliveryCost"`
}

// загрузка списка ресторанов
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

	result := GetRestaurant(Limit, Offset, c)
	return c.JSON(http.StatusOK, result)
}

// в будущем здесь будет поход в базу данных
func GetRestaurant(limit, offset int, c echo.Context) []RestaurantResponse {
	cc := c.(*api.CustomContext)
	var result []RestaurantResponse

	for _, rest := range *cc.Restaurants {
		if rest.ID >= offset && rest.ID < offset+limit {
			restaurant := RestaurantResponse{
				ID:           rest.ID,
				Name:         rest.Name,
				DeliveryCost: rest.DeliveryCost,
			}
			result = append(result, restaurant)
		}
	}

	return result
}
