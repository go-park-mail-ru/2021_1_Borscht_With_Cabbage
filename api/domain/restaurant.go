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

type RestaurantResponse struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Rating       float64 `json:"rating"`
	DeliveryTime int     `json:"time"`
	AvgCheck     int     `json:"cost"`
	DeliveryCost int     `json:"deliveryCost"`
}

type RestaurantUsecase interface {
	GetSlice(ctx *CustomContext, limit, offset int) ([]RestaurantResponse, error)
	GetById(ctx *CustomContext, id string) (Restaurant, bool)
}

type RestaurantRepo interface {
	GetSlice(ctx *CustomContext, limit, offset int) ([]RestaurantResponse, error)
	GetById(ctx *CustomContext, id string) (Restaurant, bool)
}
