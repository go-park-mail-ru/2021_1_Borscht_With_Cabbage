package api

import (
	errors "backend/models"
	"github.com/labstack/echo/v4"
	"math/rand"
	"net/http"
	"strconv"
)

const (
	Host = "http://89.208.197.150"
	Repository = Host + ":5000/"
	DefaultAvatar = Repository + "static/avatar/stas.jpg"
	SessionCookie = "borscht_session"
)

type CustomContext struct {
	echo.Context
	Users       *[]User
	Restaurants *map[string]Restaurant // [id]RestaurantStruct
	Sessions    *map[string]string     // [session]user's phone number
}

type message struct {
	Code int `json:"code"`
	Data interface{} `json:"data"`
}

func (c *CustomContext) SendOK(data interface{}) error {
	return c.JSON(http.StatusOK, message{200, data})
}

func (c *CustomContext) SendERR(err error) error {
	// проверяем можно ли преобразовать в кастомную ошибку
	if customErr, ok := err.(*errors.CustomError); ok {
		return c.JSON(http.StatusOK, customErr.SendError)
	}

	customErr := errors.FailServer(err.Error())
	return c.JSON(http.StatusOK, customErr.SendError)
}


type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"number"`
	Avatar   string `json:"avatar"`
}

type Dish struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Description string `json:"description"`
	Weight      int    `json:"weight"`
	Image       string `json:"image"`
}

type Restaurant struct {
	ID           int     `json:"id"`
	AvgCheck     int     `json:"cost"`
	Name         string  `json:"title"`
	DeliveryTime int     `json:"time"`
	Description  string  `json:"description"`
	Dishes       []Dish  `json:"foods"`
	DeliveryCost int     `json:"deliveryCost"`
	Rating       float64 `json:"rating"`
}

type Session struct {
	Session string `json:"session"`
	Number  string `json:"number"`
}

var (
	restaurantCount = 10
	dishCostRange = 1000
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
