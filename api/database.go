package api

import (
	"github.com/labstack/echo/v4"
)

type CustomContext struct {
	echo.Context
	Users       *[]User
	Restaurants *map[string]Restaurant // [id]RestaurantStruct
	Sessions    *map[string]string     // [session]user's phone number
}

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"number"`
	Avatar   string `json:"avatar"`
}

type Dish struct {
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Description string `json:"description"`
	Weight      string `json:"weight"`
}

type Restaurant struct {
	ID           int    `json:"id"`
	Name         string `json:"title"`
	Dishes       []Dish `json:"foods"`
	DeliveryCost int `json:"deliveryCost"`
}

type Session struct {
	Session string `json:"session"`
	Number  string `json:"number"`
}

var restaurantCount = 10

func InitData(cc CustomContext) {
	for i := 0; i < restaurantCount; i++ {
		res := Restaurant{}
		res.DeliveryCost = restaurantCount * i
		res.Name = "Restaurant #" + string(res.ID)
		res.ID = i

		(*cc.Restaurants)[string(res.ID)] = res
	}

	user := User{"Oleg", "oleg@mail.ru", "1111", "88005553535", ""}
	session := "olegssession"
	(*cc.Sessions)[session] = user.Phone
	*cc.Users = append(*cc.Users, user)
}
