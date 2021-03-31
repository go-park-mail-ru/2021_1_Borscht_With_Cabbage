package models

var StatusOrderAdded = "оформлено"
var StatusOrderCooking = "готовится"
var StatusOrderDelivering = "едет к вам"
var StatusOrderDone = "доставлен"

type Order struct {
	OID          int    `json:"orderID"`
	UID          int    `json:"user"`
	Restaurant   string `json:"store"`
	Address      int    `json:"address"`
	OrderTime    string `json:"orderTime"`
	DeliveryCost int    `json:"ship"`
	DeliveryTime string `json:"deliveryTime"`
	Summary      string `json:"description"`
	Status       string `json:"status"`
	Foods        []Dish `json:"foods"`
}

type CreateOrder struct {
	Address string `json:"address"`
}
