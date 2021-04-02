package repository

import (
	"context"
	"database/sql"

	"github.com/borscht/backend/internal/models"
	restModel "github.com/borscht/backend/internal/restaurant"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
)

type restaurantRepo struct {
	DB *sql.DB
}

func NewRestaurantRepo(db *sql.DB) restModel.RestaurantRepo {
	return &restaurantRepo{
		DB: db,
	}
}

func (r *restaurantRepo) GetVendor(ctx context.Context, limit, offset int) ([]models.RestaurantResponse, error) {
	restaurantsDB, err := r.DB.Query("select rid, name, deliveryCost, avgCheck, description, rating, avatar from restaurants "+
		"where rid >= $1 and rid <= $2", offset, limit+offset)
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return []models.RestaurantResponse{}, failError
	}

	var restaurants []models.RestaurantResponse
	for restaurantsDB.Next() {
		restaurant := new(models.RestaurantResponse)
		err = restaurantsDB.Scan(
			&restaurant.ID,
			&restaurant.Title,
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

func (r *restaurantRepo) GetById(ctx context.Context, id string) (models.Restaurant, error) {
	restaurant := new(models.Restaurant)
	err := r.DB.QueryRow("select name, deliveryCost, avgCheck, description, rating, avatar from restaurants where rid=$1",
		id).Scan(&restaurant.Title, &restaurant.DeliveryCost, &restaurant.AvgCheck, &restaurant.Description, &restaurant.Rating, &restaurant.Avatar)
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return models.Restaurant{}, failError
	}

	dishesDB, errr := r.DB.Query("select name, price, weight, description, image from dishes where did = $1", id)
	if errr != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return models.Restaurant{}, failError
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
			failError := errors.FailServerError(err.Error())
			logger.RepoLevel().ErrorLog(ctx, failError)
			return models.Restaurant{}, failError
		}

		dishes = append(dishes, *dish)
	}
	err = dishesDB.Close()
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return models.Restaurant{}, failError
	}

	restaurant.Dishes = dishes
	return *restaurant, nil
}
