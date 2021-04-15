package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/order"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
	"time"
)

type orderRepo struct {
	DB *sql.DB
}

func NewOrderRepo(db *sql.DB) order.OrderRepo {
	return &orderRepo{
		DB: db,
	}
}

func (o orderRepo) AddToBasket(ctx context.Context, dishToBasket models.DishToBasket, uid int) error {
	var basketRestaurant string
	var basketID int
	// ищем нужную корзину по юзеру
	err := o.DB.QueryRow("select basketID from basket_users where userID = $1", uid).Scan(&basketID)
	// если к юзеру пока не привязана корзина
	if err == sql.ErrNoRows {
		// то мы ищем ресторан, к которому привязать новую корзину
		err = o.DB.QueryRow("select restaurant from dishes where did = $1", dishToBasket.DishID).Scan(&basketRestaurant)
		if err != nil {
			logger.RepoLevel().InlineInfoLog(ctx, "Error with finding restaurant through dish")
			return errors.BadRequestError("Error with finding restaurant through dish")
		}

		// создаем новую корзину
		err = o.DB.QueryRow("insert into baskets (restaurant) values ($1) returning bid", basketRestaurant).Scan(&basketID)
		if err != nil {
			insertError := errors.FailServerError(err.Error())
			logger.RepoLevel().ErrorLog(ctx, insertError)
			return insertError
		}

		// вносим ее в табличку связи юзер-корзина
		_, err = o.DB.Exec("insert into basket_users (basketid, userid) values ($1, $2)", basketID, uid)
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
		err := o.DB.QueryRow("select dish from baskets_food where dish = $1 and basket = $2", dishToBasket.DishID, basketID).Scan(&dishID)

		if err == sql.ErrNoRows {
			_, err = o.DB.Exec("insert into baskets_food (dish, basket, number) values ($1, $2, 1)", dishToBasket.DishID, basketID)
			if err != nil {
				logger.RepoLevel().InlineInfoLog(ctx, "Error while adding dish to basket")
				return errors.BadRequestError("Error while adding dish to basket")
			}
			return nil
		}

		// если есть - увеличиваем количество в корзине
		_, err = o.DB.Exec("update baskets_food set number=number+1 where dish = $1 and basket = $2", dishToBasket.DishID, basketID)
		if err != nil {
			logger.RepoLevel().InlineInfoLog(ctx, "Error while inc dish count in basket")
			return errors.BadRequestError("Error while inc dish count in basket")
		}
		return nil
	}

	// если ресторан теперь другой, то чистим старую корзину и добавляем блюдо
	_, err = o.DB.Exec("delete from baskets_food where basket = $1", basketID)
	if err != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Error with deleting previous dishes from basket")
		return errors.BadRequestError("Error with deleting previous dishes from basket")
	}

	_, err = o.DB.Exec("insert into baskets_food (basket, dish, number) values ($1, $2, 1)", basketID, dishToBasket.DishID)
	if err != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Error while adding dish to basket")
		return errors.BadRequestError("Error while adding dish to basket")
	}

	_, err = o.DB.Exec("update baskets set restaurant = (select restaurant from dishes where did=$1) where bid=$2", dishToBasket.DishID, basketID)
	if err != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Error while adding dish to basket")
		return errors.BadRequestError("Error while adding dish to basket")
	}

	return nil
}

func (o orderRepo) DeleteFromBasket(ctx context.Context, dish models.DishToBasket, uid int) error {
	var basketID int
	err := o.DB.QueryRow("select basketID from basket_users where userID = $1", uid).Scan(&basketID)
	if err != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Error with basket through user")
		return errors.BadRequestError("Error with basket through user")
	}

	var number int
	err = o.DB.QueryRow("select number from baskets_food where dish = $1 and basket = $2", dish.DishID, basketID).Scan(&number)
	if err != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Error get dish count in basket")
		return errors.BadRequestError("Error get dish count in basket")
	}

	if number == 1 {
		_, err = o.DB.Exec("delete from baskets_food where basket = $1 and dish = $2", basketID, dish.DishID)
		if err != nil {
			logger.RepoLevel().InlineInfoLog(ctx, "Error with deleting dish from basket")
			return errors.BadRequestError("Error with deleting dish from basket")
		}

		return nil
	}

	_, err = o.DB.Exec("update baskets_food set number=number-1 where basket = $1 and dish = $2", basketID, dish.DishID)
	if err != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Error with deleting dish from basket")
		return errors.BadRequestError("Error with deleting dish from basket")
	}

	return nil
}

