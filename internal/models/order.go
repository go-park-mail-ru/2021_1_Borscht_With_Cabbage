package models

var StatusOrderAdded = "created"
var StatusOrderCooking = "cooking"
var StatusOrderDelivering = "delivering"
var StatusOrderDone = "done"

type Order struct {
	OID             int           `json:"orderID"`
	UID             int           `json:"user"`
	RID             int           `json:"restaurant"`
	UserName        string        `json:"userName"`
	UserPhone       string        `json:"userPhone"`
	Restaurant      string        `json:"store"`
	RestaurantImage string        `json:"restaurantImage"`
	Review          string        `json:"review"`
	Stars           int           `json:"stars"`
	Address         string        `json:"address"`
	OrderTime       string        `json:"orderTime"`
	DeliveryCost    int           `json:"ship"`
	DeliveryTime    string        `json:"deliveryTime"`
	Summary         int           `json:"summary"`
	Status          string        `json:"status"`
	Foods           []DishInOrder `json:"foods"`
}

type DishInOrder struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Price  int    `json:"price"`
	Number int    `json:"num"`
	Image  string `json:"image"`
	Weight int    `json:"weight"`
}

type CreateOrder struct {
	Address  string `json:"address"`
	BasketID int    `json:"basketID"`
}

type DishToBasket struct {
	DishID       int  `json:"dishID"`
	RestaurantID int  `json:"restaurantID"`
	IsPlus       bool `json:"isPlus"`
}

type SetNewStatus struct {
	OID          int    `json:"order"`
	Status       string `json:"status"`
	DeliveryTime string `json:"deliveryTime"`
	Restaurant   string
}

type SetNewReview struct {
	OID    int    `json:"oid"`
	Review string `json:"review"`
	Stars  int    `json:"stars"`
	User   int
}
