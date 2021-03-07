package api

import "github.com/labstack/echo/v4"

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
	Dishes       []Dish
	DeliveryCost int `json:"deliveryCost"`
}

type Session struct {
	Session string `json:"session"`
	Number  string `json:"number"`
}
