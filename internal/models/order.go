package models

var StatusOrderAdded = "оформлено"
var StatusOrderCooking = "готовится"
var StatusOrderDelivering = "едет к вам"
var StatusOrderDone = "доставлен"

type Order struct {
	OID          int           `json:"orderID"`
	UID          int           `json:"user"`
	Restaurant   string        `json:"store"`
	Address      string        `json:"address"`
	OrderTime    string        `json:"orderTime"`
	DeliveryCost int           `json:"ship"`
	DeliveryTime string        `json:"deliveryTime"`
	Summary      int           `json:"summary"`
	Status       string        `json:"status"`
	Foods        []DishInOrder `json:"foods"`
}

type DishInOrder struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Price  int    `json:"price"`
	Number int    `json:"num"`
	Image  string `json:"image"`
}

type CreateOrder struct {
	Address string `json:"address"`
}

type DishToBasket struct {
	DishID     int  `json:"dishID"`
	IsPlus     bool `json:"isPlus"`
	SameBasket bool `json:"same"`
}
