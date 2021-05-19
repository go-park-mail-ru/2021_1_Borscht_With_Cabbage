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
	Address           Address `json:"address"`
	ID                int     `json:"id"`
	AdminEmail        string  `json:"email"`
	AdminPhone        string  `json:"number"`
	AdminPassword     string  `json:"password"`
	AdminHashPassword []byte
	AvgCheck          int     `json:"cost"`
	Title             string  `json:"title"`
	Description       string  `json:"description"`
	DeliveryCost      int     `json:"deliveryCost"`
	DeliveryTime      int     `json:"deliveryTime"`
	Rating            float64 `json:"rating"`
	Avatar            string  `json:"avatar"`
}

type RestaurantImageResponse struct {
	Filename string `json:"avatar"`
}

type RestaurantUpdateData struct {
	Address           Address `json:"address"`
	ID                int     `json:"id"`
	AdminEmail        string  `json:"email"`
	AdminPhone        string  `json:"number"`
	AdminPassword     string  `json:"password"`
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

type RestaurantRequest struct {
	Limit         int      `json:"limit"`
	Offset        int      `json:"offset"`
	Categories    []string `json:"categories"`
	Time          int      `json:"time"`
	Receipt       int      `json:"receipt"`
	Rating        float64  `json:"rating"`
	LatitudeUser  string   `json:"latitude"`
	LongitudeUser string   `json:"longitude"`
	Address       bool     `json:"address"`
}

type Categories struct {
	CategoriesID []string `json:"categories"`
}

type RecommendationsParams struct {
	Id            int    `json:"id"`
	LatitudeUser  string `json:"latitude"`
	LongitudeUser string `json:"longitude"`
}
