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
	restaurantsDB, err := r.DB.Query("select rid, name, deliveryCost, avgCheck, description, rating, avatar from restaurants "+
		"where rid >= $1 and rid <= $2", offset, limit+offset)
	if err != nil {
		return []models.RestaurantResponse{}, _errors.FailServer(err.Error())
	}

	var restaurants []models.RestaurantResponse
	for restaurantsDB.Next() {
		restaurant := new(models.RestaurantResponse)
		err = restaurantsDB.Scan(
			&restaurant.ID,
			&restaurant.Name,
			&restaurant.DeliveryCost,
			&restaurant.AvgCheck,
			&restaurant.Description,
			&restaurant.Rating,
			&restaurant.Avatar,
		)
		restaurants = append(restaurants, *restaurant)
	}

	return restaurants, nil
}

func (r *restaurantRepo) GetById(id string) (models.Restaurant, error) {
	restaurant := new(models.Restaurant)
	err := r.DB.QueryRow("select name, deliveryCost, avgCheck, description, rating, avatar from restaurants where rid=$1",
		id).Scan(&restaurant.Name, &restaurant.DeliveryCost, &restaurant.AvgCheck, &restaurant.Description, &restaurant.Rating, &restaurant.Avatar)
	if err != nil {
		return models.Restaurant{}, _errors.FailServer(err.Error())
	}

	dishesDB, errr := r.DB.Query("select name, price, weight, description, image from dishes where did = $1", id)
	if errr != nil {
		return models.Restaurant{}, _errors.FailServer(errr.Error())
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
			return models.Restaurant{}, _errors.FailServer(err.Error())
		}

		dishes = append(dishes, *dish)
	}
	err = dishesDB.Close()
	if err != nil {
		return models.Restaurant{}, _errors.FailServer(err.Error())
	}

	restaurant.Dishes = dishes
	return *restaurant, nil
}
