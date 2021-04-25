package models

type RestaurantAuth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RestaurantWithDishes struct {
	RestaurantInfo
	Sections []Section `json:"sections"`
	Dishes   []Dish    `json:"foods"`
}

type SuccessRestaurantResponse struct {
	RestaurantInfo
	Role string `json:"role"`
}

type RestaurantReview struct {
	Review   string `json:"review"`
	Stars    int    `json:"stars"`
	Time     string `json:"deliveryTime"`
	UserName string `json:"user"`
}

type RestaurantInfo struct {
	ID                int    `json:"id"`
	AdminEmail        string `json:"email"`
	AdminPhone        string `json:"number"`
	AdminPassword     string `json:"password"`
	AdminHashPassword []byte
	AvgCheck          int     `json:"cost"`
	Title             string  `json:"title"`
	Description       string  `json:"description"`
	DeliveryCost      int     `json:"deliveryCost"`
	Rating            float64 `json:"rating"`
	Avatar            string  `json:"avatar"`
}

type RestaurantImageResponse struct {
	Filename string `json:"avatar"`
}

type RestaurantUpdateData struct {
	ID                int    `json:"id"`
	AdminEmail        string `json:"email"`
	AdminPhone        string `json:"number"`
	AdminPassword     string `json:"password"`
	AdminHashPassword []byte
	Title             string `json:"title"`
	Description       string `json:"description"`
	DeliveryCost      int    `json:"deliveryCost"`
}

type CheckRestaurantExists struct {
	Email         string
	Number        string
	Name          string
	CurrentRestId int
}

type DeleteSuccess struct {
	ID int `json:"id"`
}
