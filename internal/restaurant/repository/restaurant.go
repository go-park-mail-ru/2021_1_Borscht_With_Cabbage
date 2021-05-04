package repository

import (
	"context"
	"database/sql"
	"github.com/borscht/backend/internal/models"
	restModel "github.com/borscht/backend/internal/restaurant"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
	"math"
	"strconv"
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

func deg2rad(deg float64) float64 {
	return deg * (math.Pi / 180)
}

func getDistanceFromLatLonInKm(coordinates twoAddresses) float64 {
	R := 6371.0 // Radius of the Earth in km
	dLat := deg2rad(coordinates.latitude2 - coordinates.latitude1)
	dLon := deg2rad(coordinates.longitude2 - coordinates.longitude1)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(deg2rad(coordinates.latitude1))*math.Cos(deg2rad(coordinates.latitude2))*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := R * c // Distance in km
	return d
}

func getDeliveryTime(latitudeUser, longitudeUser, latitudeRest, longitudeRest string) int {
	latitudeU, latitudeErrU := strconv.ParseFloat(latitudeUser, 64)
	longitudeU, longitudeErrU := strconv.ParseFloat(longitudeUser, 64)
	latitudeR, latitudeErrR := strconv.ParseFloat(latitudeRest, 64)
	longitudeR, longitudeErrR := strconv.ParseFloat(longitudeRest, 64)
	if longitudeErrU == nil && latitudeErrU == nil && latitudeErrR == nil && longitudeErrR == nil {
		distanse := getDistanceFromLatLonInKm(twoAddresses{latitudeU, longitudeU, latitudeR, longitudeR})
		return int(restModel.MinutesInHour*distanse/restModel.CourierSpeed + restModel.CookingTime)
	}
	return 0
}

func (r *restaurantRepo) GetVendor(ctx context.Context, params restModel.GetVendorParams) ([]models.RestaurantInfo, error) {
	query :=
		`
	SELECT r.rid, r.name, deliveryCost, avgCheck, description, avatar, ratingssum, reviewscount, a.latitude, a.longitude
	FROM restaurants as r
	JOIN addresses a on r.rid = a.rid
	WHERE r.rid >= $1 and r.rid <= $2 
	`
	// TODO как сделать ровное количество записей, которые подходят по адресу?

	var queryParametres []interface{}
	queryParametres = append(queryParametres, params.Offset, params.Limit+params.Offset)

	// если запрос с фильтрацией по адресу
	if params.Address {
		logger.RepoLevel().InlineInfoLog(ctx, "vendors request with address")
		query += ` and ST_DWithin(
					   Geography(ST_SetSRID(ST_POINT(a.latitude::float, a.longitude::float), 4326)),
					   ST_GeogFromText($3), a.radius)`
		userAddress := "SRID=4326; POINT(" + params.Latitude + " " + params.Longitude + ")"
		queryParametres = append(queryParametres, userAddress)
	}

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
		logger.RepoLevel().InlineInfoLog(ctx, "start scan")

		var restaurantLongitude, restaurantLatitude string
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
		)
		restaurant.DeliveryTime = getDeliveryTime(params.Latitude, params.Longitude, restaurantLatitude, restaurantLongitude)

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
