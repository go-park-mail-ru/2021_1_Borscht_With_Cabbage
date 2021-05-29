package repository

import (
	"context"
	"database/sql"
	"github.com/borscht/backend/utils/calcDistance"
	"time"

	"github.com/borscht/backend/configProject"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/order"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
)

type orderRepo struct {
	DB *sql.DB
}

func NewOrderRepo(db *sql.DB) order.OrderRepo {
	return &orderRepo{
		DB: db,
	}
}

// TODO транзакция
func (o orderRepo) Create(ctx context.Context, uid int, orderParams models.CreateOrder, coordinates models.Coordinates) error {
	var basketID, restID int
	var basketRestaurant string
	// находим что за корзина и из какого ресторана привязана к юзеру
	err := o.DB.QueryRow("select basketID, restaurantid from basket_users where basketID = $1", orderParams.BasketID).Scan(&basketID, &restID)
	if err != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Error with getting restaurant name")
		return errors.BadRequestError("Error with getting restaurant name")
	}

	var latitude, longitude float64
	var radius int
	err = o.DB.QueryRow("select latitude, longitude, radius from addresses where rid=$1", restID).Scan(&latitude, &longitude, &radius)
	if err != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Error with getting restaurant address")
		return errors.BadRequestError("Error with getting restaurant address")
	}
	deliveryTime := calcDistance.GetDeliveryTime(coordinates.Latitude, coordinates.Longitude, latitude, longitude, radius)
	doneTime := time.Now().Add(time.Duration(deliveryTime) * time.Minute)

	// и для изменения кол-ва заказов и общей суммы заказа
	var deliveryCost, ordersSum, ordersCount int
	err = o.DB.QueryRow("select name, deliverycost, orderscount, orderssum from restaurants where rid=$1", restID).
		Scan(&basketRestaurant, &deliveryCost, &ordersCount, &ordersSum)
	if err != nil {
		failErr := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failErr)
		return failErr
	}

	// сумма текущего заказа
	var orderSum int
	err = o.DB.QueryRow("select sum(d.price) from baskets_food bf join dishes d on d.did = bf.dish where bf.basket=$1", basketID).
		Scan(&orderSum)
	if err != nil {
		failErr := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failErr)
		return failErr
	}

	ordersSum += orderSum + deliveryCost
	ordersCount++

	// сохраняем кол-во заказов и сумарная цена заказов
	_, err = o.DB.Exec("UPDATE restaurants SET orderscount=$1, orderssum=$2 WHERE rid=$3",
		ordersCount, ordersSum, restID)
	if err != nil {
		failErr := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failErr)
		return failErr
	}

	// формируем в бд новый заказ
	var orderID int
	err = o.DB.QueryRow("insert into orders (restaurant, userID, ordertime, address, deliverycost, sum, status, deliverytime)"+
		"values ($1,$2,$3,$4,$5,$6,$7,$8) returning oid;", basketRestaurant, uid, time.Now(), orderParams.Address, deliveryCost, 0,
		models.StatusOrderAdded, doneTime).Scan(&orderID) // todo решить что с временем доставки
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
	query :=
		`
		SELECT o.oid, o.restaurant, o.orderTime, o.address, o.deliverycost, o.sum, o.status, o.deliverytime, o.review, o.stars, r.rid
		FROM orders o JOIN restaurants r
		ON o.restaurant = r.name 
		WHERE userID=$1 ORDER BY orderTime DESC
	`
	ordersDB, err := o.DB.Query(query, uid)
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
			&order.Review,
			&order.Stars,
			&order.RID,
		)

		logger.RepoLevel().DebugLog(ctx, logger.Fields{"order": order, "rid": order.RID})

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
		order.Summary = sum + order.DeliveryCost

		var restaurantImage string

		var rid int
		err = o.DB.QueryRow("select avatar, rid from restaurants where name=$1", order.Restaurant).Scan(&restaurantImage, &rid)
		if err != nil {
			logger.RepoLevel().InlineInfoLog(ctx, "Error with getting user's info")
			return nil, errors.BadRequestError("Error with getting user's inf")
		}

		order.RestaurantImage = restaurantImage
		order.RID = rid

		orders = append(orders, *order)
	}

	return orders, nil
}

func (o orderRepo) GetRestaurantOrders(ctx context.Context, restaurantName string) ([]models.Order, error) {
	ordersDB, err := o.DB.Query("select oid, userID, orderTime, address, deliverycost, sum, status, deliverytime from orders where restaurant=$1 "+
		"order by orderTime desc", restaurantName)
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

		err = o.DB.QueryRow("select name, phone from users where uid=$1", order.UID).Scan(&order.UserName, &order.UserPhone)
		if err != nil {
			logger.RepoLevel().InlineInfoLog(ctx, "Error with getting user's info")
			return nil, errors.BadRequestError("Error with getting user's info")
		}

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

func (o orderRepo) SetNewStatus(ctx context.Context, newStatus models.SetNewStatus) (int, error) {
	timeToDB, err := time.Parse(configProject.TimeFormat, newStatus.DeliveryTime)
	if err != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Error while converting time")
		return 0, errors.BadRequestError("Error while converting time")
	}

	var uid int
	err = o.DB.QueryRow("UPDATE orders SET status=$1, deliverytime=$2 where restaurant=$3 and oid=$4 returning userid",
		newStatus.Status, timeToDB, newStatus.Restaurant, newStatus.OID).Scan(&uid)
	if err != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Error with updating order status in DB")
		return 0, errors.BadRequestError("Error with updating order status in DB")
	}

	return uid, nil
}

func (o orderRepo) CreateReview(ctx context.Context, newReview models.SetNewReview) error {
	var restaurant string
	err := o.DB.QueryRow("UPDATE orders SET review=$1, stars=$2 WHERE oid=$3 returning restaurant",
		newReview.Review, newReview.Stars, newReview.OID).Scan(&restaurant)
	if err != nil {
		logger.RepoLevel().InlineInfoLog(ctx, "Error with setting order review in DB")
		return errors.BadRequestError("Error with setting order review in DB")
	}

	_, err = o.DB.Exec("UPDATE restaurants SET ratingsSum=ratingsSum+$1, reviewsCount=reviewsCount+1 where name=$2",
		newReview.Stars, restaurant)

	return nil
}
