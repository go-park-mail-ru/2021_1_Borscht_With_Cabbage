package repository

import (
	"context"

	"github.com/borscht/backend/database/local"
	"github.com/borscht/backend/internal/models"
	restModel "github.com/borscht/backend/internal/restaurant"
)

type restaurantRepo struct {
	db local.Database
}

func NewRestaurantRepo() restModel.RestaurantRepo {
	return &restaurantRepo{
		db: local.GetInstance(),
	}
}

func (r *restaurantRepo) GetVendor(ctx context.Context, limit, offset int) ([]models.RestaurantResponse, error) {
	var result []models.RestaurantResponse

	for _, rest := range *r.db.GetModels().Restaurants {
		if rest.ID >= offset && rest.ID < offset+limit {
			restaurant := models.RestaurantResponse{
				ID:           rest.ID,
				Name:         rest.Name,
				Rating:       rest.Rating,
				DeliveryCost: rest.DeliveryCost,
				AvgCheck:     rest.AvgCheck,
				Description:  rest.Description,
				DeliveryTime: rest.DeliveryTime,
			}
			result = append(result, restaurant)
		}
	}

	return result, nil
}

func (r *restaurantRepo) GetById(ctx context.Context, id string) (models.Restaurant, bool) {
	restaurant, isItExists := (*r.db.GetModels().Restaurants)[id]

	return restaurant, isItExists
}
