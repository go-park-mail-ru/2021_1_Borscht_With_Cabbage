package repository

import (
	"context"
	"database/sql"

	"github.com/borscht/backend/internal/models"
	restModel "github.com/borscht/backend/internal/restaurant"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
	"github.com/lib/pq"
)

type restaurantRepo struct {
	DB *sql.DB
}

func NewRestaurantRepo(db *sql.DB) restModel.RestaurantRepo {
	return &restaurantRepo{
		DB: db,
	}
}

func (r *restaurantRepo) GetVendorWithCategory(ctx context.Context, request models.RestaurantRequest) ([]models.RestaurantInfo, error) {
	query := `
	SELECT DISTINCT r.rid, r.name, r.deliveryCost, r.avgCheck, r.description, r.rating, r.avatar
	FROM restaurants AS r
	JOIN categories_restaurants AS cr
	ON r.rid = cr.restaurantID
	JOIN categories AS c
	ON cr.categoryID = c.cid
	WHERE c.cid = ANY ($1)
	AND r.rating >= $2 AND r.avgCheck <= $3
	ORDER BY r.rating DESC, r.rid OFFSET $4 LIMIT $5;
  	`

	restaurantsDB, err := r.DB.Query(query, pq.Array(request.Categories),
		request.Rating, request.Receipt, request.Offset, request.Limit)

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
		if err != nil {
			failError := errors.FailServerError(err.Error())
			logger.RepoLevel().ErrorLog(ctx, failError)
			return []models.RestaurantInfo{}, failError
		}

		logger.RepoLevel().InlineDebugLog(ctx, *restaurant)
		restaurants = append(restaurants, *restaurant)
	}

	return restaurants, nil
}

func (r *restaurantRepo) GetVendor(ctx context.Context, request models.RestaurantRequest) ([]models.RestaurantInfo, error) {
	query := `
	SELECT DISTINCT rid, name, deliveryCost, avgCheck, description, rating, avatar FROM restaurants
	WHERE rating >= $1 AND avgCheck <= $2
	ORDER BY rating DESC, rid OFFSET $3 LIMIT $4;
  	`

	restaurantsDB, err := r.DB.Query(query, request.Rating, request.Receipt,
		request.Offset, request.Limit)

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
		if err != nil {
			failError := errors.FailServerError(err.Error())
			logger.RepoLevel().ErrorLog(ctx, failError)
			return []models.RestaurantInfo{}, failError
		}

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
