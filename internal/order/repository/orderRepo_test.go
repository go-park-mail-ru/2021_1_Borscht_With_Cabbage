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

func TestOrderRepo_Create_GetBasketError(t *testing.T) {
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
		WillReturnError(sql.ErrNoRows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	orderParams := models.CreateOrder{
		Address: "Prospekt mira 23",
	}
	err = orderRepo.Create(ctx, 1, orderParams)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestOrderRepo_Create_GetRestError(t *testing.T) {
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
		ExpectQuery("select restaurant from baskets").
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	orderParams := models.CreateOrder{
		Address: "Prospekt mira 23",
	}
	err = orderRepo.Create(ctx, 1, orderParams)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestOrderRepo_Create_GetDelivCostError(t *testing.T) {
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
		WillReturnError(sql.ErrNoRows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	orderParams := models.CreateOrder{
		Address: "Prospekt mira 23",
	}
	err = orderRepo.Create(ctx, 1, orderParams)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestOrderRepo_Create_InsertError(t *testing.T) {
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
		WillReturnError(sql.ErrNoRows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	orderParams := models.CreateOrder{
		Address: "Prospekt mira 23",
	}
	err = orderRepo.Create(ctx, 1, orderParams)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestOrderRepo_Create_DeleteUserBindError(t *testing.T) {
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
		WillReturnError(sql.ErrNoRows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	orderParams := models.CreateOrder{
		Address: "Prospekt mira 23",
	}
	err = orderRepo.Create(ctx, 1, orderParams)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestOrderRepo_Create_InsertBasketError(t *testing.T) {
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
		WillReturnError(sql.ErrNoRows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	orderParams := models.CreateOrder{
		Address: "Prospekt mira 23",
	}
	err = orderRepo.Create(ctx, 1, orderParams)
	if err == nil {
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

	orders := sqlmock.NewRows([]string{"oid", "restaurant", "orderTime", "address", "deliverycost", "sum", "status", "deliverytime", "review", "stars"})
	expectOrders := []models.Order{
		{OID: 1, Restaurant: "rest1", OrderTime: "15:00", Address: "Prospekt mira 2", DeliveryCost: 200, Summary: 1200, Status: models.StatusOrderAdded, DeliveryTime: "17:00"},
	}
	for _, item := range expectOrders {
		orders = orders.AddRow(item.OID, item.Restaurant, item.OrderTime, item.Address, item.DeliveryCost, item.Summary, item.Status, item.DeliveryTime, item.Review, item.Stars)
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
	basketID = basketID.AddRow(1)

	user := sqlmock.NewRows([]string{"name", "phone"})
	user = user.AddRow("Dasha", 89111111111)

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
		ExpectQuery("select name, phone").
		WithArgs(expectOrders[0].UID).
		WillReturnRows(user)
	mock.
		ExpectQuery("select basketid from").
		WithArgs(expectOrders[0].OID).
		WillReturnRows(basketID)
	mock.
		ExpectQuery("select d.name, d.price,").
		WithArgs("1").
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

func TestOrderRepo_SetNewStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	orderRepo := &orderRepo{
		DB: db,
	}

	newStatus := models.SetNewStatus{
		Status:       models.StatusOrderDone,
		DeliveryTime: "2012-11-01T22:08:41",
		Restaurant:   "rest1",
		OID:          1,
	}

	mock.
		ExpectExec("UPDATE orders SET").
		WillReturnResult(sqlmock.NewResult(1, 1))

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	err = orderRepo.SetNewStatus(ctx, newStatus)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestOrderRepo_SetNewStatus_ErrorConvertingTime(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	orderRepo := &orderRepo{
		DB: db,
	}

	newStatus := models.SetNewStatus{
		Status:       models.StatusOrderDone,
		DeliveryTime: "2021-05-01 18:00:00",
		Restaurant:   "rest1",
		OID:          1,
	}

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	err = orderRepo.SetNewStatus(ctx, newStatus)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestOrderRepo_CreateReview(t *testing.T) {
}
