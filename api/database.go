package api

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Number   string `json:"number"`
}

type Dish struct {
	Name        string `json:"name"`
	Price       int `json:"price"`
	Description string `json:"description"`
	Weight string `json:"weight"`
}

type Restaurant struct {
	ID           int `json:"id"`
	Name         string `json:"title"`
	Dishes       []Dish `json:"foods"`
	DeliveryCost int `json:"deliveryCost"`
}

type Session struct {
	Session string `json:"session"`
	Number  string `json:"number"`
}
