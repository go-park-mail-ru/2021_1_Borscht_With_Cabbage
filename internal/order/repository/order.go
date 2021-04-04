package repository

import (
	"context"
	"database/sql"
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
	err := o.DB.QueryRow("select bid, restaurant from baskets where userid = $1", uid).Scan(&basketID, &basketRestaurant)
	if basketRestaurant == "" {
		//если у текущей корзины нет ресторана,то мы записываем в нее ресторан, блюдо из которого добавляем сейчас
		var restaurantName string
		err = o.DB.QueryRow("select restaurant from dishes where did = $1", dishToBasket.DishID).Scan(&restaurantName)
		if err != nil {
			logger.RepoLevel().InlineInfoLog(ctx, "Error with finding restaurant through dish")
			return errors.BadRequestError("Error with finding restaurant through dish")
		}
		_, err = o.DB.Exec("update baskets set restaurant = $1 where bid = $2", restaurantName, basketID)
		if err != nil {
			logger.RepoLevel().InlineInfoLog(ctx, "Error with set restaurant name to basket")
			return errors.BadRequestError("Error with set restaurant name to basket")
		}
	}

	// если добавляем в корзину с тем же рестораном, то просто добавляем блюдо
	if dishToBasket.SameBasket {
		_, err = o.DB.Exec("insert into baskets_food (dish, basket) values ($1, $2)", dishToBasket.DishID, basketID)
		if err != nil {
			logger.RepoLevel().InlineInfoLog(ctx, "Error while adding dish to basket")
			return errors.BadRequestError("Error while adding dish to basket")
		}
		return nil
	}

	// если ресeторан теперь другой, то добавляем блюдо и чистим старую корзину
	_, err = o.DB.Exec("delete from baskets_food where basket = $1", basketID)
	if err != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Error with deleting previous dishes from basket")
		return errors.BadRequestError("Error with deleting previous dishes from basket")
	}

	_, err = o.DB.Exec("insert into baskets_food (basket, dish) values ($1, $2)", basketID, dishToBasket.DishID)
	if err != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Error while adding dish to basket")
		return errors.BadRequestError("Error while adding dish to basket")
	}
	return nil
}

func (o orderRepo) Create(ctx context.Context, uid int, orderParams models.CreateOrder) error {
	var basketID int
	var basketRestaurant string
	// находим что за корзина и из какого ресторана привязана к юзеру
	err := o.DB.QueryRow("select bid, restaurant from baskets where user=$1", uid).Scan(&basketID, &basketRestaurant)
	if err != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Error with getting restaurant name")
		return errors.BadRequestError("Error with getting restaurant name")
	}

	// цена доставки ресторана для формирования заказа
	var deliveryCost int
	err = o.DB.QueryRow("select deliverycost from restaurants where name=$1", basketRestaurant).Scan(&deliveryCost)
	if err != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Error with getting delivery cost")
		return errors.BadRequestError("Error with getting delivery cost")
	}

	// формируем в бд новый заказ
	var orderID int
	err = o.DB.QueryRow("insert into orders (restaurant, user, ordertime, address, deliverycost, sum, status, deliverytime)"+
		"values ($1,$2,$3,$4,$5,$6,$7,$8) returning oid;", basketRestaurant, uid, time.Now(), orderParams.Address, deliveryCost, 0,
		models.StatusOrderAdded, time.Now()).Scan(orderID) //todo изменить время
	if err != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Error with inserting order in DB")
		return errors.BadRequestError("Error with inserting order in DB")
	}

	// теперь корзина привязана к сформированному заказу, а не пользователю
	_, err = o.DB.Exec("update baskets set userid=null, orderid=$1 where user=$2", orderID, uid)

	return nil
}

func (o orderRepo) GetUserOrders(ctx context.Context, uid int) ([]models.Order, error) {
	//var oid int
	//err := o.DB.QueryRow("insert into orders (restaurant, userID, ordertime, address, deliverycost, sum, status, deliverytime)"+
	//	" values ($1,$2,$3,$4,$5,$6,$7,$8) returning oid", "yum", 1, time.Now(), "бауманская 2", 200, 1900, "едет к вам", time.Now()).Scan(oid)
	//fmt.Println(err)
	//fmt.Println(oid)
	//return nil, err

	ordersDB, err := o.DB.Query("select oid, restaurant, orderTime, address, deliverycost, sum, status, deliverytime "+
		"from orders where userID=$1", uid)
	if err != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Error with getting restaurant orders")
		return nil, errors.BadRequestError("Error with getting restaurant orders")
	}
	orders := make([]models.Order, 0)
	for ordersDB.Next() {
		order := new(models.Order)
		err = ordersDB.Scan(
			&order.OID,
			&order.Restaurant,
			&order.OrderTime,
			&order.Address,
			&order.DeliveryCost,
			&order.Summary,
			&order.Status,
			&order.DeliveryTime,
		)
		orders = append(orders, *order)
	}

	return orders, nil
}

func (o orderRepo) GetRestaurantOrders(ctx context.Context, restaurantName string) ([]models.Order, error) {
	ordersDB, err := o.DB.Query("select oid, userID, orderTime, address, deliverycost, sum, status, deliverytime "+
		"from orders where restaurant=$1", restaurantName)
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
		orders = append(orders, *order)
	}

	return orders, nil
}
