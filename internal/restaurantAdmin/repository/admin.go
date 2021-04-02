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

type adminRepo struct {
	DB *sql.DB
}

func NewAdminRepo(db *sql.DB) restaurantAdmin.AdminRepo {
	return &adminRepo{
		DB: db,
	}
}

func (a adminRepo) Update(ctx context.Context, restaurant models.RestaurantUpdate) error {
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

func (a adminRepo) UpdateDish(ctx context.Context, dish models.Dish) error {
	restaurant, ok := ctx.Value("Restaurant").(models.Restaurant)
	if !ok {
		failError := errors.FailServerError("failed to convert to models.Restaurant")
		logger.RepoLevel().ErrorLog(ctx, failError)
		return failError
	}
	dataToExistingCheck := models.CheckDishExists{
		Name:         dish.Name,
		RestaurantId: restaurant.ID,
		Id:           dish.ID,
	}
	err := a.checkExistingDish(ctx, dataToExistingCheck)
	if err != nil {
		return err
	}

	_, err = a.DB.Exec(`update dishes set name = $1, price = $2, weight = $3, 
						description = $4, image = $5
						where did = $6`,
		dish.Name, dish.Price, dish.Weight, dish.Description, dish.Image, dish.ID)
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return failError
	}

	return nil
}

func (a adminRepo) GetDish(ctx context.Context, did int) (models.Dish, error) {
	DBdish, err := a.DB.Query("select did, restaurant, name, price, weight, description, image from dishes where did=$1", did)
	if err != nil {
		return models.Dish{}, errors.AuthorizationError("dish not found")
	}

	dish := new(models.Dish)
	for DBdish.Next() {
		err = DBdish.Scan(
			&dish.ID,
			&dish.Restaurant,
			&dish.Name,
			&dish.Price,
			&dish.Weight,
			&dish.Description,
			&dish.Image,
		)
		if err != nil {
			failError := errors.FailServerError(err.Error())
			logger.RepoLevel().ErrorLog(ctx, failError)
			return models.Dish{}, failError
		}
	}

	return *dish, nil
}

func (a adminRepo) DeleteDish(ctx context.Context, did int) error {
	_, err := a.DB.Exec("delete from dishes where did = $1", did)
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return failError
	}

	return nil
}

func (a adminRepo) checkExistingDish(ctx context.Context, dishData models.CheckDishExists) error {
	dishes, err := a.DB.Query("select did, name from dishes where restaurant = $1", dishData.RestaurantId)
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return failError
	}
	for dishes.Next() {
		nameDish := new(string)
		didDish := new(int)
		dishes.Scan(&didDish, &nameDish)
		if *nameDish == dishData.Name && *didDish != dishData.Id {
			return errors.NewCustomError(http.StatusBadRequest, "There is already such a dish")
		}
	}

	return nil
}

func (a adminRepo) AddDish(ctx context.Context, dish models.Dish) (int, error) {
	restaurant, ok := ctx.Value("Restaurant").(models.Restaurant)
	if !ok {
		failError := errors.FailServerError("failed to convert to models.Restaurant")
		logger.RepoLevel().ErrorLog(ctx, failError)
		return 0, failError
	}
	dataToExistingCheck := models.CheckDishExists{
		Name:         dish.Name,
		RestaurantId: restaurant.ID,
	}
	err := a.checkExistingDish(ctx, dataToExistingCheck)
	if err != nil {
		return 0, err
	}

	var did int
	err = a.DB.QueryRow("insert into dishes (restaurant, name, price, weight, description, image) values ($1, $2, $3, $4, $5, $6) returning did",
		restaurant.ID, dish.Name, dish.Price, dish.Weight, dish.Description, config.DefaultAvatar).Scan(&did)
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return 0, failError
	}

	return did, nil
}

func (a adminRepo) checkExistingRestaurant(ctx context.Context, restaurantData models.CheckRestaurantExists) error {
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

func (a adminRepo) Create(ctx context.Context, newRestaurant models.Restaurant) (int, error) {
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

func (a adminRepo) CheckRestaurantExists(ctx context.Context, restaurantToCheck models.RestaurantAuth) (models.Restaurant, error) {
	restaurant := new(models.Restaurant)
	err := a.DB.QueryRow("select rid, name, avatar from restaurants where (adminphone=$1 or adminemail=$1) and adminpassword=$2",
		restaurantToCheck.Login, restaurantToCheck.Password).Scan(&restaurant.ID, &restaurant.Title, &restaurant.Avatar)

	if err == sql.ErrNoRows {
		return models.Restaurant{}, errors.NewCustomError(http.StatusBadRequest, "user not found")
	}
	if err != nil {
		custError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, custError)
		return models.Restaurant{}, custError
	}

	return *restaurant, nil
}

func (a adminRepo) GetByRid(ctx context.Context, rid int) (models.Restaurant, error) {
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
