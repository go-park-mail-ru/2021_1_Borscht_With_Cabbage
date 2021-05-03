package repository

import (
	"context"
	"database/sql"
	"math"

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

func (r *restaurantRepo) GetVendor(ctx context.Context, params restModel.GetVendorParams) ([]models.RestaurantInfo, error) {
	query :=
		`
	SELECT r.rid, r.name, deliveryCost, avgCheck, description, avatar, ratingssum, reviewscount 
	FROM restaurants as r
	JOIN addresses a on r.rid = a.rid
	WHERE r.rid >= $1 and r.rid <= $2 
	`
	// TODO как сделать ровное количество записей, которые подходят по адресу?

	// если запрос с фильтрацией по адресу
	if params.Address {
		query += ` and ST_DWithin(r.coordinates, ST_SetSRID(ST_Point(a.longitude, a.latitude), 432), a.radius * 1000)`
	}

	restaurantsDB, err := r.DB.Query(query, params.Offset, params.Limit+params.Offset, params.Longitude, params.Latitude)
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return []models.RestaurantInfo{}, failError
	}

	restaurants := make([]models.RestaurantInfo, 0)
	for restaurantsDB.Next() {
		var ratingsSum, reviewsCount int
		restaurant := new(models.RestaurantInfo)
		logger.RepoLevel().InlineInfoLog(ctx, "start scan")

		err = restaurantsDB.Scan(
			&restaurant.ID,
			&restaurant.Title,
			&restaurant.DeliveryCost,
			&restaurant.AvgCheck,
			&restaurant.Description,
			&restaurant.Avatar,
			&ratingsSum,
			&reviewsCount,
		)

		if reviewsCount != 0 {
			restaurant.Rating = math.Round(float64(ratingsSum) / float64(reviewsCount))
		}

		logger.RepoLevel().InlineDebugLog(ctx, restaurant)
		restaurants = append(restaurants, *restaurant)
		logger.RepoLevel().InlineDebugLog(ctx, "stop scan")
	}

	return restaurants, nil
}

func (r *restaurantRepo) GetById(ctx context.Context, id int) (*models.RestaurantWithDishes, error) {
	restaurant := new(models.RestaurantWithDishes)
	var ratingsSum, reviewsCount int
	query :=
		`
	SELECT rid, name, deliveryCost, avgCheck, description, avatar, ratingsSum, reviewsCount
	FROM restaurants 
	WHERE rid=$1
	`

	err := r.DB.QueryRow(query, id).
		Scan(&restaurant.ID, &restaurant.Title, &restaurant.DeliveryCost, &restaurant.AvgCheck,
			&restaurant.Description, &restaurant.Avatar, &ratingsSum, &reviewsCount)
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return nil, failError
	}
	if reviewsCount != 0 {
		restaurant.Rating = math.Round(float64(ratingsSum) / float64(reviewsCount))
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

func (r *restaurantRepo) GetReviews(ctx context.Context, id int) ([]models.RestaurantReview, error) {
	reviewsDB, err := r.DB.Query("select review, stars, deliveryTime, (select name from users where uid=userid) from orders"+
		" where restaurant=(select name from restaurants where rid=$1) and status=$2", id, models.StatusOrderDone)
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return nil, failError
	}

	reviews := make([]models.RestaurantReview, 0)
	for reviewsDB.Next() {
		review := models.RestaurantReview{}
		err = reviewsDB.Scan(
			&review.Review,
			&review.Stars,
			&review.Time,
			&review.UserName,
		)
		if err != nil {
			failError := errors.FailServerError(err.Error())
			logger.RepoLevel().ErrorLog(ctx, failError)
			return nil, failError
		}
		logger.RepoLevel().InlineDebugLog(ctx, review)

		reviews = append(reviews, review)
	}
	err = reviewsDB.Close()
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return nil, failError
	}

	return reviews, nil
}

func (r *restaurantRepo) GetUserAddress(ctx context.Context, uid int) (*models.Address, error) {
	queri := `SELECT name, latitude, longitude FROM addresses WHERE uid = $1`

	logger.RepoLevel().InlineDebugLog(ctx, uid)
	var address models.Address
	err := r.DB.QueryRow(queri, uid).Scan(&address.Name, &address.Latitude, &address.Longitude)
	if err == sql.ErrNoRows {
		logger.RepoLevel().InlineDebugLog(ctx, "end get address not address")
		return &models.Address{}, nil
	}
	if err != nil {
		err := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, err)
		return nil, err
	}

	logger.RepoLevel().InlineDebugLog(ctx, "end get address")
	return &address, nil
}
