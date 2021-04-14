package repository

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/borscht/backend/internal/models"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"testing"
)

func TestNewOrderRepo(t *testing.T) {

}

func TestOrderRepo_AddToNewBasket(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	orderRepo := &orderRepo{
		DB: db,
	}

	restaurantName := sqlmock.NewRows([]string{"restaurant"})
	expectRestName := []models.RestaurantInfo{
		{Title: "rest1"},
	}
	for _, item := range expectRestName {
		restaurantName = restaurantName.AddRow(item.Title)
	}

	basketID := sqlmock.NewRows([]string{"restaurant"})
	expectBasketID := []models.BasketForUser{
		{BID: 1},
	}
	for _, item := range expectBasketID {
		basketID = basketID.AddRow(item.BID)
	}

	dishID := sqlmock.NewRows([]string{"dish"})
	expectDishID := []models.Dish{
		{ID: 1},
	}
	for _, item := range expectDishID {
		dishID = dishID.AddRow(item.ID)
	}

	// если к юзеру пока не привязана корзина - создаем
	mock.
		ExpectQuery("select basketID from").
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)
	mock.
		ExpectQuery("select restaurant from").
		WithArgs(1).
		WillReturnRows(restaurantName)
	mock.
		ExpectQuery("insert into baskets").
		WithArgs("rest1").
		WillReturnRows(basketID)
	mock.
		ExpectExec("insert into basket_users").
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(1, 2))

	// добавляем в корзину
	mock.
		ExpectQuery("select dish from baskets_food").
		WithArgs(1, 1).
		WillReturnError(sql.ErrNoRows).WillReturnRows(dishID)
	mock.
		ExpectExec("insert into baskets_food").
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(1, 3))

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	dishToBasket := models.DishToBasket{
		DishID:     1,
		SameBasket: true,
		IsPlus:     true,
	}
	err = orderRepo.AddToBasket(ctx, dishToBasket, 1)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestOrderRepo_AddToBasket_SameRestaurant(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	orderRepo := &orderRepo{
		DB: db,
	}

	restaurantName := sqlmock.NewRows([]string{"restaurant"})
	expectRestName := []models.RestaurantInfo{
		{Title: "rest1"},
	}
	for _, item := range expectRestName {
		restaurantName = restaurantName.AddRow(item.Title)
	}

	basketID := sqlmock.NewRows([]string{"restaurant"})
	expectBasketID := []models.BasketForUser{
		{BID: 1},
	}
	for _, item := range expectBasketID {
		basketID = basketID.AddRow(item.BID)
	}

	dishID := sqlmock.NewRows([]string{"dish"})
	expectDishID := []models.Dish{
		{ID: 1},
	}
	for _, item := range expectDishID {
		dishID = dishID.AddRow(item.ID)
	}

	mock.
		ExpectQuery("select basketID from").
		WithArgs(1).
		WillReturnRows(basketID)

	// добавляем в корзину
	mock.
		ExpectQuery("select dish from baskets_food").
		WithArgs(1, 1).
		WillReturnError(sql.ErrNoRows).WillReturnRows(dishID)
	mock.
		ExpectExec("insert into baskets_food").
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(1, 3))

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	dishToBasket := models.DishToBasket{
		DishID:     1,
		SameBasket: true,
		IsPlus:     true,
	}
	err = orderRepo.AddToBasket(ctx, dishToBasket, 1)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestOrderRepo_AddToBasketSameDish_SameRestaurant(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	orderRepo := &orderRepo{
		DB: db,
	}

	basketID := sqlmock.NewRows([]string{"restaurant"})
	expectBasketID := []models.BasketForUser{
		{BID: 1},
	}
	for _, item := range expectBasketID {
		basketID = basketID.AddRow(item.BID)
	}

	dishID := sqlmock.NewRows([]string{"dish"})
	expectDishID := []models.Dish{
		{ID: 1},
	}
	for _, item := range expectDishID {
		dishID = dishID.AddRow(item.ID)
	}

	mock.
		ExpectQuery("select basketID from").
		WithArgs(1).
		WillReturnRows(basketID)

	// добавляем в корзину
	mock.
		ExpectQuery("select dish from baskets_food").
		WithArgs(1, 1).
		WillReturnError(nil).WillReturnRows(dishID)
	mock.
		ExpectExec("update baskets_food").
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	dishToBasket := models.DishToBasket{
		DishID:     1,
		SameBasket: true,
		IsPlus:     true,
	}
	err = orderRepo.AddToBasket(ctx, dishToBasket, 1)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestOrderRepo_AddToBasket_NewRestaurant(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	orderRepo := &orderRepo{
		DB: db,
	}

	basketID := sqlmock.NewRows([]string{"restaurant"})
	expectBasketID := []models.BasketForUser{
		{BID: 1},
	}
	for _, item := range expectBasketID {
		basketID = basketID.AddRow(item.BID)
	}

	mock.
		ExpectQuery("select basketID from").
		WithArgs(1).
		WillReturnRows(basketID)
	mock.
		ExpectQuery("select dish from baskets_food").
		WithArgs(1, 1).
		WillReturnError(sql.ErrNoRows)
	mock.
		ExpectExec("insert into baskets_food").
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	dishToBasket := models.DishToBasket{
		DishID:     1,
		SameBasket: true,
		IsPlus:     true,
	}
	err = orderRepo.AddToBasket(ctx, dishToBasket, 1)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestOrderRepo_DeleteFromBasket(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	orderRepo := &orderRepo{
		DB: db,
	}

	basketID := sqlmock.NewRows([]string{"restaurant"})
	expectBasketID := []models.BasketForUser{
		{BID: 1},
	}
	for _, item := range expectBasketID {
		basketID = basketID.AddRow(item.BID)
	}

	dishNum := sqlmock.NewRows([]string{"dish"})
	expectDishNum := []models.DishInOrder{
		{Number: 1},
	}
	for _, item := range expectDishNum {
		dishNum = dishNum.AddRow(item.Number)
	}

	mock.
		ExpectQuery("select basketID from").
		WithArgs(1).
		WillReturnRows(basketID)

	mock.
		ExpectQuery("select number from baskets_food").
		WithArgs(1, 1).
		WillReturnRows(dishNum)
	mock.
		ExpectExec("delete from baskets_food").
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	dishToBasket := models.DishToBasket{
		DishID:     1,
		SameBasket: true,
		IsPlus:     false,
	}
	err = orderRepo.DeleteFromBasket(ctx, dishToBasket, 1)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestOrderRepo_DecDishCountInBasket(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	orderRepo := &orderRepo{
		DB: db,
	}

	basketID := sqlmock.NewRows([]string{"restaurant"})
	expectBasketID := []models.BasketForUser{
		{BID: 1},
	}
	for _, item := range expectBasketID {
		basketID = basketID.AddRow(item.BID)
	}

	dishNum := sqlmock.NewRows([]string{"dish"})
	expectDishNum := []models.DishInOrder{
		{Number: 2},
	}
	for _, item := range expectDishNum {
		dishNum = dishNum.AddRow(item.Number)
	}

	mock.
		ExpectQuery("select basketID from").
		WithArgs(1).
		WillReturnRows(basketID)

	mock.
		ExpectQuery("select number from baskets_food").
		WithArgs(1, 1).
		WillReturnRows(dishNum)
	mock.
		ExpectExec("update baskets_food set").
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	dishToBasket := models.DishToBasket{
		DishID:     1,
		SameBasket: true,
		IsPlus:     false,
	}
	err = orderRepo.DeleteFromBasket(ctx, dishToBasket, 1)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestOrderRepo_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	orderRepo := &orderRepo{
		DB: db,
	}

	basketID := sqlmock.NewRows([]string{"restaurant"})
	expectBasketID := []models.BasketForUser{
		{BID: 1},
	}
	for _, item := range expectBasketID {
		basketID = basketID.AddRow(item.BID)
	}

	basketRestaurant := sqlmock.NewRows([]string{"dish"})
	expectBasketRestaurant := []models.RestaurantInfo{
		{Title: "rest1"},
	}
	for _, item := range expectBasketRestaurant {
		basketRestaurant = basketRestaurant.AddRow(item.Title)
	}

	deliveryCost := sqlmock.NewRows([]string{"dish"})
	expectDeliveryCost := []models.RestaurantInfo{
		{DeliveryCost: 200},
	}
	for _, item := range expectDeliveryCost {
		deliveryCost = deliveryCost.AddRow(item.DeliveryCost)
	}

	orderID := sqlmock.NewRows([]string{"dish"})
	expectOrderID := []models.Order{
		{OID: 1},
	}
	for _, item := range expectOrderID {
		orderID = orderID.AddRow(item.OID)
	}

	mock.
		ExpectQuery("select basketID from").
		WithArgs(1).
		WillReturnRows(basketID)
	mock.
		ExpectQuery("select restaurant from baskets").
		WithArgs(1).
		WillReturnRows(basketRestaurant)
	mock.
		ExpectQuery("select deliverycost from restaurants").
		WithArgs("rest1").
		WillReturnRows(deliveryCost)
	mock.
		ExpectQuery("insert into orders").
		WillReturnRows(orderID)
	mock.
		ExpectExec("delete from basket_users").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.
		ExpectExec("insert into basket_orders").
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	orderParams := models.CreateOrder{
		Address: "Prospekt mira 23",
	}
	err = orderRepo.Create(ctx, 1, orderParams)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestOrderRepo_GetUserOrders(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	orderRepo := &orderRepo{
		DB: db,
	}

	orders := sqlmock.NewRows([]string{"oid", "restaurant", "orderTime", "address", "deliverycost", "sum", "status", "deliverytime"})
	expectOrders := []models.Order{
		{OID: 1, Restaurant: "rest1", OrderTime: "15:00", Address: "Prospekt mira 2", DeliveryCost: 200, Summary: 1200, Status: models.StatusOrderAdded, DeliveryTime: "17:00"},
	}
	for _, item := range expectOrders {
		orders = orders.AddRow(item.OID, item.Restaurant, item.OrderTime, item.Address, item.DeliveryCost, item.Summary, item.Status, item.DeliveryTime)
	}

	dishes := sqlmock.NewRows([]string{"name", "price", "image", "number"})
	expectDishes := []models.DishInOrder{
		{Name: "dish1", Price: 250, Image: "img.jpg", Number: 1},
		{Name: "dish2", Price: 350, Image: "img2.jpg", Number: 2},
	}
	for _, item := range expectDishes {
		dishes = dishes.AddRow(item.Name, item.Price, item.Image, item.Number)
	}

	basketID := sqlmock.NewRows([]string{"dish"})
	expectBasketID := []models.BasketForUser{
		{BID: 1},
	}
	for _, item := range expectBasketID {
		basketID = basketID.AddRow(item.BID)
	}

	restaurantAvatar := sqlmock.NewRows([]string{"dish"})
	expectRestaurantAvatar := []models.RestaurantInfo{
		{Avatar: "img.jpg"},
	}
	for _, item := range expectRestaurantAvatar {
		restaurantAvatar = restaurantAvatar.AddRow(item.Avatar)
	}

	mock.
		ExpectQuery("select oid, restaurant,").
		WithArgs(1).
		WillReturnRows(orders)
	mock.
		ExpectQuery("select basketid from").
		WithArgs(1).
		WillReturnRows(basketID)
	mock.
		ExpectQuery("select d.name, d.price,").
		WillReturnRows(dishes)
	mock.
		ExpectQuery("select avatar from restaurants").
		WithArgs("rest1").
		WillReturnRows(restaurantAvatar)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	ordersResult, errr := orderRepo.GetUserOrders(ctx, 1)
	if errr != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	require.EqualValues(t, ordersResult[0].OID, 1)
}

func TestOrderRepo_GetRestaurantOrders(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	orderRepo := &orderRepo{
		DB: db,
	}

	orders := sqlmock.NewRows([]string{"oid", "userid", "orderTime", "address", "deliverycost", "sum", "status", "deliverytime"})
	expectOrders := []models.Order{
		{OID: 1, UID: 1, OrderTime: "15:00", Address: "Prospekt mira 2", DeliveryCost: 200, Summary: 1200, Status: models.StatusOrderAdded, DeliveryTime: "17:00"},
	}
	for _, item := range expectOrders {
		orders = orders.AddRow(item.OID, item.UID, item.OrderTime, item.Address, item.DeliveryCost, item.Summary, item.Status, item.DeliveryTime)
	}

	basketID := sqlmock.NewRows([]string{"dish"})
	expectBasketID := []models.BasketForUser{
		{BID: 1},
	}
	for _, item := range expectBasketID {
		basketID = basketID.AddRow(item.BID)
	}

	dishes := sqlmock.NewRows([]string{"name", "price", "image", "number"})
	expectDishes := []models.DishInOrder{
		{Name: "dish1", Price: 250, Image: "img.jpg", Number: 1},
		{Name: "dish2", Price: 350, Image: "img2.jpg", Number: 2},
	}
	for _, item := range expectDishes {
		dishes = dishes.AddRow(item.Name, item.Price, item.Image, item.Number)
	}

	mock.
		ExpectQuery("select oid, userID,").
		WithArgs("rest1").
		WillReturnRows(orders)
	mock.
		ExpectQuery("select basketid from").
		WillReturnRows(basketID)
	mock.
		ExpectQuery("select d.name, d.price,").
		WillReturnRows(dishes)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	ordersResult, errr := orderRepo.GetRestaurantOrders(ctx, "rest1")
	if errr != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	require.EqualValues(t, ordersResult[0].OID, 1)
}

func TestOrderRepo_GetBasket(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	orderRepo := &orderRepo{
		DB: db,
	}

	basketID := sqlmock.NewRows([]string{"dish"})
	expectBasketID := []models.BasketForUser{
		{BID: 1},
	}
	for _, item := range expectBasketID {
		basketID = basketID.AddRow(item.BID)
	}

	restaurantName := sqlmock.NewRows([]string{"restaurant"})
	expectRestName := []models.RestaurantInfo{
		{Title: "rest1"},
	}
	for _, item := range expectRestName {
		restaurantName = restaurantName.AddRow(item.Title)
	}

	restaurantAvatar := sqlmock.NewRows([]string{"restaurant"})
	expectRestaurantAvatar := []models.RestaurantInfo{
		{Avatar: "img.jpg"},
	}
	for _, item := range expectRestaurantAvatar {
		restaurantAvatar = restaurantAvatar.AddRow(item.Avatar)
	}

	restaurantInfo := sqlmock.NewRows([]string{"rid", "deliveryCost"})
	expectRestaurantInfo := []models.RestaurantInfo{
		{ID: 1, DeliveryCost: 200},
	}
	for _, item := range expectRestaurantInfo {
		restaurantInfo = restaurantInfo.AddRow(item.ID, item.DeliveryCost)
	}

	dishes := sqlmock.NewRows([]string{"id", "name", "price", "number", "image"})
	expectDishes := []models.DishInOrder{
		{ID: 1, Name: "dish1", Price: 250, Image: "img.jpg", Number: 1},
		{ID: 2, Name: "dish2", Price: 350, Image: "img2.jpg", Number: 2},
	}
	for _, item := range expectDishes {
		dishes = dishes.AddRow(item.ID, item.Name, item.Price, item.Number, item.Image)
	}

	mock.
		ExpectQuery("select basketID from basket_users").
		WithArgs(1).
		WillReturnRows(basketID)
	mock.
		ExpectQuery("select restaurant from baskets").
		WithArgs(1).
		WillReturnRows(restaurantName)
	mock.
		ExpectQuery("select avatar from restaurants").
		WithArgs("rest1").
		WillReturnRows(restaurantAvatar)
	mock.
		ExpectQuery("select rid, deliverycost from restaurants").
		WithArgs("rest1").
		WillReturnRows(restaurantInfo)
	mock.
		ExpectQuery("select d.did, d.name,").
		WithArgs(1).
		WillReturnRows(dishes)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	basketResult, errr := orderRepo.GetBasket(ctx, 1)
	if errr != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	require.EqualValues(t, basketResult.RID, 1)
	require.EqualValues(t, basketResult.BID, 1)
	require.EqualValues(t, basketResult.Restaurant, "rest1")
	require.EqualValues(t, basketResult.RestaurantImage, "img.jpg")
	require.EqualValues(t, basketResult.DeliveryCost, 200)
}
