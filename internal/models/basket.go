package models

type DishInBasket struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Price  int    `json:"price"`
	Number int    `json:"num"`
	Image  string `json:"image"`
}

type BasketForUser struct {
	BID             int `json:"id"`
	UID             int
	Restaurant      string         `json:"restaurantName"`
	RestaurantImage string         `json:"restaurantImage"`
	RID             int            `json:"restaurantID"`
	DeliveryCost    int            `json:"deliveryPrice"`
	Summary         int            `json:"totalPrice"`
	Foods           []DishInBasket `json:"foods"`
	Address         Address        `json:"address"`
}
