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
	SectionName string `json:"sectionName"`
	SectionId   int    `json:"sectionId"`
	Dishes      []Dish `json:"dishes"`
}
