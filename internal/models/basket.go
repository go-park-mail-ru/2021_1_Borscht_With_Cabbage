package models

type DishInBasket struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Price  int    `json:"price"`
	Number int    `json:"num"`
	Image  string `json:"image"`
}

type BasketForUser struct {
	BID          int            `json:"id"`
	Restaurant   string         `json:"restaurantName"`
	RID          int            `json:"restaurantID"`
	DeliveryCost int            `json:"deliveryPrice"`
	Summary      int            `json:"totalPrice"`
	Foods        []DishInBasket `json:"foods"`
}
