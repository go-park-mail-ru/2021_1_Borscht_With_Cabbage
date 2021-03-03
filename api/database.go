package api

type User struct {
	Name string
	Email string
	Password string
	Number string
}

type Dish struct {
	Name string
	Cost int
	Description string
}

type Restaurant struct {
	ID           int
	Name         string
	Dishes       []Dish
	DeliveryCost int
}
