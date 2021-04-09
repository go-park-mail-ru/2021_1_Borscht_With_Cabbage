package repository

import (
	"context"
	"database/sql"
	"fmt"
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

		fmt.Println(basketID, uid)
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
		_, err = o.DB.Exec("insert into baskets_food (dish, basket) values ($1, $2)", dishToBasket.DishID, basketID)
		if err != nil {
			logger.RepoLevel().InlineInfoLog(ctx, "Error while adding dish to basket")
			return errors.BadRequestError("Error while adding dish to basket")
		}
		return nil
	}

	// если ресторан теперь другой, то чистим старую корзину и добавляем блюдо
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
	err = o.DB.QueryRow("select restaurant from baskets where bid = $1", basketID).Scan(basketRestaurant)

	// цена доставки ресторана для формирования заказа
	var deliveryCost int
	err = o.DB.QueryRow("select deliverycost from restaurants where name=$1", basketRestaurant).Scan(&deliveryCost)
	if err != nil {
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

		dishesDB, errr := o.DB.Query("select d.name, d.price, d.image " +
			"from baskets_food bf" +
			"join dishes d where d.dish = bf.did")
		if errr != nil {
			logger.RepoLevel().InlineInfoLog(ctx, "Error with getting order's dishes")
			return nil, errors.BadRequestError("Error with getting order's dishes")
		}

		dishes := make([]models.DishInOrder, 0)
		for dishesDB.Next() {
			dish := new(models.DishInOrder)
			err = dishesDB.Scan(
				&dish.Name,
				&dish.Price,
				&dish.Image)
			dishes = append(dishes, *dish)
		}
		order.Foods = dishes

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
