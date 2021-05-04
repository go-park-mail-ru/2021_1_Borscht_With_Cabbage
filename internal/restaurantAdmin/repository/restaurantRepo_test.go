package repository

import (
	"context"
	"database/sql"
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
