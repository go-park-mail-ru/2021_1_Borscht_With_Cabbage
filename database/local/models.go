package local

import (
	"github.com/borscht/backend/internal/models"
)

type Model struct {
	Users       *[]models.User
	Restaurants *map[string]models.Restaurant // [id]RestaurantStruct
	Sessions    *map[string]string            // [session]user's phone number
}

type Database interface {
	GetModels() *Model
}
