package api

type User struct {
	name string
	email string
	password string
	number string
}

func (info User) GetName() string {
	return info.name
}

func (info User) GetEmail() string {
	return info.email
}

func (info User) GetPassword() string {
	return info.password
}

func (info User) GetNumber() string {
	return info.number
}

func (info *User) SetName(name string) {
	info.name = name
}

func (info *User) SetEmail(email string) {
	info.email = email
}

func (info *User) SetPassword(password string) { // what for? ok let it be...
	info.password = password
}

func (info *User) SetNumber(number string) {
	info.number = number
}


type Dish struct {
	name string
	cost int
	description string
}

func (info Dish) GetName() string {
	return info.name
}

func (info Dish) GetCost() int {
	return info.cost
}

func (info Dish) GetDescription() string {
	return info.description
}

func (info *Dish) SetName(name string) {
	info.name = name
}

func (info *Dish) SetCost(cost int) {
	info.cost = cost
}

func (info *Dish) SetDescription(description string) {
	info.description = description
}


type Restourant struct {
	pk int
	name string
	dishes []Dish
	deliveryCost int
}

func (info Restourant) GetPK() int{
	return info.pk
}

func (info Restourant) GetName() string {
	return info.name
}

func (info Restourant) GetDishes() []Dish {
	return info.dishes
}

func (info Restourant) GetCost() int {
	return info.deliveryCost
}

func (info *Restourant) SetPK(pk int) {
	info.pk = pk
}

func (info *Restourant) SetName(name string) {
	info.name = name
}

func (info *Restourant) SetDishes(dishes []Dish) {
	info.dishes = dishes
}

func (info *Restourant) SetCost(cost int) {
	info.deliveryCost = cost
}

