package repository

import (
	"context"
	"database/sql"
	"fmt"
	"math"

	"github.com/borscht/backend/internal/models"
	restModel "github.com/borscht/backend/internal/restaurant"
	"github.com/borscht/backend/utils/calcDistance"
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

type twoAddresses struct {
	latitude1  float64
	longitude1 float64
	latitude2  float64
	longitude2 float64
}

func (r *restaurantRepo) GetVendor(ctx context.Context, request models.RestaurantRequest) ([]models.RestaurantInfo, error) {
	query :=
		`
		SELECT DISTINCT r.rid, r.name, r.deliveryCost, r.avgCheck, r.description, r.avatar, 
			r.ratingssum, r.reviewscount, r.ordersCount, r.ordersSum, a.latitude, a.longitude, a.radius
		FROM restaurants AS r
		JOIN categories_restaurants AS cr
		ON r.rid = cr.restaurantID
		JOIN categories AS c
		ON cr.categoryID = c.cid
		JOIN addresses a on r.rid = a.rid
		WHERE c.cid = ANY ($1)
		AND r.avgCheck <= $2
	`

	var queryParametres []interface{}
	queryParametres = append(queryParametres, pq.Array(request.Categories), request.Receipt)

	// если запрос с фильтрацией по адресу
	if request.Address {
		logger.RepoLevel().InlineInfoLog(ctx, "vendors request with address")
		query += ` and ST_DWithin(
					    Geography(ST_SetSRID(ST_POINT(a.latitude, a.longitude), 4326)),
						ST_GeomFromText('SRID=4326; POINT(` + fmt.Sprintf("%f", request.LatitudeUser) + ` ` + fmt.Sprintf("%f", request.LongitudeUser) + `)'),
						a.radius)`
	}

	query += `ORDER BY r.rid OFFSET $3 LIMIT $4;`
	queryParametres = append(queryParametres, request.Offset, request.Limit)

	restaurantsDB, err := r.DB.Query(query, queryParametres...)
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return []models.RestaurantInfo{}, failError
	}

	restaurants := make([]models.RestaurantInfo, 0)
	for restaurantsDB.Next() {
		var ratingsSum, reviewsCount int
		restaurant := new(models.RestaurantInfo)

		var restaurantLongitude, restaurantLatitude float64
		var radius int
		var ordersSum, ordersCount int
		err = restaurantsDB.Scan(
			&restaurant.ID,
			&restaurant.Title,
			&restaurant.DeliveryCost,
			&restaurant.AvgCheck,
			&restaurant.Description,
			&restaurant.Avatar,
			&ratingsSum,
			&reviewsCount,
			&ordersCount,
			&ordersSum,
			&restaurantLatitude,
			&restaurantLongitude,
			&radius,
		)

		restaurant.DeliveryTime = calcDistance.GetDeliveryTime(request.LatitudeUser, request.LongitudeUser,
			restaurantLatitude, restaurantLongitude, radius)

		if reviewsCount != 0 {
			restaurant.Rating = math.Round(float64(ratingsSum) / float64(reviewsCount))
		}

		if ordersCount != 0 {
			restaurant.AvgCheck = math.Round(float64(ordersSum) / float64(ordersCount))
		}

		logger.RepoLevel().InlineDebugLog(ctx, restaurant)

		if restaurant.Rating >= request.Rating && restaurant.AvgCheck <= float64(request.Receipt) {
			restaurants = append(restaurants, *restaurant)
		}
		logger.RepoLevel().InlineDebugLog(ctx, "stop scan")
	}

	return restaurants, nil
}

func (r *restaurantRepo) GetById(ctx context.Context, id int, coordinates models.Coordinates) (*models.RestaurantWithDishes, error) {
	restaurant := new(models.RestaurantWithDishes)
	var ratingsSum, reviewsCount int
	query :=
		`
	SELECT r.rid, r.name, deliveryCost, avgCheck, description, avatar, ratingsSum, reviewsCount, a.latitude, a.longitude, a.radius
	FROM restaurants as r
	JOIN addresses a on r.rid = a.rid
	WHERE r.rid=$1
	`

	var restaurantLongitude, restaurantLatitude float64
	var radius int
	err := r.DB.QueryRow(query, id).
		Scan(&restaurant.ID, &restaurant.Title, &restaurant.DeliveryCost, &restaurant.AvgCheck,
			&restaurant.Description, &restaurant.Avatar, &ratingsSum, &reviewsCount, &restaurantLatitude, &restaurantLongitude, &radius)
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return nil, failError
	}
	if reviewsCount != 0 {
		restaurant.Rating = math.Round(float64(ratingsSum) / float64(reviewsCount))
	}

	if coordinates.Latitude != 0 && coordinates.Longitude != 0 {
		restaurant.DeliveryTime = calcDistance.GetDeliveryTime(coordinates.Latitude, coordinates.Longitude, restaurantLatitude, restaurantLongitude, radius)
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
		" where restaurant=(select name from restaurants where rid=$1) and status=$2 and review is not null", id, models.StatusOrderDone)
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
	query := `SELECT name, latitude, longitude FROM addresses WHERE uid = $1`

	logger.RepoLevel().InlineDebugLog(ctx, uid)
	var address models.Address
	err := r.DB.QueryRow(query, uid).Scan(&address.Name, &address.Latitude, &address.Longitude)
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

