package basketServiceRepo

import (
	"context"
	"database/sql"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
)

type BasketRepo interface {
	AddToBasket(ctx context.Context, dish models.DishToBasket, uid int) error
	DeleteFromBasket(ctx context.Context, dish models.DishToBasket, uid int) error
	GetBasket(ctx context.Context, uid int) (*models.BasketForUser, error)
	AddBasket(ctx context.Context, userID, restaurantID int) (basketID int, err error)
	DeleteBasket(ctx context.Context, userID, basketID int) error
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
	var basketRestaurant string
	var basketID int
	// ищем нужную корзину по юзеру
	err := b.DB.QueryRow("select basketID from basket_users where userID = $1", uid).Scan(&basketID)
	// если к юзеру пока не привязана корзина
	if err == sql.ErrNoRows {
		// то мы ищем ресторан, к которому привязать новую корзину
		err = b.DB.QueryRow("select restaurant from dishes where did = $1", dishToBasket.DishID).Scan(&basketRestaurant)
		if err != nil {
			logger.RepoLevel().InlineInfoLog(ctx, "Error with finding restaurant through dish")
			return errors.BadRequestError("Error with finding restaurant through dish")
		}

		// создаем новую корзину
		err = b.DB.QueryRow("insert into baskets (restaurant) values ($1) returning bid", basketRestaurant).Scan(&basketID)
		if err != nil {
			insertError := errors.FailServerError(err.Error())
			logger.RepoLevel().ErrorLog(ctx, insertError)
			return insertError
		}

		// вносим ее в табличку связи юзер-корзина
		_, err = b.DB.Exec("insert into basket_users (basketid, userid) values ($1, $2)", basketID, uid)
		if err != nil {
			insertError := errors.FailServerError(err.Error())
			logger.RepoLevel().ErrorLog(ctx, insertError)
			return insertError
		}
	}

	// если добавляем в корзину с тем же рестораном, то просто добавляем блюдо
	if dishToBasket.SameBasket {
		// если такого блюда в корзине нет
		var dishID int
		err := b.DB.QueryRow("select dish from baskets_food where dish = $1 and basket = $2", dishToBasket.DishID, basketID).Scan(&dishID)

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

	// если ресторан теперь другой, то чистим старую корзину и добавляем блюдо
	_, err = b.DB.Exec("delete from baskets_food where basket = $1", basketID)
	if err != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Error with deleting previous dishes from basket")
		return errors.BadRequestError("Error with deleting previous dishes from basket")
	}

	_, err = b.DB.Exec("insert into baskets_food (basket, dish, number) values ($1, $2, 1)", basketID, dishToBasket.DishID)
	if err != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Error while adding dish to basket")
		return errors.BadRequestError("Error while adding dish to basket")
	}

	_, err = b.DB.Exec("update baskets set restaurant = (select restaurant from dishes where did=$1) where bid=$2", dishToBasket.DishID, basketID)
	if err != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Error while adding dish to basket")
		return errors.BadRequestError("Error while adding dish to basket")
	}

	return nil
}

func (b basketRepository) DeleteFromBasket(ctx context.Context, dish models.DishToBasket, uid int) error {
	var basketID int
	err := b.DB.QueryRow("select basketID from basket_users where userID = $1", uid).Scan(&basketID)
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

func (b basketRepository) GetBasket(ctx context.Context, uid int) (*models.BasketForUser, error) {
	var basketRestaurant, imageRestaurant string
	var basketID, restaurantID, deliveryCost int
	err := b.DB.QueryRow("select basketID from basket_users where userID = $1", uid).Scan(&basketID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Error with getting basket through user")
		return nil, errors.BadRequestError("Error with getting basket through user")
	}

	basketResponse := models.BasketForUser{
		BID: basketID,
	}

	err = b.DB.QueryRow("select restaurant from baskets where bid = $1", basketID).Scan(&basketRestaurant)
	if err != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Error with finding restaurant through basket")
		return nil, errors.BadRequestError("Error with finding restaurant through basket")
	}
	basketResponse.Restaurant = basketRestaurant

	err = b.DB.QueryRow("select avatar from restaurants where name = $1", basketRestaurant).Scan(&imageRestaurant)
	if err != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Error with getting restaurant image")
		return nil, errors.BadRequestError("Error with getting restaurant image")
	}
	basketResponse.RestaurantImage = imageRestaurant

	err = b.DB.QueryRow("select rid, deliverycost from restaurants where name = $1", basketRestaurant).Scan(&restaurantID, &deliveryCost)
	if err != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Error with finding restaurantID through name")
		return nil, errors.BadRequestError("Error with finding restaurantID through name")
	}
	basketResponse.RID = restaurantID
	basketResponse.DeliveryCost = deliveryCost

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

func (b basketRepository) AddBasket(ctx context.Context, userID, restaurantID int) (basketID int, err error) {
	// TODO: убрать когда будут внешние ключи id
	var restaurantName string
	err = b.DB.QueryRow("select name from restaurants where rid = $1", restaurantID).Scan(&restaurantName)

	err = b.DB.QueryRow("insert into baskets (restaurant) values ($1) returning bid", restaurantName).Scan(&basketID)
	if err != nil {
		insertError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, insertError)
		return 0, insertError
	}

	// вносим ее в табличку связи юзер-корзина
	_, err = b.DB.Exec("insert into basket_users (basketid, userid) values ($1, $2)", basketID, userID)
	if err != nil {
		insertError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, insertError)
		return 0, insertError
	}

	return basketID, nil
}

func (b basketRepository) DeleteBasket(ctx context.Context, userID, basketID int) error {
	_, err := b.DB.Exec("delete from basket_users where basketID = $1 and userID = $2", basketID, userID)
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return failError
	}

	_, err = b.DB.Exec("delete from baskets where bid = $1", basketID)
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
	queri := `SELECT name, latitude, longitude, radius FROM addresses WHERE rid = $1`

	logger.RepoLevel().DebugLog(ctx, logger.Fields{"rid": rid})
	var address models.Address
	err := b.DB.QueryRow(queri, rid).Scan(&address.Name, &address.Latitude,
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
