package api

var Users []User
var Restaurants []Restaurant
var Sessions []Session

const restaurantCount = 15

// заполнение хоть чем-то
func InitData()  {
	for i := 0; i < restaurantCount; i++ {
		res := Restaurant{}
		res.DeliveryCost = restaurantCount * i
		res.Name = "My rest"
		res.ID = i

		Restaurants = append(Restaurants, res)
	}
}