// TODO транзакция
func (o orderRepo) Create(ctx context.Context, uid int, orderParams models.CreateOrder) error {
	var basketID int
	var basketRestaurant string
	// находим что за корзина и из какого ресторана привязана к юзеру
	err := o.DB.QueryRow("select basketID from basket_users where userID = $1", uid).Scan(&basketID)
	if err != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Error with getting restaurant name")
		return errors.BadRequestError("Error with getting restaurant name")
	}

	// ищем ресторан по корзине
	err = o.DB.QueryRow("select restaurant from baskets where bid = $1", basketID).Scan(&basketRestaurant)

	// цена доставки ресторана для формирования заказа
	var deliveryCost int
	err = o.DB.QueryRow("select deliverycost from restaurants where name=$1", basketRestaurant).Scan(&deliveryCost)
	if err != nil {
		fmt.Println(err)
		logger.RepoLevel().InlineInfoLog(ctx, "Error with getting delivery cost")
		return errors.BadRequestError("Error with getting delivery cost")
	}

	// формируем в бд новый заказ
	var orderID int
	err = o.DB.QueryRow("insert into orders (restaurant, userID, ordertime, address, deliverycost, sum, status, deliverytime)"+
		"values ($1,$2,$3,$4,$5,$6,$7,$8) returning oid;", basketRestaurant, uid, time.Now(), orderParams.Address, deliveryCost, 0,
		models.StatusOrderAdded, time.Now()).Scan(&orderID) // todo решить что с временем доставки
	if err != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Error with inserting order in DB")
		return errors.BadRequestError("Error with inserting order in DB")
	}

	// удаляем связь корзина-юзер
	_, err = o.DB.Exec("delete from basket_users where basketid = $1", basketID)
	if err != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Error with inserting order in DB")
		return errors.BadRequestError("Error with inserting order in DB")
	}

	// создаем связь корзина-заказ
	_, err = o.DB.Exec("insert into basket_orders (basketid, orderid) values ($1, $2)", basketID, orderID)
	if err != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Error with inserting order in DB")
		return errors.BadRequestError("Error with inserting order in DB")
	}

	return nil
}

func (o orderRepo) GetUserOrders(ctx context.Context, uid int) ([]models.Order, error) {
	ordersDB, err := o.DB.Query("select oid, restaurant, orderTime, address, deliverycost, sum, status, deliverytime "+
		"from orders where userID=$1", uid)
	if err != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Error with getting restaurant orders")
		return nil, errors.BadRequestError("Error with getting restaurant orders")
	}
	orders := make([]models.Order, 0)
	for ordersDB.Next() {
		order := new(models.Order)

		var deliveryTime, orderTime time.Time
		err = ordersDB.Scan(
			&order.OID,
			&order.Restaurant,
			&orderTime,
			&order.Address,
			&order.DeliveryCost,
			&order.Summary,
			&order.Status,
			&deliveryTime,
		)
		order.OrderTime = orderTime.Format(config.TimeFormat)
		order.DeliveryTime = deliveryTime.Format(config.TimeFormat)

		var basketID string
		err = o.DB.QueryRow("select basketid from basket_orders where orderid=$1", order.OID).Scan(&basketID)
		if err != nil {
			logger.RepoLevel().InlineInfoLog(ctx, "Error with getting order's dishes")
			return nil, errors.BadRequestError("Error with getting order's dishes")
		}

		dishesDB, errr := o.DB.Query("select d.name, d.price, d.image, bf.number from baskets_food bf join dishes d on d.did = bf.dish where bf.basket=$1", basketID)
		if errr != nil {
			logger.RepoLevel().InlineInfoLog(ctx, "Error with getting order's dishes")
			return nil, errors.BadRequestError("Error with getting order's dishes")
		}

		dishes := make([]models.DishInOrder, 0)
		sum := 0
		for dishesDB.Next() {
			dish := new(models.DishInOrder)
			err = dishesDB.Scan(
				&dish.Name,
				&dish.Price,
				&dish.Image,
				&dish.Number)
			sum += dish.Number * dish.Price
			dishes = append(dishes, *dish)
		}
		order.Foods = dishes
		order.Summary = sum

		var restaurantImage string
		err = o.DB.QueryRow("select avatar from restaurants where name=$1", order.Restaurant).Scan(&restaurantImage)
		order.RestaurantImage = restaurantImage

		orders = append(orders, *order)
	}

	return orders, nil
}

