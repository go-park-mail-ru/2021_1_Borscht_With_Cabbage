package basketServiceRepo

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/utils/calcDistance"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
)

type BasketRepo interface {
	AddToBasket(ctx context.Context, dish models.DishToBasket, uid int) error
	DeleteFromBasket(ctx context.Context, dish models.DishToBasket, uid int) error
	GetBasket(ctx context.Context, uid, rid int) (*models.BasketForUser, error)
	GetBaskets(ctx context.Context, params models.GetBasketParams) (models.BasketsForUser, error)
	AddBasket(ctx context.Context, userID, restaurantID int) (basketID int, err error)
	DeleteBasket(ctx context.Context, basketID int) error
	AddDishToBasket(ctx context.Context, basketID int, dish models.DishInBasket) error
	GetAddress(ctx context.Context, rid int) (*models.Address, error)
}

type basketRepository struct {
	DB *sql.DB
}

func NewBasketRepository(db *sql.DB) BasketRepo {
	return &basketRepository{
		DB: db,
	}
}

func (b basketRepository) AddToBasket(ctx context.Context, dishToBasket models.DishToBasket, uid int) error {
	var basketID int
	// ищем нужную корзину по юзеру
	err := b.DB.QueryRow("select basketID from basket_users where userID = $1 and restaurantID = $2", uid, dishToBasket.RestaurantID).Scan(&basketID)
	if err != nil && err != sql.ErrNoRows {
		logger.RepoLevel().InlineInfoLog(ctx, "Error with finding basket")
		return errors.BadRequestError("Error with finding basket")
	}

	// если к юзеру пока не привязана корзина
	if err == sql.ErrNoRows {
		// то мы ищем ресторан, к которому привязать новую корзину
		err = b.DB.QueryRow("insert into baskets(restaurant) values ($1) returning bid", dishToBasket.RestaurantID).Scan(&basketID)
		fmt.Println(err)
		if err != nil {
			logger.RepoLevel().InlineInfoLog(ctx, "Error with basket creating")
			return errors.BadRequestError("Error with basket creating")
		}

		// создаем новую корзину
		_, err = b.DB.Exec("insert into basket_users (restaurantID, userID, basketid) values ($1, $2, $3)", dishToBasket.RestaurantID, uid, basketID)
		if err != nil {
			insertError := errors.FailServerError(err.Error())
			logger.RepoLevel().ErrorLog(ctx, insertError)
			return insertError
		}
	}

	var dishID int
	err = b.DB.QueryRow("select dish from baskets_food where dish = $1 and basket = $2", dishToBasket.DishID, basketID).Scan(&dishID)
	if err == sql.ErrNoRows {
		_, err = b.DB.Exec("insert into baskets_food (dish, basket, number) values ($1, $2, 1)", dishToBasket.DishID, basketID)
		if err != nil {
			logger.RepoLevel().InlineInfoLog(ctx, "Error while adding dish to basket")
			return errors.BadRequestError("Error while adding dish to basket")
		}
		return nil
	}

	// если есть - увеличиваем количество в корзине
	_, err = b.DB.Exec("update baskets_food set number=number+1 where dish = $1 and basket = $2", dishToBasket.DishID, basketID)
	if err != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Error while inc dish count in basket")
		return errors.BadRequestError("Error while inc dish count in basket")
	}

	return nil
}

func (b basketRepository) DeleteFromBasket(ctx context.Context, dish models.DishToBasket, uid int) error {
	var basketID int
	err := b.DB.QueryRow("select basketID from basket_users where userID = $1 and restaurantid = $2", uid, dish.RestaurantID).Scan(&basketID)
	if err != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Error with basket through user")
		return errors.BadRequestError("Error with basket through user")
	}

	var number int
	err = b.DB.QueryRow("select number from baskets_food where dish = $1 and basket = $2", dish.DishID, basketID).Scan(&number)
	if err != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Error get dish count in basket")
		return errors.BadRequestError("Error get dish count in basket")
	}

	if number == 1 {
		_, err = b.DB.Exec("delete from baskets_food where basket = $1 and dish = $2", basketID, dish.DishID)
		if err != nil {
			logger.RepoLevel().InlineInfoLog(ctx, "Error with deleting dish from basket")
			return errors.BadRequestError("Error with deleting dish from basket")
		}

		return nil
	}

	_, err = b.DB.Exec("update baskets_food set number=number-1 where basket = $1 and dish = $2", basketID, dish.DishID)
	if err != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Error with deleting dish from basket")
		return errors.BadRequestError("Error with deleting dish from basket")
	}

	return nil
}

func (b basketRepository) GetBasket(ctx context.Context, uid, rid int) (*models.BasketForUser, error) {
	var basketRestaurant, imageRestaurant string
	var basketID, restaurantID, deliveryCost int
	err := b.DB.QueryRow("select basketID from basket_users where userID = $1 and restaurantID=$2", uid, rid).Scan(&basketID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Error with getting basket through user and restaurant")
		return nil, errors.BadRequestError("Error with getting basket through user and restaurant")
	}

	basketResponse := models.BasketForUser{
		BID: basketID,
	}

	err = b.DB.QueryRow("select name, deliverycost, avatar from restaurants where rid = $1", rid).Scan(&basketRestaurant, &deliveryCost, &imageRestaurant)
	if err != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Error with finding restaurantID through name")
		return nil, errors.BadRequestError("Error with finding restaurantID through name")
	}
	basketResponse.Restaurant = basketRestaurant
	basketResponse.RID = restaurantID
	basketResponse.DeliveryCost = deliveryCost
	basketResponse.RestaurantImage = imageRestaurant

	dishesDB, errr := b.DB.Query("select d.did, d.name, d.price, bf.number, d.image from baskets_food bf join dishes d on d.did = bf.dish where bf.basket=$1", basketID)
	if errr != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Error with getting basket's dishes")
		return nil, errors.BadRequestError("Error with getting basket's dishes")
	}

	sum := 0
	dishes := make([]models.DishInBasket, 0)
	for dishesDB.Next() {
		dish := new(models.DishInBasket)
		err = dishesDB.Scan(
			&dish.ID,
			&dish.Name,
			&dish.Price,
			&dish.Number,
			&dish.Image)
		sum += dish.Price * dish.Number

		dishes = append(dishes, *dish)
	}
	basketResponse.Foods = dishes
	basketResponse.Summary = sum

	return &basketResponse, nil
}

