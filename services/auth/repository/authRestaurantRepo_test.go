package authRepo

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/utils/logger"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewRestaurantAuthRepo(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	defer db.Close()
	restRepo := NewRestaurantAuthRepo(db)
	if restRepo != nil {
		return
	}
}

func TestExistingRestaurant_Email(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	restaurantRepo := &authRestaurantRepo{
		DB: db,
	}

	ridInDB := sqlmock.NewRows([]string{"rid"})
	ridInDB.AddRow(2)

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

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)
	logger.InitLogger()

	err = restaurantRepo.checkExistingRestaurant(ctx, restaurant)
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
	restaurantRepo := &authRestaurantRepo{
		DB: db,
	}

	ridInDB := sqlmock.NewRows([]string{"rid"})
	ridInDB.AddRow(2)

	restaurant := models.CheckRestaurantExists{
		CurrentRestId: 1,
		Email:         "dasha@mail.ru",
		Name:          "rest1",
		Number:        "89111111111",
	}

	mock.
		ExpectQuery("select rid from restaurants where adminemail").
		WithArgs("dasha@mail.ru").
		WillReturnError(sql.ErrNoRows)
	mock.
		ExpectQuery("select rid from restaurants where adminphone").
		WithArgs("89111111111").
		WillReturnError(nil).WillReturnRows(ridInDB)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)
	logger.InitLogger()

	err = restaurantRepo.checkExistingRestaurant(ctx, restaurant)
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
	restaurantRepo := &authRestaurantRepo{
		DB: db,
	}

	ridInDB := sqlmock.NewRows([]string{"rid"})
	ridInDB.AddRow(2)

	restaurant := models.CheckRestaurantExists{
		CurrentRestId: 1,
		Email:         "dasha@mail.ru",
		Name:          "rest1",
		Number:        "89111111111",
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
		WillReturnError(nil).WillReturnRows(ridInDB)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)
	logger.InitLogger()

	err = restaurantRepo.checkExistingRestaurant(ctx, restaurant)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}
func TestRestaurantRepo_CreateRestaurant(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	restaurantRepo := &authRestaurantRepo{
		DB: db,
	}

	rid := sqlmock.NewRows([]string{"rid"})
	rid.AddRow(1)

	mock.
		ExpectQuery("insert into restaurants").
		WillReturnRows(rid)
	mock.ExpectExec("insert into categories").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	//mock.
	//	ExpectQuery("select rid from restaurants where adminemail").
	//	WithArgs("dasha@mail.ru").
	//	WillReturnError(sql.ErrNoRows)
	//mock.
	//	ExpectQuery("select rid from restaurants where adminphone").
	//	WithArgs("89111111111").
	//	WillReturnError(sql.ErrNoRows)
	//mock.
	//	ExpectQuery("select rid from restaurants where name").
	//	WithArgs("rest1").
	//	WillReturnError(sql.ErrNoRows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	logger.InitLogger()
	restaurant := models.RestaurantInfo{
		AdminPhone:    "89111111111",
		AdminEmail:    "dasha@mail.ru",
		AdminPassword: "111111",
		DeliveryCost:  200,
		Description:   "yum",
		Title:         "rest1",
	}
	var id int
	id, err = restaurantRepo.CreateRestaurant(ctx, restaurant)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	require.EqualValues(t, id, 1)
}

func TestRestaurantRepo_CheckRestaurantExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	restaurantRepo := &authRestaurantRepo{
		DB: db,
	}

	restaurant := models.RestaurantAuth{
		Login:    "dasha@mail.ru",
		Password: "111111",
	}
	restaurantFromBD := sqlmock.NewRows([]string{"rid", "name", "adminemail", "adminphone", "deliveryCost", "avgCheck",
		"description", "avatar", "adminpassword"})
	restaurantFromBD.AddRow(1, "rest1", "dasha@mail.ru", "89111111111", 200, 1200, "yum", config.DefaultRestaurantImage, "111111")

	mock.
		ExpectQuery("select rid, name, adminemail,").
		WithArgs(restaurant.Login).
		WillReturnRows(restaurantFromBD)

	ctx := new(context.Context)

	restaurantResponse := new(models.RestaurantInfo)
	restaurantResponse, err = restaurantRepo.GetByLogin(*ctx, restaurant.Login)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	require.EqualValues(t, restaurantResponse.ID, 1)
}

func TestRestaurantRepo_GetByRid(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	restaurantRepo := &authRestaurantRepo{
		DB: db,
	}

	restaurantFromBD := sqlmock.NewRows([]string{"title", "adminphone", "adminemail", "avatar"})
	restaurantFromBD.AddRow("rest1", "89111111111", "dasha@mail.ru", config.DefaultRestaurantImage)

	mock.
		ExpectQuery("select name, adminphone,").
		WithArgs(1).
		WillReturnRows(restaurantFromBD)

	ctx := new(context.Context)

	restaurantResponse := new(models.RestaurantInfo)
	restaurantResponse, err = restaurantRepo.GetByRid(*ctx, 1)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	require.EqualValues(t, restaurantResponse.Title, "rest1")
}