func (r *restaurantRepo) GetAllCategories(ctx context.Context) ([]string, error) {
	query :=
		`
		SELECT cid FROM categories;
	`
	categoriesDB, err := r.DB.Query(query)
	categories := make([]string, 0)
	for categoriesDB.Next() {
		var category string
		err = categoriesDB.Scan(
			&category,
		)
		if err != nil {
			failError := errors.FailServerError(err.Error())
			logger.RepoLevel().ErrorLog(ctx, failError)
			return nil, failError
		}

		categories = append(categories, category)
	}
	err = categoriesDB.Close()
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return nil, failError
	}
	return categories, nil
}

func (r *restaurantRepo) getRestaurantCategories(ctx context.Context, id int) ([]string, error) {
	query :=
		`
	SELECT categoryID FROM categories_restaurants
	WHERE restaurantID = $1
	`
	categoriesDB, err := r.DB.Query(query, id)
	if err == sql.ErrNoRows {
		logger.RepoLevel().InlineDebugLog(ctx, "getting restaurant categories failed")
		return nil, err
	}

	categories := make([]string, 0)
	for categoriesDB.Next() {
		var categoryID string
		err = categoriesDB.Scan(&categoryID)
		if err != nil {
			logger.RepoLevel().InlineDebugLog(ctx, "getting restaurant categories failed")
			return nil, err
		}
		categories = append(categories, categoryID)
	}

	return categories, nil
}

func (r *restaurantRepo) GetRecommendations(ctx context.Context, params models.RecommendationsParams) ([]models.RestaurantInfo, error) {
	categories, err := r.getRestaurantCategories(ctx, params.Id)
	if err != nil {
		return nil, nil
	}
	logger.RepoLevel().InlineDebugLog(ctx, categories)

	query :=
		`
		SELECT DISTINCT r.rid, r.name, r.deliveryCost, r.avgCheck, r.description, r.avatar, 
			r.ratingssum, r.reviewscount, a.latitude, a.longitude, a.radius
		FROM restaurants AS r
		JOIN categories_restaurants AS cr
		ON r.rid = cr.restaurantID
		JOIN categories AS c
		ON cr.categoryID = c.cid
		JOIN addresses a on r.rid = a.rid
		WHERE c.cid = ANY ($1) AND r.rid != $2
	`

	if params.LongitudeUser != 0 {
		logger.RepoLevel().InlineInfoLog(ctx, "vendors request with address")
		query += ` and ST_DWithin(
					    Geography(ST_SetSRID(ST_POINT(CAST(CAST(a.latitude::float AS varchar) as float), CAST(CAST(a.longitude::float AS varchar) as float)), 4326)),
						ST_GeomFromText('SRID=4326; POINT(` + fmt.Sprintf("%f", params.LatitudeUser) + ` ` + fmt.Sprintf("%f", params.LongitudeUser) + `)'),
						a.radius)`
	}

	query += " LIMIT 6; "

	restaurantsDB, err := r.DB.Query(query, pq.Array(categories), params.Id)
	if err == sql.ErrNoRows {
		logger.RepoLevel().InlineDebugLog(ctx, "getting restaurants for recommendations failed")
		return nil, nil
	}
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return nil, failError
	}

	restaurants := make([]models.RestaurantInfo, 0)
	for restaurantsDB.Next() {
		var ratingsSum, reviewsCount int
		restaurant := new(models.RestaurantInfo)

		var restaurantLongitude, restaurantLatitude float64
		var radius int
		err = restaurantsDB.Scan(
			&restaurant.ID,
			&restaurant.Title,
			&restaurant.DeliveryCost,
			&restaurant.AvgCheck,
			&restaurant.Description,
			&restaurant.Avatar,
			&ratingsSum,
			&reviewsCount,
			&restaurantLatitude,
			&restaurantLongitude,
			&radius,
		)

		restaurant.DeliveryTime = calcDistance.GetDeliveryTime(params.LatitudeUser, params.LongitudeUser,
			restaurantLatitude, restaurantLongitude, radius)

		if reviewsCount != 0 {
			restaurant.Rating = math.Round(float64(ratingsSum) / float64(reviewsCount))
		}

		logger.RepoLevel().InlineDebugLog(ctx, restaurant)
		restaurants = append(restaurants, *restaurant)
	}

	return restaurants, nil
}
