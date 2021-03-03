package main
// мб можно и без геттеров-сеттеров, но потом будет проще, если, когда мы подключим бд, мы в эти функции
// запишем логику работы с базой данных

type user struct {
	name string
	email string
	password string
	number string
}

func (info user) GetName() string {
	return info.name
}

func (info user) GetEmail() string {
	return info.email
}

func (info user) GetPassword() string {
	return info.password
}

func (info user) GetNumber() string {
	return info.number
}

func (info user) SetName(name string) {
	info.name = name
}

func (info user) SetEmail(email string) {
	info.email = email
}

func (info user) SetPassword(password string) { // what for? ok let it be...
	info.password = password
}

func (info user) SetNumber(number string) {
	info.number = number
}


type dish struct {
	name string
	cost int
	description string
}

func (info dish) GetName() string {
	return info.name
}

func (info dish) GetCost() int {
	return info.cost
}

func (info dish) GetDescription() string {
	return info.description
}

func (info dish) SetName(name string) {
	info.name = name
}

func (info dish) SetCost(cost int) {
	info.cost = cost
}

func (info dish) SetDescription(description string) {
	info.description = description
}


type restourant struct {
	name string
	dishes []dish
	deliveryCost int
}

func (info restourant) GetName() string {
	return info.name
}

func (info restourant) GetDishes() []dish {
	return info.dishes
}

func (info restourant) GetCost() int {
	return info.deliveryCost
}

func (info restourant) SetName(name string) {
	info.name = name
}

func (info restourant) SetDishes(dishes []dish) {
	info.dishes = dishes
}

func (info restourant) SetCost(cost int) {
	info.deliveryCost = cost
}

