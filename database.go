package main

type user struct {
	name string
	email string
	password string
	number string
}

type dish struct {
	name string
	cost int
	description string
}

type restourant struct {
	name string
	dishes []dish
	deliveryCost int
}

