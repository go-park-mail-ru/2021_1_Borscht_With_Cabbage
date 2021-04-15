package models

type Dish struct {
	ID          int    `json:"id"`
	Restaurant  int    `json:"restaurant"`
	Section     int    `json:"section"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Description string `json:"description"`
	Weight      int    `json:"weight"`
	Image       string `json:"image"`
}

type DishDelete struct {
	ID int `json:"id"`
}

type DishImageResponse struct {
	ID       int    `json:"id"`
	Filename string `json:"filename"`
}

type CheckDishExists struct {
	Id           int
	Name         string
	RestaurantId int
}
