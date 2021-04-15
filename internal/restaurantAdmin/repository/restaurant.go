package repository

import (
	"context"
	"database/sql"
	"net/http"

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

	_, err = a.DB.Exec(
		`update restaurants set name = $1, adminemail = $2, adminphone = $3,
		deliverycost = $4, description = $5
		where rid = $6`,
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
		return errors.NewCustomError(http.StatusBadRequest, "Restaurant with this email already exists")
	}

	err = a.DB.QueryRow("select rid from restaurants where adminphone = $1", restaurantData.Number).Scan(&userInDB)
	if err != sql.ErrNoRows && userInDB != restaurantData.CurrentRestId {
		return errors.NewCustomError(http.StatusBadRequest, "Restaurant with this number already exists")
	}

	err = a.DB.QueryRow("select rid from restaurants where name = $1", restaurantData.Name).Scan(&userInDB)
	if err != sql.ErrNoRows && userInDB != restaurantData.CurrentRestId {
		return errors.NewCustomError(http.StatusBadRequest, "Restaurant with this name already exists")
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

	var rid int
	err = a.DB.QueryRow(`insert into restaurants (name, adminphone, adminemail, adminpassword, 
		avatar, deliveryCost, avgCheck, description, rating) 
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning rid`,
		newRestaurant.Title, newRestaurant.AdminPhone, newRestaurant.AdminEmail,
		newRestaurant.AdminHashPassword, newRestaurant.Avatar, newRestaurant.DeliveryCost,
		newRestaurant.AvgCheck, newRestaurant.Description, newRestaurant.Rating).Scan(&rid)

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
	err := a.DB.QueryRow(`select rid, name, adminemail, adminphone, deliveryCost, avgCheck, 
	description, rating, avatar, adminpassword from restaurants where (adminphone=$1 or adminemail=$1)`,
		login).
		Scan(&restaurant.ID, &restaurant.Title, &restaurant.AdminEmail, &restaurant.AdminPhone,
			&restaurant.DeliveryCost, &restaurant.AvgCheck, &restaurant.Description,
			&restaurant.Rating, &restaurant.Avatar, &restaurant.AdminHashPassword)

	if err == sql.ErrNoRows {
		return nil, errors.NewCustomError(http.StatusBadRequest, "user not found")
	}
	if err != nil {
		custError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, custError)
		return nil, custError
	}

	return restaurant, nil
}

func (a restaurantRepo) GetByRid(ctx context.Context, rid int) (*models.RestaurantInfo, error) {
	DBuser, err := a.DB.Query("select name, adminphone, adminemail, avatar from restaurants where rid=$1", rid)
	if err != nil {
		return nil, errors.AuthorizationError("user not found")
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
