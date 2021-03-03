package page

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"backend/api"
)

var restourants []api.Restourant

type restourantJSON struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	DeliveryCost int    `json:"deliveryCost"`
}

func initRestourants() {
	restourants = restourants[0:0]

	for i := 0; i < 10; i++ {
		res := api.Restourant{}
		res.SetCost(10)
		res.SetName("My rest")
		res.SetPK(i)

		restourants = append(restourants, res)
	}
}

// загрузка главной страницы с ресторанами
func MainPage(c echo.Context) error {
	return c.String(http.StatusOK, "It will be the main page")
}

// загрузка списка рестаранов
func GetVendor(c echo.Context) error {
	var quotaVendor struct{
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
	}

	err := echo.QueryParamsBinder(c).
		Int("limit", &quotaVendor.Limit).
		Int("offset", &quotaVendor.Offset).
		BindError()

	if err != nil {
		return err
	}

	result := getRestaurant(quotaVendor.Limit, quotaVendor.Offset)
	return c.JSON(http.StatusOK, result)
}

// в будущем здесь будет поход в базу данных
func getRestaurant(limit, offset int) []restourantJSON {
	initRestourants()

	var result []restourantJSON

	for _, val := range restourants {
		if val.GetPK() >= offset && val.GetPK() < offset + limit {
			restaurant := restourantJSON {
				ID: val.GetPK(),
				Name: val.GetName(),
				DeliveryCost: val.GetCost(),
			}
			result = append(result, restaurant)
		}
	}

	return result
}
