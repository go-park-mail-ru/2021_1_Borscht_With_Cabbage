package repository

import (
	"context"
	"database/sql"

	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/restaurantAdmin"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
)

type dishRepo struct {
	DB *sql.DB
}

func NewDishRepo(db *sql.DB) restaurantAdmin.AdminDishRepo {
	return &dishRepo{
		DB: db,
	}
}

func (a dishRepo) GetAllDishes(ctx context.Context, idSection int) ([]models.Dish, error) {
	dishesDB, err := a.DB.Query("select did, section, name, price, weight, description, image from dishes "+
		"where section = $1", idSection)
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return []models.Dish{}, failError
	}

	var dishes []models.Dish
	for dishesDB.Next() {
		dish := new(models.Dish)
		err = dishesDB.Scan(
			&dish.ID,
			&dish.Section,
			&dish.Name,
			&dish.Price,
			&dish.Weight,
			&dish.Description,
			&dish.Image,
		)
		dishes = append(dishes, *dish)
	}

	return dishes, nil
}

func (a dishRepo) UpdateDishData(ctx context.Context, dish models.Dish) error {
	dataToExistingCheck := models.CheckDishExists{
		Name:         dish.Name,
		RestaurantId: dish.Restaurant,
		Id:           dish.ID,
	}
	err := a.checkExistingDish(ctx, dataToExistingCheck)
	if err != nil {
		return err
	}

	_, err = a.DB.Exec(`update dishes set name = $1, price = $2, weight = $3, description = $4
						where did = $5`,
		dish.Name, dish.Price, dish.Weight, dish.Description, dish.ID)
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return failError
	}

	return nil
}

func (a dishRepo) GetDish(ctx context.Context, did int) (*models.Dish, error) {
	DBdish, err := a.DB.Query("select did, restaurantId, name, price, weight, description, image from dishes where did=$1", did)
	if err != nil {
		return nil, errors.AuthorizationError("dish not found")
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
			return nil, failError
		}
	}

	return dish, nil
}

func (a dishRepo) DeleteDish(ctx context.Context, did int) error {
	_, err := a.DB.Exec("delete from dishes where did = $1", did)
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return failError
	}

	return nil
}

func (a dishRepo) checkExistingDish(ctx context.Context, dishData models.CheckDishExists) error {
	dishes, err := a.DB.Query("select did, name from dishes where restaurantId = $1", dishData.RestaurantId)
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
			return errors.NewErrorWithMessage("There is already such a dish")
		}
	}

	return nil
}

func (a dishRepo) AddDish(ctx context.Context, dish models.Dish) (int, error) {
	dataToExistingCheck := models.CheckDishExists{
		Name:         dish.Name,
		RestaurantId: dish.Restaurant,
	}
	err := a.checkExistingDish(ctx, dataToExistingCheck)
	if err != nil {
		return 0, err
	}

	var did int
	err = a.DB.QueryRow(`insert into dishes (restaurantId, section, name, price, 
		weight, description, image) values ($1, $2, $3, $4, $5, $6, $7) returning did`,
		dish.Restaurant, dish.Section, dish.Name, dish.Price,
		dish.Weight, dish.Description, config.ConfigStatic.DefaultDishImage).Scan(&did)
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return 0, failError
	}

	var restaurantName string
	err = a.DB.QueryRow(`select name from restaurants where rid =$1`, dish.Restaurant).Scan(&restaurantName)
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return 0, failError
	}

	_, err = a.DB.Exec(`update dishes set restaurant=$1 where restaurantid=$2`, restaurantName, dish.Restaurant)
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return 0, failError
	}

	return did, nil
}

func (a dishRepo) UpdateDishImage(ctx context.Context, idDish int, filename string) error {
	_, err := a.DB.Exec("UPDATE dishes SET image = $1 where did = $2",
		filename, idDish)
	if err != nil {
		dbError := errors.AuthorizationError("curUser not found")
		logger.RepoLevel().ErrorLog(ctx, dbError)
		return dbError
	}

	return nil
}
