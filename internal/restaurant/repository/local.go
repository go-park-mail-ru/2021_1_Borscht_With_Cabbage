package repository

import (
	"database/sql"
	"github.com/borscht/backend/internal/models"
	restModel "github.com/borscht/backend/internal/restaurant"
	_errors "github.com/borscht/backend/utils"
)

type restaurantRepo struct {
	DB *sql.DB
}

func NewRestaurantRepo(db *sql.DB) restModel.RestaurantRepo {
	return &restaurantRepo{
		DB: db,
	}
}

func (r *restaurantRepo) GetVendor(limit, offset int) ([]models.RestaurantResponse, error) {
	//var result []models.RestaurantResponse
	//
	//for _, rest := range *r.db.GetModels().Restaurants {
	//	if rest.ID >= offset && rest.ID < offset+limit {
	//		restaurant := models.RestaurantResponse{
	//			ID:           rest.ID,
	//			Name:         rest.Name,
	//			Rating:       rest.Rating,
	//			DeliveryCost: rest.DeliveryCost,
	//			AvgCheck:     rest.AvgCheck,
	//			Description:  rest.Description,
	//		}
	//		result = append(result, restaurant)
	//	}
	//}
	//
	//return result, nil
	return []models.RestaurantResponse{}, nil
}

func (r *restaurantRepo) GetById(id string) (models.Restaurant, bool, error) {
	restaurantDB, err := r.DB.Query("select name, deliveryCost, avgCheck, description, rating, avatar from restaurants where rid=$1", id)
	if err != nil {
		return models.Restaurant{}, false, _errors.FailServer(err.Error())
	}

	restaurant := new(models.Restaurant)
	for restaurantDB.Next() {
		err = restaurantDB.Scan(
			&restaurant.Name,
			&restaurant.DeliveryCost,
			&restaurant.AvgCheck,
			&restaurant.Description,
			&restaurant.Rating,
			&restaurant.Avatar,
		)
		if err != nil {
			return models.Restaurant{}, false, _errors.FailServer(err.Error())
		}
	}
	err = restaurantDB.Close()
	if err != nil {
		return models.Restaurant{}, false, _errors.FailServer(err.Error())

	}

	dishesDB, errr := r.DB.Query("select name, price, weight, description, image from dishes where did = $1", id)
	if errr != nil {
		return models.Restaurant{}, false, _errors.FailServer(errr.Error())
	}

	dishes := make([]models.Dish, 0)
	for dishesDB.Next() {
		dish := new(models.Dish)
		err = dishesDB.Scan(
			&dish.Name,
			&dish.Price,
			&dish.Weight,
			&dish.Description,
			&dish.Image,
		)
		if err != nil {
			return models.Restaurant{}, false, _errors.FailServer(err.Error())
		}

		dishes = append(dishes, *dish)
	}
	err = dishesDB.Close()
	if err != nil {
		return models.Restaurant{}, false, _errors.FailServer(err.Error())
	}

	restaurant.Dishes = dishes
	return *restaurant, true, nil
}
