package models

type Section struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Restaurant int    `json:"restaurant"`
}

type CheckSectionExists struct {
	Name         string
	RestaurantId int
	ID           int
}

type SectionWithDishes struct {
	SectionName string `json:"name"`
	SectionId   int    `json:"id"`
	Dishes      []Dish `json:"dishes"`
}

type ArraySectionWithDishes struct {
	Section []SectionWithDishes `json:""`
}
