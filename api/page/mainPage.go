package page

import (
	"backend/api"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type restaurantResponse struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	DeliveryCost int    `json:"deliveryCost"`
}

func initRestaurants() []api.Restaurant {
	restaurants := make([]api.Restaurant, 0, 10)

	for i := 0; i < 10; i++ {
		res := api.Restaurant{}
		res.DeliveryCost = 10
		res.Name = "My rest"
		res.ID = i

		restaurants = append(restaurants, res)
	}

	return restaurants
}

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

	result := getRestaurant(Limit, Offset)
	return c.JSON(http.StatusOK, result)
}

// в будущем здесь будет поход в базу данных
func getRestaurant(limit, offset int) []restaurantResponse {
	var result []restaurantResponse

	restaurants := initRestaurants()

	for _, val := range restaurants {
		if val.ID >= offset && val.ID < offset+limit {
			restaurant := restaurantResponse{
				ID:           val.ID,
				Name:         val.Name,
				DeliveryCost: val.DeliveryCost,
			}
			result = append(result, restaurant)
		}
	}

	return result
}
