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

func (r *restaurantRepo) GetVendor(ctx context.Context, limit, offset int) ([]models.RestaurantInfo, error) {
	restaurantsDB, err := r.DB.Query("select rid, name, deliveryCost, avgCheck, description, rating, avatar from restaurants "+
		"where rid >= $1 and rid <= $2", offset, limit+offset)
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return []models.RestaurantInfo{}, failError
	}

	var restaurants []models.RestaurantInfo
	for restaurantsDB.Next() {
		restaurant := new(models.RestaurantInfo)
		err = restaurantsDB.Scan(
			&restaurant.ID,
			&restaurant.Title,
			&restaurant.DeliveryCost,
			&restaurant.AvgCheck,
			&restaurant.Description,
			&restaurant.Rating,
			&restaurant.Avatar,
		)

		logger.RepoLevel().InlineDebugLog(ctx, *restaurant)
		restaurants = append(restaurants, *restaurant)
	}

	return restaurants, nil
}

func (r *restaurantRepo) GetById(ctx context.Context, id int) (*models.RestaurantWithDishes, error) {
	restaurant := new(models.RestaurantWithDishes)
	err := r.DB.QueryRow("select rid, name, deliveryCost, avgCheck, description, rating, avatar from restaurants where rid=$1",
		id).Scan(&restaurant.ID, &restaurant.Title, &restaurant.DeliveryCost, &restaurant.AvgCheck, &restaurant.Description, &restaurant.Rating, &restaurant.Avatar)
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return nil, failError
	}

	logger.RepoLevel().InlineDebugLog(ctx, restaurant)

	// Get sections
	sectionsDB, ok := r.DB.Query("select sid, name from sections where restaurant = $1", restaurant.ID)
	if ok != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return nil, failError
	}
	sections := make([]models.Section, 0)
	for sectionsDB.Next() {
		section := new(models.Section)
		err = sectionsDB.Scan(
			&section.ID,
			&section.Name,
		)
		if err != nil {
			failError := errors.FailServerError(err.Error())
			logger.RepoLevel().ErrorLog(ctx, failError)
			return nil, failError
		}
		logger.RepoLevel().InlineDebugLog(ctx, section)

		sections = append(sections, *section)
	}
	err = sectionsDB.Close()
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return nil, failError
	}
	restaurant.Sections = sections

	// get dishes
	dishesDB, errr := r.DB.Query("select did, name, price, weight, description, section, image from dishes where restaurantId = $1", id)
	if errr != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return nil, failError
	}

	dishes := make([]models.Dish, 0)
	for dishesDB.Next() {
		dish := new(models.Dish)
		err = dishesDB.Scan(
			&dish.ID,
			&dish.Name,
			&dish.Price,
			&dish.Weight,
			&dish.Description,
			&dish.Section,
			&dish.Image,
		)
		if err != nil {
			failError := errors.FailServerError(err.Error())
			logger.RepoLevel().ErrorLog(ctx, failError)
			return nil, failError
		}
		logger.RepoLevel().InlineDebugLog(ctx, dish)

		dishes = append(dishes, *dish)
	}
	err = dishesDB.Close()
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return nil, failError
	}

	restaurant.Dishes = dishes
	return restaurant, nil
}
