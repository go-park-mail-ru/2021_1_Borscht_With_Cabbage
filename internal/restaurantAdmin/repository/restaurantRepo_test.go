package repository

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/utils/logger"
)

func TestNewRestaurantRepo(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	defer db.Close()
	restRepo := NewRestaurantRepo(db)
	if restRepo != nil {
		return
	}
}

func TestRestaurantRepo_GetAddress(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	restaurantRepo := &restaurantRepo{
		DB: db,
	}

	rowAddress := sqlmock.NewRows([]string{"name", "latitude", "longitude", "radius"})
	rowAddress.AddRow("Бауманская 2", 0, 0, 1500)

	restID := 1
	mock.
		ExpectQuery("SELECT").
		WithArgs(restID).
		WillReturnRows(rowAddress)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	logger.InitLogger()
	address, err := restaurantRepo.GetAddress(ctx, restID)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	require.EqualValues(t, address.Name, "Бауманская 2")
}

func TestRestaurantRepo_AddAddress(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	restaurantRepo := &restaurantRepo{
		DB: db,
	}

	address := models.Address{
		Name:      "Бауманская 2",
		Longitude: 0,
		Latitude:  0,
		Radius:    1500,
	}
	restID := 1
	mock.
		ExpectExec("INSERT").
		WithArgs(1, address.Name, address.Latitude, address.Longitude, address.Radius).
		WillReturnResult(sqlmock.NewResult(1, 1))

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	logger.InitLogger()

	err = restaurantRepo.AddAddress(ctx, restID, address)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestRestaurantRepo_AddAddress_SqlErr(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	restaurantRepo := &restaurantRepo{
		DB: db,
	}

	address := models.Address{}
	restID := 0
	mock.
		ExpectExec("INSERT").
		WithArgs(restID, address.Name, address.Latitude, address.Longitude, address.Radius).
		WillReturnError(sql.ErrNoRows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	logger.InitLogger()

	err = restaurantRepo.AddAddress(ctx, restID, address)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestRestaurantRepo_UpdateAddress(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	restaurantRepo := &restaurantRepo{
		DB: db,
	}

	address := models.Address{
		Name:      "Бауманская 2",
		Longitude: 0,
		Latitude:  0,
		Radius:    1500,
	}
	restID := 1
	mock.
		ExpectExec("UPDATE").
		WithArgs(address.Name, address.Latitude, address.Longitude, address.Radius, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	logger.InitLogger()

	err = restaurantRepo.UpdateAddress(ctx, restID, address)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestRestaurantRepo_UpdateAddress_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	restaurantRepo := &restaurantRepo{
		DB: db,
	}

	address := models.Address{}
	restID := 1
	mock.
		ExpectExec("UPDATE").
		WithArgs(address.Name, address.Latitude, address.Longitude, address.Radius, 1).
		WillReturnError(sql.ErrNoRows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	logger.InitLogger()

	err = restaurantRepo.UpdateAddress(ctx, restID, address)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestRestaurantRepo_UpdateRestaurantData(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	restaurantRepo := &restaurantRepo{
		DB: db,
	}

	mock.
		ExpectQuery("select rid from restaurants where adminemail").
		WithArgs("dasha@mail.ru").
		WillReturnError(sql.ErrNoRows)
	mock.
		ExpectQuery("select rid from restaurants where adminphone").
		WithArgs("89111111111").
		WillReturnError(sql.ErrNoRows)
	mock.
		ExpectQuery("select rid from restaurants where name").
		WithArgs("rest1").
		WillReturnError(sql.ErrNoRows)
	mock.
		ExpectExec("UPDATE restaurants SET name").
		WithArgs("rest1", "dasha@mail.ru", "89111111111", 200, "yum", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := new(context.Context)

	restaurant := models.RestaurantUpdateData{
		ID:            1,
		AdminPhone:    "89111111111",
		AdminEmail:    "dasha@mail.ru",
		AdminPassword: "111111",
		DeliveryCost:  200,
		Description:   "yum",
		Title:         "rest1",
	}
	err = restaurantRepo.UpdateRestaurantData(*ctx, restaurant)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestRestaurantRepo_UpdateRestaurantImage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	restaurantRepo := &restaurantRepo{
		DB: db,
	}

	mock.
		ExpectExec("UPDATE restaurants SET").
		WithArgs("img.png", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	logger.InitLogger()
	err = restaurantRepo.UpdateRestaurantImage(ctx, 1, "img.png")
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestRestaurantRepo_UpdateRestaurantImageError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	restaurantRepo := &restaurantRepo{
		DB: db,
	}

	mock.
		ExpectExec("UPDATE restaurants SET").
		WithArgs("img.png", 1).
		WillReturnError(sql.ErrNoRows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	logger.InitLogger()
	err = restaurantRepo.UpdateRestaurantImage(ctx, 1, "img.png")
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestRestaurantRepo_DeleteAllCategories(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	restaurantRepo := &restaurantRepo{
		DB: db,
	}

	restID := 1
	mock.
		ExpectExec("DELETE").
		WithArgs(restID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	logger.InitLogger()

	err = restaurantRepo.DeleteAllCategories(ctx, restID)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestRestaurantRepo_DeleteAllCategories_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	restaurantRepo := &restaurantRepo{
		DB: db,
	}

	restID := 1
	mock.
		ExpectExec("DELETE").
		WithArgs(restID).
		WillReturnError(sql.ErrNoRows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	logger.InitLogger()

	err = restaurantRepo.DeleteAllCategories(ctx, restID)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestRestaurantRepo_AddCategories(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	restaurantRepo := &restaurantRepo{
		DB: db,
	}

	mock.
		ExpectQuery("INSERT INTO ").
		WithArgs("burgers", 1).
		WillReturnError(sql.ErrNoRows)

	ctx := new(context.Context)

	var categories []string
	categories = append(categories, "burgers")
	err = restaurantRepo.AddCategories(*ctx, 1, categories)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestRestaurantRepo_GetCategories(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	restaurantRepo := &restaurantRepo{
		DB: db,
	}

	RID := 1

	rows := sqlmock.NewRows([]string{"CATEGORYID"})
	rows.AddRow("food")

	mock.
		ExpectQuery("SELECT").
		WithArgs(RID).
		WillReturnRows(rows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	logger.InitLogger()
	_, err = restaurantRepo.GetCategories(ctx, RID)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestRestaurantRepo_GetCategories_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	restaurantRepo := &restaurantRepo{
		DB: db,
	}

	RID := 1

	mock.
		ExpectQuery("SELECT").
		WithArgs(RID).
		WillReturnError(sql.ErrNoRows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	logger.InitLogger()
	_, err = restaurantRepo.GetCategories(ctx, RID)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}
