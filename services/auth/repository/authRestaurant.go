package authRepo

import (
	"context"
	"database/sql"

	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/services/auth"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
)

type authRestaurantRepo struct {
	DB *sql.DB
}

func NewRestaurantAuthRepo(db *sql.DB) auth.RestaurantAuthRepo {
	return &authRestaurantRepo{
		DB: db,
	}
}

func (a authRestaurantRepo) checkExistingRestaurant(ctx context.Context, restaurantData models.CheckRestaurantExists) error {
	var userInDB int
	err := a.DB.QueryRow("select rid from restaurants where adminemail = $1", restaurantData.Email).Scan(&userInDB)
	if err != sql.ErrNoRows && userInDB != restaurantData.CurrentRestId {
		custErr := errors.NewErrorWithMessage("Restaurant with this email already exists")
		logger.RepoLevel().ErrorLog(ctx, custErr)
		return custErr
	}

	err = a.DB.QueryRow("select rid from restaurants where adminphone = $1", restaurantData.Number).Scan(&userInDB)
	if err != sql.ErrNoRows && userInDB != restaurantData.CurrentRestId {
		custErr := errors.NewErrorWithMessage("Restaurant with this number already exists")
		logger.RepoLevel().ErrorLog(ctx, custErr)
		return custErr
	}

	err = a.DB.QueryRow("select rid from restaurants where name = $1", restaurantData.Name).Scan(&userInDB)
	if err != sql.ErrNoRows && userInDB != restaurantData.CurrentRestId {
		custErr := errors.NewErrorWithMessage("Restaurant with this name already exists")
		logger.RepoLevel().ErrorLog(ctx, custErr)
		return custErr
	}

	return nil
}

func (a authRestaurantRepo) CreateRestaurant(ctx context.Context, newRestaurant models.RestaurantInfo) (int, error) {
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
			avatar, deliveryCost, avgCheck, description)
			values ($1, $2, $3, $4, $5, $6, $7, $8) returning rid`,
		newRestaurant.Title, newRestaurant.AdminPhone, newRestaurant.AdminEmail,
		newRestaurant.AdminHashPassword, newRestaurant.Avatar, newRestaurant.DeliveryCost,
		newRestaurant.AvgCheck, newRestaurant.Description).Scan(&rid)

	if err != nil {
		custError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, custError)
		return 0, custError
	}

	//TODO: временное, пока на фронте нет возможности добавить категории
	_, err = a.DB.Exec(`insert into categories_restaurants (categoryID, restaurantID)
		values('sushi', $1), ('pizza', $1), ('burgers', $1),
			  ('meat', $1), ('fast_food', $1), ('zosh', $1)`, rid)
	if err != nil {
		custError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, custError)
		return 0, custError
	}

	return rid, nil
}

func (a authRestaurantRepo) GetByLogin(ctx context.Context, login string) (*models.RestaurantInfo, error) {
	restaurant := new(models.RestaurantInfo)
	err := a.DB.QueryRow(`select rid, name, adminemail, adminphone, deliveryCost, avgCheck,
		description, avatar, adminpassword from restaurants where (adminphone=$1 or adminemail=$1)`,
		login).
		Scan(&restaurant.ID, &restaurant.Title, &restaurant.AdminEmail, &restaurant.AdminPhone,
			&restaurant.DeliveryCost, &restaurant.AvgCheck, &restaurant.Description,
			&restaurant.Avatar, &restaurant.AdminHashPassword)

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

func (a authRestaurantRepo) GetByRid(ctx context.Context, rid int) (*models.RestaurantInfo, error) {
	DBuser, err := a.DB.Query("select name, adminphone, adminemail, avatar from restaurants where rid=$1", rid)
	if err != nil {
		return nil, errors.NewErrorWithMessage("error with getting restaurant's data").SetDescription("user not found")
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

func (a authRestaurantRepo) GetAddress(ctx context.Context, rid int) (*models.Address, error) {
	queri := `SELECT name, latitude, longitude, radius FROM addresses WHERE rid = $1`

	logger.RepoLevel().DebugLog(ctx, logger.Fields{"rid": rid})
	var address models.Address
	err := a.DB.QueryRow(queri, rid).Scan(&address.Name, &address.Latitude,
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
