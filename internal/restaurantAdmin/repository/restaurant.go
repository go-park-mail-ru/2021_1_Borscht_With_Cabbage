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
