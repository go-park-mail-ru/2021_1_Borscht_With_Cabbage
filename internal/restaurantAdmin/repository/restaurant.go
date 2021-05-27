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

func (r *restaurantRepo) GetCategories(ctx context.Context, rid int) ([]string, error) {
	query :=
		`
		SELECT categoryId 
		FROM categories_restaurants
		WHERE restaurantId = $1
	`
	logger.RepoLevel().DebugLog(ctx, logger.Fields{"rid": rid})
	categories := make([]string, 0)
	categoriesDB, err := r.DB.Query(query, rid)
	if err != nil {
		failErr := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failErr)
		return nil, failErr
	}

	for categoriesDB.Next() {
		var category string

		err = categoriesDB.Scan(
			&category,
		)
		if err != nil {
			failErr := errors.FailServerError(err.Error())
			logger.RepoLevel().ErrorLog(ctx, failErr)
			return nil, failErr
		}

		categories = append(categories, category)
	}

	logger.RepoLevel().DebugLog(ctx, logger.Fields{"categories": categories})
	return categories, nil
}

func (r restaurantRepo) GetAddress(ctx context.Context, rid int) (*models.Address, error) {
	query := `SELECT name, latitude, longitude, radius FROM addresses WHERE rid = $1`
	logger.RepoLevel().DebugLog(ctx, logger.Fields{"rid": rid})
	var address models.Address
	err := r.DB.QueryRow(query, rid).Scan(&address.Name, &address.Latitude,
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

func (r restaurantRepo) AddAddress(ctx context.Context, rid int, address models.Address) error {
	query :=
		`
	INSERT INTO addresses (rid, name, latitude, longitude, radius)
	VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.DB.Exec(query, rid, address.Name, address.Latitude, address.Longitude, address.Radius)
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return failError
	}

	return nil
}

func (r restaurantRepo) UpdateAddress(ctx context.Context, rid int, address models.Address) error {
	query :=
		`
	UPDATE addresses SET name = $1, latitude = $2, longitude = $3, radius = $4
	WHERE rid = $5
	`
	_, err := r.DB.Exec(query, address.Name, address.Latitude, address.Longitude, address.Radius, rid)
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return failError
	}

	return nil
}

// TODO: update password
func (r restaurantRepo) UpdateRestaurantData(ctx context.Context, restaurant models.RestaurantUpdateData) error {
	dataToExistingCheck := models.CheckRestaurantExists{
		CurrentRestId: restaurant.ID,
		Email:         restaurant.AdminEmail,
		Number:        restaurant.AdminPhone,
		Name:          restaurant.Title,
	}
	err := r.checkExistingRestaurant(ctx, dataToExistingCheck)
	if err != nil {
		return err
	}

	query :=
		`
	UPDATE restaurants SET name = $1, adminemail = $2, adminphone = $3,
		deliverycost = $4, description = $5
	WHERE rid = $6
	`

	_, err = r.DB.Exec(query,
		restaurant.Title, restaurant.AdminEmail, restaurant.AdminPhone,
		restaurant.DeliveryCost, restaurant.Description, restaurant.ID)
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return failError
	}

	return nil
}

func (r restaurantRepo) checkExistingRestaurant(ctx context.Context, restaurantData models.CheckRestaurantExists) error {
	var userInDB int
	err := r.DB.QueryRow("select rid from restaurants where adminemail = $1", restaurantData.Email).Scan(&userInDB)
	if err != sql.ErrNoRows && userInDB != restaurantData.CurrentRestId {
		return errors.NewErrorWithMessage("Restaurant with this email already exists")
	}

	err = r.DB.QueryRow("select rid from restaurants where adminphone = $1", restaurantData.Number).Scan(&userInDB)
	if err != sql.ErrNoRows && userInDB != restaurantData.CurrentRestId {
		return errors.NewErrorWithMessage("Restaurant with this number already exists")
	}

	err = r.DB.QueryRow("select rid from restaurants where name = $1", restaurantData.Name).Scan(&userInDB)
	if err != sql.ErrNoRows && userInDB != restaurantData.CurrentRestId {
		return errors.NewErrorWithMessage("Restaurant with this name already exists")
	}

	return nil
}

func (r restaurantRepo) UpdateRestaurantImage(ctx context.Context, idRestaurant int, filename string) error {
	_, err := r.DB.Exec("UPDATE restaurants SET avatar = $1 where rid = $2",
		filename, idRestaurant)
	if err != nil {
		dbError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, dbError)
		return dbError
	}

	return nil
}

func (a restaurantRepo) DeleteAllCategories(ctx context.Context, idRestaurant int) error {
	queri :=
		`
	DELETE FROM categories_restaurants
	WHERE restaurantID = $1
	`

	_, err := a.DB.Exec(queri, idRestaurant)
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return failError
	}

	return nil
}

func (a restaurantRepo) AddCategories(ctx context.Context, idRestaurant int, nameCategories []string) error {
	queri :=
		`
	INSERT INTO categories_restaurants (categoryID, restaurantID)
	VALUES ($1, $2)
	`

	// TODO: подумать как это можно сделать одним запросом
	for _, values := range nameCategories {
		a.DB.QueryRow(queri, values, idRestaurant)
	}

	return nil
}
