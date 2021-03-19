package domain

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

type RestaurantUsecase interface {

}