func (o orderRepo) GetRestaurantOrders(ctx context.Context, restaurantName string) ([]models.Order, error) {
	ordersDB, err := o.DB.Query("select oid, userID, orderTime, address, deliverycost, sum, status, deliverytime from orders where restaurant=$1", restaurantName)
	if err != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Error with getting restaurant orders")
		return nil, errors.BadRequestError("Error with getting restaurant orders")
	}

	orders := make([]models.Order, 0)
	for ordersDB.Next() {
		order := new(models.Order)
		err = ordersDB.Scan(
			&order.OID,
			&order.UID,
			&order.OrderTime,
			&order.Address,
			&order.DeliveryCost,
			&order.Summary,
			&order.Status,
			&order.DeliveryTime,
		)

		var basketID string
		err = o.DB.QueryRow("select basketid from basket_orders where orderid=$1", order.OID).Scan(&basketID)
		if err != nil {
			logger.RepoLevel().InlineInfoLog(ctx, "Error with getting order's dishes")
			return nil, errors.BadRequestError("Error with getting order's dishes")
		}

		dishesDB, errr := o.DB.Query("select d.name, d.price, d.image, bf.number from baskets_food bf join dishes d on d.did = bf.dish where bf.basket=$1", basketID)
		if errr != nil {
			logger.RepoLevel().InlineInfoLog(ctx, "Error with getting order's dishes")
			return nil, errors.BadRequestError("Error with getting order's dishes")
		}

		dishes := make([]models.DishInOrder, 0)
		sum := 0
		for dishesDB.Next() {
			dish := new(models.DishInOrder)
			err = dishesDB.Scan(
				&dish.Name,
				&dish.Price,
				&dish.Image,
				&dish.Number)
			sum += dish.Number * dish.Price
			dishes = append(dishes, *dish)
		}
		order.Foods = dishes
		order.Summary = sum

		orders = append(orders, *order)
	}

	return orders, nil
}

func (o orderRepo) GetBasket(ctx context.Context, uid int) (models.BasketForUser, error) {
	var basketRestaurant, imageRestaurant string
	var basketID, restaurantID, deliveryCost int
	err := o.DB.QueryRow("select basketID from basket_users where userID = $1", uid).Scan(&basketID)
	if err == sql.ErrNoRows {
		return models.BasketForUser{}, nil
	}
	if err != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Error with getting basket through user")
		return models.BasketForUser{}, errors.BadRequestError("Error with getting basket through user")
	}

	basketResponse := models.BasketForUser{
		BID: basketID,
	}

	err = o.DB.QueryRow("select restaurant from baskets where bid = $1", basketID).Scan(&basketRestaurant)
	if err != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Error with finding restaurant through basket")
		return models.BasketForUser{}, errors.BadRequestError("Error with finding restaurant through basket")
	}
	basketResponse.Restaurant = basketRestaurant

	err = o.DB.QueryRow("select avatar from restaurants where name = $1", basketRestaurant).Scan(&imageRestaurant)
	if err != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Error with getting restaurant image")
		return models.BasketForUser{}, errors.BadRequestError("Error with getting restaurant image")
	}
	basketResponse.RestaurantImage = imageRestaurant

	err = o.DB.QueryRow("select rid, deliverycost from restaurants where name = $1", basketRestaurant).Scan(&restaurantID, &deliveryCost)
	if err != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Error with finding restaurantID through name")
		return models.BasketForUser{}, errors.BadRequestError("Error with finding restaurantID through name")
	}
	basketResponse.RID = restaurantID
	basketResponse.DeliveryCost = deliveryCost

	dishesDB, errr := o.DB.Query("select d.did, d.name, d.price, bf.number, d.image from baskets_food bf join dishes d on d.did = bf.dish where bf.basket=$1", basketID)
	if errr != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Error with getting basket's dishes")
		return models.BasketForUser{}, errors.BadRequestError("Error with getting basket's dishes")
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

	return basketResponse, nil
}
