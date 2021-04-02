package models

// TODO много одинаковых полей, надо объеденить
type Dish struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Description string `json:"description"`
	Weight      int    `json:"weight"`
	Image       string `json:"image"`
}

type RestaurantAuth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Restaurant struct {
	ID            int     `json:"id"`
	AdminEmail    string  `json:"email"`
	AdminPhone    string  `json:"number"`
	AdminPassword string  `json:"password"`
	AvgCheck      int     `json:"cost"`
	Name          string  `json:"title"`
	Description   string  `json:"description"`
	Dishes        []Dish  `json:"foods"`
	DeliveryCost  int     `json:"deliveryCost"`
	Rating        float64 `json:"rating"`
	Avatar        string  `json:"avatar"`
}

type RestaurantResponse struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Rating       float64 `json:"rating"`
	DeliveryTime int     `json:"time"`
	AvgCheck     int     `json:"cost"`
	DeliveryCost int     `json:"deliveryCost"`
	Avatar       string  `json:"avatar"`
}

type DishResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

type CheckRestaurantExists struct {
	Email         string
	Number        string
	Name          string
	CurrentRestId int
}

type CheckDishExists struct {
	Name         string
	RestaurantId int
}
