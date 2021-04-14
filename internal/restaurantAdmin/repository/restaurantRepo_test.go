package repository

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/borscht/backend/internal/models"
	"testing"
)

func TestNewRestaurantRepo(t *testing.T) {

}

func TestExistingRestaurant_Email(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	restaurantRepo := &restaurantRepo{
		DB: db,
	}

	ridInDB := sqlmock.NewRows([]string{"rid"})
	expectItem := []*models.RestaurantInfo{
		{ID: 1},
	}
	for _, item := range expectItem {
		ridInDB = ridInDB.AddRow(item.ID)
	}

	restaurant := models.CheckRestaurantExists{
		CurrentRestId: 1,
		Email:         "dasha@mail.ru",
		Name:          "rest1",
		Number:        "89111111111",
	}
	mock.
		ExpectQuery("select rid from restaurants where adminemail").
		WithArgs("dasha@mail.ru").
		WillReturnError(nil).WillReturnRows(ridInDB)

	ctx := new(context.Context)

	err = restaurantRepo.checkExistingRestaurant(*ctx, restaurant)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestExistingRestaurant_Phone(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	restaurantRepo := &restaurantRepo{
		DB: db,
	}

	ridInDB := sqlmock.NewRows([]string{"rid"})
	expectItem := []*models.RestaurantInfo{
		{ID: 1},
	}
	for _, item := range expectItem {
		ridInDB = ridInDB.AddRow(item.ID)
	}

	restaurant := models.CheckRestaurantExists{
		CurrentRestId: 1,
		Email:         "dasha@mail.ru",
		Name:          "rest1",
		Number:        "89111111111",
	}
	mock.
		ExpectQuery("select rid from restaurants where adminphone").
		WithArgs("89111111111").
		WillReturnError(nil).WillReturnRows(ridInDB)

	ctx := new(context.Context)

	err = restaurantRepo.checkExistingRestaurant(*ctx, restaurant)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestExistingRestaurant_Number(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	restaurantRepo := &restaurantRepo{
		DB: db,
	}

	ridInDB := sqlmock.NewRows([]string{"rid"})
	expectItem := []*models.RestaurantInfo{
		{ID: 1},
	}
	for _, item := range expectItem {
		ridInDB = ridInDB.AddRow(item.ID)
	}

	restaurant := models.CheckRestaurantExists{
		CurrentRestId: 1,
		Email:         "dasha@mail.ru",
		Name:          "rest1",
		Number:        "89111111111",
	}
	mock.
		ExpectQuery("select rid from restaurants where name").
		WithArgs("rest1").
		WillReturnError(nil).WillReturnRows(ridInDB)

	ctx := new(context.Context)

	err = restaurantRepo.checkExistingRestaurant(*ctx, restaurant)
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
		ExpectExec("update restaurants set name").
		WithArgs("rest1", "dasha@mail.ru", "89111111111", "111111", 200, "yum", 1).
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

func TestRestaurantRepo_CreateRestaurant(t *testing.T) {

}

func TestRestaurantRepo_UpdateRestaurantImage(t *testing.T) {

}

func TestRestaurantRepo_CheckRestaurantExists(t *testing.T) {

}

func TestRestaurantRepo_GetByRid(t *testing.T) {

}
