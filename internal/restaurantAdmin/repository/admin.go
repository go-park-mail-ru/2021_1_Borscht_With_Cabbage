package repository

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/restaurantAdmin"
	_errors "github.com/borscht/backend/utils"
)

type adminRepo struct {
	DB *sql.DB
}

func NewAdminRepo(db *sql.DB) restaurantAdmin.AdminRepo {
	return &adminRepo{
		DB: db,
	}
}

func (a adminRepo) checkExistingRestaurant(ctx context.Context, restaurantData models.CheckRestaurantExists) error {
	var userInDB int
	err := a.DB.QueryRow("select rid from restaurants where adminemail = $1", restaurantData.Email).Scan(&userInDB)
	if err != sql.ErrNoRows && userInDB != restaurantData.CurrentRestId {
		return _errors.NewCustomError(ctx, http.StatusBadRequest, "Restaurant with this email already exists")
	}

	err = a.DB.QueryRow("select rid from restaurants where adminphone = $1", restaurantData.Number).Scan(&userInDB)
	if err != sql.ErrNoRows && userInDB != restaurantData.CurrentRestId {
		return _errors.NewCustomError(ctx, http.StatusBadRequest, "Restaurant with this number already exists")
	}

	err = a.DB.QueryRow("select rid from restaurants where name = $1", restaurantData.Name).Scan(&userInDB)
	if err != sql.ErrNoRows && userInDB != restaurantData.CurrentRestId {
		return _errors.NewCustomError(ctx, http.StatusBadRequest, "Restaurant with this name already exists")
	}

	return nil
}

func (a adminRepo) Create(ctx context.Context, newRestaurant models.Restaurant) (int, error) {
	dataToExistingCheck := models.CheckRestaurantExists{
		Email:  newRestaurant.AdminEmail,
		Number: newRestaurant.AdminPhone,
		Name:   newRestaurant.Name,
	}
	err := a.checkExistingRestaurant(ctx, dataToExistingCheck)
	if err != nil {
		return 0, _errors.FailServer(ctx, err.Error())
	}

	var rid int
	err = a.DB.QueryRow("insert into restaurants (name, adminphone, adminemail, adminpassword, avatar) values ($1, $2, $3, $4, $5) returning rid",
		newRestaurant.Name, newRestaurant.AdminPhone, newRestaurant.AdminEmail, newRestaurant.AdminPassword, config.DefaultAvatar).Scan(&rid)
	if err != nil {
		return 0, _errors.FailServer(ctx, err.Error())
	}

	return rid, nil
}

func (a adminRepo) CheckRestaurantExists(ctx context.Context, restaurantToCheck models.RestaurantAuth) (models.Restaurant, error) {
	restaurant := new(models.Restaurant)
	err := a.DB.QueryRow("select rid, name, avatar from restaurants where (adminphone=$1 or adminemail=$1) and adminpassword=$2",
		restaurantToCheck.Login, restaurantToCheck.Password).Scan(&restaurant.ID, &restaurant.Name, &restaurant.Avatar)

	if err == sql.ErrNoRows {
		return models.Restaurant{}, _errors.NewCustomError(ctx, http.StatusBadRequest, "user not found")
	}
	if err != nil {
		return models.Restaurant{}, _errors.FailServer(ctx, err.Error())
	}

	return *restaurant, nil
}

func (a adminRepo) GetByRid(ctx context.Context, rid int) (models.Restaurant, error) {
	DBuser, err := a.DB.Query("select name, adminphone, adminemail, avatar from restaurants where rid=$1", rid)
	if err != nil {
		return models.Restaurant{}, _errors.Authorization(ctx, "user not found")
	}

	restaurant := new(models.Restaurant)
	for DBuser.Next() {
		err = DBuser.Scan(
			&restaurant.Name,
			&restaurant.AdminPhone,
			&restaurant.AdminEmail,
			&restaurant.Avatar,
		)
		if err != nil {
			return models.Restaurant{}, _errors.FailServer(ctx, err.Error())
		}
	}
	return *restaurant, nil
}
