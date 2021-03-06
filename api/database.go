package api

type User struct {
	Name string
	Email string
	Password string
	Number string
	Avatar string
}

type Dish struct {
	name string
	cost int
	description string
}

type Restaurant struct {
	name string
	dishes []Dish
	deliveryCost int
}

