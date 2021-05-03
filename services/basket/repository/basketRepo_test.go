package basketServiceRepo

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/borscht/backend/internal/models"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewBasketRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	defer db.Close()
	basketRepo := NewBasketRepository(db)
	if basketRepo != nil {
		return
	}
}

func TestOrderRepo_AddToNewBasket(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	basketRepo := &basketRepository{
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
	err = basketRepo.AddToBasket(ctx, dishToBasket, 1)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestOrderRepo_AddToNewBasket_GetRestError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	basketRepo := &basketRepository{
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
		WillReturnError(sql.ErrNoRows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	dishToBasket := models.DishToBasket{
		DishID:     1,
		SameBasket: true,
		IsPlus:     true,
	}
	err = basketRepo.AddToBasket(ctx, dishToBasket, 1)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestOrderRepo_AddToNewBasket_CreateBasketError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	basketRepo := &basketRepository{
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
		WillReturnError(sql.ErrNoRows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	dishToBasket := models.DishToBasket{
		DishID:     1,
		SameBasket: true,
		IsPlus:     true,
	}
	err = basketRepo.AddToBasket(ctx, dishToBasket, 1)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestOrderRepo_AddToNewBasket_BindUserError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	basketRepo := &basketRepository{
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
		WillReturnError(sql.ErrNoRows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	dishToBasket := models.DishToBasket{
		DishID:     1,
		SameBasket: true,
		IsPlus:     true,
	}
	err = basketRepo.AddToBasket(ctx, dishToBasket, 1)
	if err == nil {
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
	basketRepo := &basketRepository{
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
	err = basketRepo.AddToBasket(ctx, dishToBasket, 1)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestOrderRepo_AddToBasket_SameRestaurant_AddError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	basketRepo := &basketRepository{
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
		WillReturnError(sql.ErrNoRows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	dishToBasket := models.DishToBasket{
		DishID:     1,
		SameBasket: true,
		IsPlus:     true,
	}
	err = basketRepo.AddToBasket(ctx, dishToBasket, 1)
	if err == nil {
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
	basketRepo := &basketRepository{
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
	err = basketRepo.AddToBasket(ctx, dishToBasket, 1)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestOrderRepo_AddToBasketSameDish_SameRestaurant_AddError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	basketRepo := &basketRepository{
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
		WillReturnError(sql.ErrNoRows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	dishToBasket := models.DishToBasket{
		DishID:     1,
		SameBasket: true,
		IsPlus:     true,
	}
	err = basketRepo.AddToBasket(ctx, dishToBasket, 1)
	if err == nil {
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
	basketRepo := &basketRepository{
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
		ExpectExec("delete from baskets_food ").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.
		ExpectExec("insert into baskets_food").
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.
		ExpectExec("update baskets set restaurant").
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	dishToBasket := models.DishToBasket{
		DishID:     1,
		SameBasket: false,
		IsPlus:     true,
	}
	err = basketRepo.AddToBasket(ctx, dishToBasket, 1)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestOrderRepo_AddToBasket_NewRestaurant_DeleteOldError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	basketRepo := &basketRepository{
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
		ExpectExec("delete from baskets_food ").
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	dishToBasket := models.DishToBasket{
		DishID:     1,
		SameBasket: false,
		IsPlus:     true,
	}
	err = basketRepo.AddToBasket(ctx, dishToBasket, 1)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestOrderRepo_AddToBasket_NewRestaurant_InsertError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	basketRepo := &basketRepository{
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
		ExpectExec("delete from baskets_food ").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.
		ExpectExec("insert into baskets_food").
		WithArgs(1, 1).
		WillReturnError(sql.ErrNoRows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	dishToBasket := models.DishToBasket{
		DishID:     1,
		SameBasket: false,
		IsPlus:     true,
	}
	err = basketRepo.AddToBasket(ctx, dishToBasket, 1)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestOrderRepo_AddToBasket_NewRestaurant_UpdRestError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	basketRepo := &basketRepository{
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
		ExpectExec("delete from baskets_food ").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.
		ExpectExec("insert into baskets_food").
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.
		ExpectExec("update baskets set restaurant").
		WithArgs(1, 1).
		WillReturnError(sql.ErrNoRows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	dishToBasket := models.DishToBasket{
		DishID:     1,
		SameBasket: false,
		IsPlus:     true,
	}
	err = basketRepo.AddToBasket(ctx, dishToBasket, 1)
	if err == nil {
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
	basketRepo := &basketRepository{
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
	err = basketRepo.DeleteFromBasket(ctx, dishToBasket, 1)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestOrderRepo_DeleteFromBasket_GetBasketError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	basketRepo := &basketRepository{
		DB: db,
	}

	mock.
		ExpectQuery("select basketID from").
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	dishToBasket := models.DishToBasket{
		DishID:     1,
		SameBasket: true,
		IsPlus:     false,
	}
	err = basketRepo.DeleteFromBasket(ctx, dishToBasket, 1)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestOrderRepo_DeleteFromBasket_GetNumberError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	basketRepo := &basketRepository{
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
		ExpectQuery("select number from baskets_food").
		WithArgs(1, 1).
		WillReturnError(sql.ErrNoRows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	dishToBasket := models.DishToBasket{
		DishID:     1,
		SameBasket: true,
		IsPlus:     false,
	}
	err = basketRepo.DeleteFromBasket(ctx, dishToBasket, 1)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestOrderRepo_DeleteFromBasket_DeleteError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	basketRepo := &basketRepository{
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
		WillReturnError(sql.ErrNoRows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	dishToBasket := models.DishToBasket{
		DishID:     1,
		SameBasket: true,
		IsPlus:     false,
	}
	err = basketRepo.DeleteFromBasket(ctx, dishToBasket, 1)
	if err == nil {
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
	basketRepo := &basketRepository{
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
	err = basketRepo.DeleteFromBasket(ctx, dishToBasket, 1)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestOrderRepo_DecDishCountInBasket_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	basketRepo := &basketRepository{
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
		WillReturnError(sql.ErrNoRows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	dishToBasket := models.DishToBasket{
		DishID:     1,
		SameBasket: true,
		IsPlus:     false,
	}
	err = basketRepo.DeleteFromBasket(ctx, dishToBasket, 1)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestOrderRepo_GetBasket(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	basketRepo := &basketRepository{
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

	basketResult, errr := basketRepo.GetBasket(ctx, 1)
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

func TestOrderRepo_GetBasket_SelectBidError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	basketRepo := &basketRepository{
		DB: db,
	}

	mock.
		ExpectQuery("select basketID from basket_users").
		WithArgs(1).
		WillReturnError(sql.ErrConnDone)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	_, errr := basketRepo.GetBasket(ctx, 1)
	if errr == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestOrderRepo_GetBasket_GetRestaurantError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	basketRepo := &basketRepository{
		DB: db,
	}

	basketID := sqlmock.NewRows([]string{"dish"})
	expectBasketID := []models.BasketForUser{
		{BID: 1},
	}
	for _, item := range expectBasketID {
		basketID = basketID.AddRow(item.BID)
	}

	mock.
		ExpectQuery("select basketID from basket_users").
		WithArgs(1).
		WillReturnRows(basketID)
	mock.
		ExpectQuery("select restaurant from baskets").
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	_, errr := basketRepo.GetBasket(ctx, 1)
	if errr == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestOrderRepo_GetBasket_GetAvatarError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	basketRepo := &basketRepository{
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
		WillReturnError(sql.ErrNoRows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	_, errr := basketRepo.GetBasket(ctx, 1)
	if errr == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestOrderRepo_GetBasket_RestInfoError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	basketRepo := &basketRepository{
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
		WillReturnError(sql.ErrNoRows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	_, errr := basketRepo.GetBasket(ctx, 1)
	if errr == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestOrderRepo_GetBasket_GetDishesErr(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	basketRepo := &basketRepository{
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
		WillReturnError(sql.ErrNoRows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	_, errr := basketRepo.GetBasket(ctx, 1)
	if errr == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestBasketRepository_AddBasket(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	basketRepo := &basketRepository{
		DB: db,
	}

	restaurantName := sqlmock.NewRows([]string{"restaurant"})
	restaurantName.AddRow("rest1")

	basketID := sqlmock.NewRows([]string{"bid"})
	basketID.AddRow(1)

	mock.
		ExpectQuery("select name from restaurants").
		WithArgs(1).
		WillReturnRows(restaurantName)
	mock.
		ExpectQuery("insert into baskets").
		WithArgs("rest1").
		WillReturnRows(basketID)
	mock.
		ExpectExec("insert into basket_users").
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	bid, errr := basketRepo.AddBasket(ctx, 1, 1)
	if errr != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	require.EqualValues(t, 1, bid)
}

func TestBasketRepository_AddBasket_ErrorInsertBasket(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	basketRepo := &basketRepository{
		DB: db,
	}

	restaurantName := sqlmock.NewRows([]string{"restaurant"})
	restaurantName.AddRow("rest1")

	basketID := sqlmock.NewRows([]string{"bid"})
	basketID.AddRow(1)

	mock.
		ExpectQuery("select name from restaurants").
		WithArgs(1).
		WillReturnRows(restaurantName)
	mock.
		ExpectQuery("insert into baskets").
		WithArgs("rest1").
		WillReturnError(sql.ErrNoRows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	_, errr := basketRepo.AddBasket(ctx, 1, 1)
	if errr == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestBasketRepository_AddBasket_InsertBasketError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	basketRepo := &basketRepository{
		DB: db,
	}

	restaurantName := sqlmock.NewRows([]string{"restaurant"})
	restaurantName.AddRow("rest1")

	basketID := sqlmock.NewRows([]string{"bid"})
	basketID.AddRow(1)

	mock.
		ExpectQuery("select name from restaurants").
		WithArgs(1).
		WillReturnRows(restaurantName)
	mock.
		ExpectQuery("insert into baskets").
		WithArgs("rest1").
		WillReturnRows(basketID)
	mock.
		ExpectExec("insert into basket_users").
		WithArgs(1, 1).
		WillReturnError(sql.ErrNoRows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	_, errr := basketRepo.AddBasket(ctx, 1, 1)
	if errr == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestBasketRepository_AddDishToBasket(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	basketRepo := &basketRepository{
		DB: db,
	}

	dish := models.DishInBasket{
		ID:     1,
		Name:   "dish1",
		Number: 1,
	}

	dishID := sqlmock.NewRows([]string{"bid"})
	dishID.AddRow(1)

	mock.
		ExpectQuery("select dish from baskets_food").
		WithArgs(1, 1).
		WillReturnRows(dishID)

	mock.
		ExpectExec("update baskets_food set").
		WithArgs(1, 1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	err = basketRepo.AddDishToBasket(ctx, 1, dish)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestBasketRepository_AddDishToBasket_updBasketError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	basketRepo := &basketRepository{
		DB: db,
	}

	dish := models.DishInBasket{
		ID:     1,
		Name:   "dish1",
		Number: 1,
	}

	dishID := sqlmock.NewRows([]string{"bid"})
	dishID.AddRow(1)

	mock.
		ExpectQuery("select dish from baskets_food").
		WithArgs(1, 1).
		WillReturnRows(dishID)

	mock.
		ExpectExec("update baskets_food set").
		WithArgs(1, 1, 1).
		WillReturnError(sql.ErrNoRows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	err = basketRepo.AddDishToBasket(ctx, 1, dish)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestBasketRepository_AddDishToBasket_NewDish(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	basketRepo := &basketRepository{
		DB: db,
	}

	dish := models.DishInBasket{
		ID:     1,
		Name:   "dish1",
		Number: 1,
	}

	dishID := sqlmock.NewRows([]string{"bid"})
	dishID.AddRow(1)

	mock.
		ExpectQuery("select dish from baskets_food").
		WithArgs(1, 1).
		WillReturnError(sql.ErrNoRows)

	mock.
		ExpectExec("insert into baskets_food").
		WithArgs(1, 1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	err = basketRepo.AddDishToBasket(ctx, 1, dish)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestBasketRepository_AddDishToBasket_NewDish_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	basketRepo := &basketRepository{
		DB: db,
	}

	dish := models.DishInBasket{
		ID:     1,
		Name:   "dish1",
		Number: 1,
	}

	dishID := sqlmock.NewRows([]string{"bid"})
	dishID.AddRow(1)

	mock.
		ExpectQuery("select dish from baskets_food").
		WithArgs(1, 1).
		WillReturnError(sql.ErrNoRows)

	mock.
		ExpectExec("insert into baskets_food").
		WithArgs(1, 1, 1).
		WillReturnError(sql.ErrNoRows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	err = basketRepo.AddDishToBasket(ctx, 1, dish)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestBasketRepository_DeleteBasket(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	basketRepo := &basketRepository{
		DB: db,
	}

	mock.
		ExpectExec("delete from basket_users").
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.
		ExpectExec("delete from baskets").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	err = basketRepo.DeleteBasket(ctx, 1, 1)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestBasketRepository_DeleteBasket_ErrorDeleteLink(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	basketRepo := &basketRepository{
		DB: db,
	}

	mock.
		ExpectExec("delete from basket_users").
		WithArgs(1, 1).
		WillReturnError(sql.ErrNoRows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	err = basketRepo.DeleteBasket(ctx, 1, 1)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestBasketRepository_DeleteBasket_ErrorDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	basketRepo := &basketRepository{
		DB: db,
	}

	mock.
		ExpectExec("delete from basket_users").
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.
		ExpectExec("delete from baskets").
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	err = basketRepo.DeleteBasket(ctx, 1, 1)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestBasketRepository_GetAddress(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	basketRepo := &basketRepository{
		DB: db,
	}

	addressDB := sqlmock.NewRows([]string{"name", "latitude", "longitude", "radius"})
	addressDB.AddRow("address1", "1234", "4321", 1500)

	mock.
		ExpectQuery("SELECT name, latitude,").
		WithArgs(1).
		WillReturnRows(addressDB)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	address, errr := basketRepo.GetAddress(ctx, 1)
	if errr != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	require.EqualValues(t, address.Name, "address1")
}

func TestBasketRepository_GetAddress_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	basketRepo := &basketRepository{
		DB: db,
	}

	addressDB := sqlmock.NewRows([]string{"name", "latitude", "longitude", "radius"})
	addressDB.AddRow("address1", "1234", "4321", 1500)

	mock.
		ExpectQuery("SELECT name, latitude,").
		WithArgs(1).
		WillReturnError(sql.ErrConnDone)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	_, err = basketRepo.GetAddress(ctx, 1)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}
