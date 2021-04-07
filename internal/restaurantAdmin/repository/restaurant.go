package repository

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/borscht/backend/config"
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

func (a restaurantRepo) Update(ctx context.Context, restaurant models.RestaurantUpdate) error {
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
		adminpassword = $4, deliverycost = $5, description = $6, avatar = $7
		where rid = $8`,
		restaurant.Title, restaurant.AdminEmail, restaurant.AdminPhone, restaurant.AdminPassword,
		restaurant.DeliveryCost, restaurant.Description, restaurant.Avatar, restaurant.ID)
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

func (a restaurantRepo) Create(ctx context.Context, newRestaurant models.Restaurant) (int, error) {
	dataToExistingCheck := models.CheckRestaurantExists{
		Email:  newRestaurant.AdminEmail,
		Number: newRestaurant.AdminPhone,
		Name:   newRestaurant.Title,
	}
	err := a.checkExistingRestaurant(ctx, dataToExistingCheck)
	if err != nil {
		return 0, err
	}

	var rid int
	err = a.DB.QueryRow("insert into restaurants (name, adminphone, adminemail, adminpassword, avatar) values ($1, $2, $3, $4, $5) returning rid",
		newRestaurant.Title, newRestaurant.AdminPhone, newRestaurant.AdminEmail, newRestaurant.AdminPassword, config.DefaultAvatar).Scan(&rid)
	if err != nil {
		custError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, custError)
		return 0, custError
	}

	return rid, nil
}

func (a restaurantRepo) CheckRestaurantExists(ctx context.Context, restaurantToCheck models.RestaurantAuth) (*models.Restaurant, error) {
	restaurant := new(models.Restaurant)
	err := a.DB.QueryRow("select rid, name, avatar from restaurants where (adminphone=$1 or adminemail=$1) and adminpassword=$2",
		restaurantToCheck.Login, restaurantToCheck.Password).Scan(&restaurant.ID, &restaurant.Title, &restaurant.Avatar)

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

func (a restaurantRepo) GetByRid(ctx context.Context, rid int) (models.Restaurant, error) {
	DBuser, err := a.DB.Query("select name, adminphone, adminemail, avatar from restaurants where rid=$1", rid)
	if err != nil {
		return models.Restaurant{}, errors.AuthorizationError("user not found")
	}

	restaurant := new(models.Restaurant)
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
			return models.Restaurant{}, custError
		}
	}
	return *restaurant, nil
}
