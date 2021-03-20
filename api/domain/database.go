package domain

import (
	"math/rand"
	"strconv"
)

const (
	Host          = "http://89.208.197.150"
	Repository    = Host + ":5000/"
	DefaultAvatar = Repository + "static/avatar/stas.jpg"
	SessionCookie = "borscht_session"
)

type Dish struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Description string `json:"description"`
	Weight      int    `json:"weight"`
	Image       string `json:"image"`
}

var (
	restaurantCount = 10
	dishCostRange   = 1000
	dishWeightRange = 700
)

var dishNames = [...]string{"Макарошки с пюрешкой", "Цезарь", "Куриные котлетки", "Ольвьешечка",
	"Макарошки с котлеткой", "Ролл Калифорния", "Хлеб", "Пiво светлое", "Пiво темное"}

func createDish() Dish {
	dish := Dish{}
	dish.Name = dishNames[rand.Intn(len(dishNames))]
	dish.Price = rand.Intn(dishCostRange)
	dish.Weight = rand.Intn(dishWeightRange)
	dish.Description = "delicious"
	dish.Image = "static/food.jpg"
	return dish
}

func InitData(cc CustomContext) {
	for i := 0; i < restaurantCount; i++ {
		res := Restaurant{}
		res.DeliveryCost = restaurantCount * i
		res.Name = "Restaurant #" + strconv.Itoa(i)
		res.ID = i
		res.Rating = float64(i % 5)
		res.AvgCheck = i * 150
		res.Description = "yum"
		res.DeliveryTime = i * 15

		for i := 0; i < rand.Intn(10); i++ {
			dish := createDish()
			dish.ID = i
			res.Dishes = append(res.Dishes, dish)
		}

		(*cc.Restaurants)[strconv.Itoa(i)] = res
	}
}
