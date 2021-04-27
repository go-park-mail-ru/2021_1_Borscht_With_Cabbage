package repository

import (
	"context"
	"database/sql"

	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/restaurantAdmin"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
)

type restaurantRepo struct {
	DB *sql.DB
}

func NewRestaurantRepo(db *sql.DB) restaurantAdmin.AdminRestaurantRepo {
	return &restaurantRepo{
		DB: db,
	}
}

func (u restaurantRepo) GetAddress(ctx context.Context, rid int) (*models.Address, error) {
	queri := `SELECT name, latitude, longitude, radius FROM addresses WHERE rid = $1`

	logger.RepoLevel().DebugLog(ctx, logger.Fields{"rid": rid})
	var address models.Address
	err := u.DB.QueryRow(queri, rid).Scan(&address.Name, &address.Latitude,
		&address.Longitude, &address.Radius)

	if err == sql.ErrNoRows {
		return &models.Address{}, nil
	}
	if err != nil {
		err := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, err)
		return nil, err
	}

	logger.RepoLevel().DebugLog(ctx, logger.Fields{"address": address})
	return &address, nil
}

func (u restaurantRepo) AddAddress(ctx context.Context, rid int, address models.Address) error {
	query :=
		`
	INSERT INTO addresses (rid, name, latitude, longitude, radius)
	VALUES ($1, $2, $3, $4, $5)
	`
	_, err := u.DB.Exec(query, rid, address.Name, address.Latitude, address.Longitude, address.Radius)
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return failError
	}

	return nil
}

func (u restaurantRepo) UpdateAddress(ctx context.Context, rid int, address models.Address) error {
	query :=
		`
	UPDATE addresses SET name = $1, latitude = $2, longitude = $3, radius = $4
	WHERE rid = $5
	`
	_, err := u.DB.Exec(query, address.Name, address.Latitude, address.Longitude, address.Radius, rid)
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return failError
	}

	return nil
}

// TODO: update password
func (a restaurantRepo) UpdateRestaurantData(ctx context.Context, restaurant models.RestaurantUpdateData) error {
	dataToExistingCheck := models.CheckRestaurantExists{
		CurrentRestId: restaurant.ID,
		Email:         restaurant.AdminEmail,
		Number:        restaurant.AdminPhone,
		Name:          restaurant.Title,
	}
	err := a.checkExistingRestaurant(ctx, dataToExistingCheck)
	if err != nil {
		return err
	}

	query :=
		`
	UPDATE restaurants SET name = $1, adminemail = $2, adminphone = $3,
		deliverycost = $4, description = $5
	WHERE rid = $6
	`

	_, err = a.DB.Exec(query,
		restaurant.Title, restaurant.AdminEmail, restaurant.AdminPhone,
		restaurant.DeliveryCost, restaurant.Description, restaurant.ID)
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return failError
	}

	return nil
}

func (a restaurantRepo) checkExistingRestaurant(ctx context.Context, restaurantData models.CheckRestaurantExists) error {
	var userInDB int
	err := a.DB.QueryRow("select rid from restaurants where adminemail = $1", restaurantData.Email).Scan(&userInDB)
	if err != sql.ErrNoRows && userInDB != restaurantData.CurrentRestId {
		return errors.NewErrorWithMessage("Restaurant with this email already exists")
	}

	err = a.DB.QueryRow("select rid from restaurants where adminphone = $1", restaurantData.Number).Scan(&userInDB)
	if err != sql.ErrNoRows && userInDB != restaurantData.CurrentRestId {
		return errors.NewErrorWithMessage("Restaurant with this number already exists")
	}

	err = a.DB.QueryRow("select rid from restaurants where name = $1", restaurantData.Name).Scan(&userInDB)
	if err != sql.ErrNoRows && userInDB != restaurantData.CurrentRestId {
		return errors.NewErrorWithMessage("Restaurant with this name already exists")
	}

	return nil
}

func (a restaurantRepo) CreateRestaurant(ctx context.Context, newRestaurant models.RestaurantInfo) (int, error) {
	logger.RepoLevel().InfoLog(ctx, logger.Fields{"create restaurant": newRestaurant})
	dataToExistingCheck := models.CheckRestaurantExists{
		Email:  newRestaurant.AdminEmail,
		Number: newRestaurant.AdminPhone,
		Name:   newRestaurant.Title,
	}
	err := a.checkExistingRestaurant(ctx, dataToExistingCheck)
	if err != nil {
		return 0, err
	}
	logger.RepoLevel().InlineDebugLog(ctx, "correct restaurant")

	query :=
		`
	INSERT INTO restaurants (name, adminphone, adminemail, adminpassword, 
		avatar, deliveryCost, avgCheck, description, rating) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
	RETURNING rid
	`
	var rid int
	err = a.DB.QueryRow(query,
		newRestaurant.Title, newRestaurant.AdminPhone, newRestaurant.AdminEmail,
		newRestaurant.AdminHashPassword, newRestaurant.Avatar, newRestaurant.DeliveryCost,
		newRestaurant.AvgCheck, newRestaurant.Description, newRestaurant.Rating).
		Scan(&rid)

	if err != nil {
		custError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, custError)
		return 0, custError
	}

	return rid, nil
}

func (a restaurantRepo) UpdateRestaurantImage(ctx context.Context, idRestaurant int, filename string) error {
	_, err := a.DB.Exec("UPDATE restaurants SET avatar = $1 where rid = $2",
		filename, idRestaurant)
	if err != nil {
		dbError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, dbError)
		return dbError
	}

	return nil
}

func (a restaurantRepo) GetByLogin(ctx context.Context, login string) (*models.RestaurantInfo, error) {
	restaurant := new(models.RestaurantInfo)
	query :=
		`
	SELECT rid, name, adminemail, adminphone, deliveryCost, avgCheck, 
		description, rating, avatar, adminpassword 
	FROM restaurants 
	WHERE (adminphone=$1 or adminemail=$1)
	`
	err := a.DB.QueryRow(query, login).
		Scan(&restaurant.ID, &restaurant.Title, &restaurant.AdminEmail, &restaurant.AdminPhone,
			&restaurant.DeliveryCost, &restaurant.AvgCheck, &restaurant.Description,
			&restaurant.Rating, &restaurant.Avatar, &restaurant.AdminHashPassword)

	if err == sql.ErrNoRows {
		return nil, errors.NewErrorWithMessage("not authorization").SetDescription("user not found")
	}
	if err != nil {
		custError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, custError)
		return nil, custError
	}

	return restaurant, nil
}

func (a restaurantRepo) GetByRid(ctx context.Context, rid int) (*models.RestaurantInfo, error) {
	query :=
		`
	SELECT name, adminphone, adminemail, avatar
	FROM restaurants 
	WHERE rid=$1
	`
	DBuser, err := a.DB.Query(query, rid)
	if err != nil {
		return nil, errors.NewErrorWithMessage("not authorization").SetDescription("user not found")
	}

	restaurant := new(models.RestaurantInfo)
	for DBuser.Next() {
		err = DBuser.Scan(
			&restaurant.Title,
			&restaurant.AdminPhone,
			&restaurant.AdminEmail,
			&restaurant.Avatar,
		)
		if err != nil {
			custError := errors.FailServerError(err.Error())
			logger.RepoLevel().ErrorLog(ctx, custError)
			return nil, custError
		}
	}
	return restaurant, nil
}
