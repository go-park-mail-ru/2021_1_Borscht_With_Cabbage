package repository

import (
	"backend/api/domain"
)

type restaurantRepo struct {
}

func NewRestaurantRepo() domain.RestaurantRepo {
	return &restaurantRepo{}
}

func (r *restaurantRepo) GetSlice(ctx *domain.CustomContext, limit, offset int) ([]domain.RestaurantResponse, error) {
	var result []domain.RestaurantResponse

	for _, rest := range *ctx.Restaurants {
		if rest.ID >= offset && rest.ID < offset+limit {
			restaurant := domain.RestaurantResponse{
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

func (r *restaurantRepo) GetById(ctx *domain.CustomContext, id string) (domain.Restaurant, bool) {
	restaurant, isItExists := (*ctx.Restaurants)[id]

	return restaurant, isItExists
}
