package local

// Это все временное, удалиться когда будет другое бд

import (
	"math/rand"
	"strconv"
	"sync"

	"github.com/borscht/backend/internal/models"
)

type LocalBd struct {
	model Model
}

var instance *LocalBd
var once sync.Once

func GetInstance() Database {
	once.Do(func() {
		instance = new(LocalBd)
		instance.initData()
		instance.initRest()
	})

	return instance
}

func (bd *LocalBd) GetModels() *Model {
	return &bd.model
}

var (
	restaurantCount = 10
	dishCostRange   = 1000
	dishWeightRange = 700
)

var dishNames = [...]string{"Макарошки с пюрешкой", "Цезарь", "Куриные котлетки", "Ольвьешечка",
	"Макарошки с котлеткой", "Ролл Калифорния", "Хлеб", "Пiво светлое", "Пiво темное"}

func createDish() models.Dish {
	dish := models.Dish{}
	dish.Name = dishNames[rand.Intn(len(dishNames))]
	dish.Price = rand.Intn(dishCostRange)
	dish.Weight = rand.Intn(dishWeightRange)
	dish.Description = "delicious"
	dish.Image = "static/food.jpg"
	return dish
}

func (bd *LocalBd) initRest() {
	for i := 0; i < restaurantCount; i++ {
		res := models.Restaurant{}
		res.DeliveryCost = restaurantCount * i
		res.Name = "Restaurant #" + strconv.Itoa(i)
		res.ID = i
		res.Rating = float64(i % 5)
		res.AvgCheck = i * 150
		res.Description = "yum"
		res.DeliveryTime = i * 15

		for j := 0; j < rand.Intn(10); j++ {
			dish := createDish()
			dish.ID = j
			res.Dishes = append(res.Dishes, dish)
		}

		(*bd.model.Restaurants)[strconv.Itoa(i)] = res
	}
}

func (bd *LocalBd) initData() {
	Users := make([]models.User, 0)
	Sessions := make(map[string]string, 0)
	Restaurants := make(map[string]models.Restaurant, 0)

	bd.model = Model{Users: &Users, Restaurants: &Restaurants, Sessions: &Sessions}
}