func (b basketRepository) GetBaskets(ctx context.Context, params models.GetBasketParams) (models.BasketsForUser, error) {
	basketsID, err := b.DB.Query("SELECT basketid, restaurantid FROM basket_users where userid=$1", params.Uid)
	if err != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Error with getting basket's IDs")
		return models.BasketsForUser{}, errors.BadRequestError("Error with getting basket's IDs")
	}

	query := `
	SELECT r.name, deliverycost, avatar, a.radius, a.latitude, a.longitude FROM restaurants r
	JOIN addresses a on r.rid = a.rid
	where r.rid = $1
`
	baskets := models.BasketsForUser{}
	for basketsID.Next() {
		basket := models.BasketForUser{}
		err = basketsID.Scan(&basket.BID, &basket.RID)
		if err != nil {
			logger.RepoLevel().InlineInfoLog(ctx, "Error with getting basket's ID")
			return models.BasketsForUser{}, errors.BadRequestError("Error with getting basket's ID")
		}

		var radius int
		var latitude, longitude string
		err = b.DB.QueryRow(query, basket.RID).Scan(
			&basket.Restaurant,
			&basket.DeliveryCost,
			&basket.RestaurantImage,
			&radius,
			&latitude,
			&longitude,
		)
		if err != nil {
			logger.RepoLevel().InlineInfoLog(ctx, "Error with finding restaurant through id")
			return models.BasketsForUser{}, errors.BadRequestError("Error with finding restaurant through id")
		}

		basket.DeliveryTime = calcDistance.GetDeliveryTime(params.Latitude, params.Longitude, latitude, longitude, radius)

		// getting dishes
		dishesDB, errr := b.DB.Query("select d.did, d.name, d.price, bf.number, d.image from baskets_food bf join dishes d on d.did = bf.dish where bf.basket=$1", basket.BID)
		if errr != nil {
			logger.RepoLevel().InlineInfoLog(ctx, "Error with getting basket's dishes")
			return models.BasketsForUser{}, errors.BadRequestError("Error with getting basket's dishes")
		}

		sum := 0
		dishes := make([]models.DishInBasket, 0)
		for dishesDB.Next() {
			dish := new(models.DishInBasket)
			err = dishesDB.Scan(
				&dish.ID,
				&dish.Name,
				&dish.Price,
				&dish.Number,
				&dish.Image)
			sum += dish.Price * dish.Number
			dishes = append(dishes, *dish)
		}
		basket.Foods = dishes
		basket.Summary = sum

		baskets.Baskets = append(baskets.Baskets, basket)
	}

	return baskets, nil
}

func (b basketRepository) AddBasket(ctx context.Context, userID, restaurantID int) (basketID int, err error) {
	// вносим ее в табличку связи юзер-корзина
	_, err = b.DB.Exec("insert into basket_users (userid, restaurantid) values ($1, $2)", userID, restaurantID)
	if err != nil {
		insertError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, insertError)
		return 0, insertError
	}

	return basketID, nil
}

func (b basketRepository) DeleteBasket(ctx context.Context, basketID int) error {
	_, err := b.DB.Exec("delete from basket_users where basketID = $1", basketID)
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return failError
	}

	return nil
}

func (b basketRepository) AddDishToBasket(ctx context.Context, basketID int, dish models.DishInBasket) error {
	// если такого блюда в корзине нет
	var dishID int
	err := b.DB.QueryRow("select dish from baskets_food where dish = $1 and basket = $2", dish.ID, basketID).Scan(&dishID)
	if err == sql.ErrNoRows {
		_, err = b.DB.Exec("insert into baskets_food (dish, basket, number) values ($1, $2, $3)", dish.ID, basketID, dish.Number)
		if err != nil {
			failError := errors.FailServerError(err.Error() + "error with add dish")
			logger.RepoLevel().ErrorLog(ctx, failError)
			return failError
		}
		return nil
	}

	// если есть - увеличиваем количество в корзине
	_, err = b.DB.Exec("update baskets_food set number=number+$1 where dish = $2 and basket = $3", dish.Number, dish.ID, basketID)
	if err != nil {
		failError := errors.FailServerError(err.Error() + "error with add dish")
		logger.RepoLevel().ErrorLog(ctx, failError)
		return failError
	}
	return nil
}

func (b basketRepository) GetAddress(ctx context.Context, rid int) (*models.Address, error) {
	query := `SELECT name, latitude, longitude, radius FROM addresses WHERE rid = $1`

	logger.RepoLevel().DebugLog(ctx, logger.Fields{"rid": rid})
	var address models.Address
	err := b.DB.QueryRow(query, rid).Scan(&address.Name, &address.Latitude,
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
