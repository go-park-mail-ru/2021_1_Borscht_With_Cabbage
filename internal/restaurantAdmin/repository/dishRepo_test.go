package repository

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/utils/logger"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewDishRepo(t *testing.T) {

}

func TestDishRepo_GetAllDishes(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	dishRepo := &dishRepo{
		DB: db,
	}

	rowsDishes := sqlmock.NewRows([]string{"did", "section", "name", "price", "weight", "description", "image"})
	expectDishes := []*models.Dish{
		{ID: 1, Section: 1, Name: "dish1", Price: 350, Weight: 150, Description: "yum", Image: "img.sql"},
	}
	for _, item := range expectDishes {
		rowsDishes = rowsDishes.AddRow(item.ID, item.Section, item.Name, item.Price, item.Weight, item.Description, item.Image)
	}

	mock.
		ExpectQuery("select did, section, name").
		WithArgs(1).
		WillReturnRows(rowsDishes)

	ctx := new(context.Context)

	var dishes []models.Dish
	dishes, err = dishRepo.GetAllDishes(*ctx, 1)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	require.EqualValues(t, dishes[0].ID, 1)
}

func TestDishRepo_UpdateDishData(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	dishRepo := &dishRepo{
		DB: db,
	}

	rowsDishes := sqlmock.NewRows([]string{"did", "name"})
	expectDishes := []*models.Dish{
		{ID: 1, Name: "dish1"},
	}
	for _, item := range expectDishes {
		rowsDishes = rowsDishes.AddRow(item.ID, item.Name)
	}

	dish := models.Dish{
		Name:        "newName",
		Restaurant:  1,
		ID:          1,
		Description: "new description",
		Price:       120,
		Weight:      100,
	}

	mock.
		ExpectQuery("select did, name from dishes").
		WithArgs(1).
		WillReturnRows(rowsDishes)
	mock.
		ExpectExec("update dishes set name").
		WithArgs("newName", 120, 100, "new description", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := new(context.Context)

	err = dishRepo.UpdateDishData(*ctx, dish)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestDishRepo_GetDish(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	dishRepo := &dishRepo{
		DB: db,
	}

	rowsDishes := sqlmock.NewRows([]string{"did", "section", "name", "price", "weight", "description", "image"})
	expectDishes := []*models.Dish{
		{ID: 1, Section: 1, Name: "dish1", Price: 350, Weight: 150, Description: "yum", Image: "img.sql"},
	}
	for _, item := range expectDishes {
		rowsDishes = rowsDishes.AddRow(item.ID, item.Section, item.Name, item.Price, item.Weight, item.Description, item.Image)
	}

	mock.
		ExpectQuery("select did, restaurantId, name,").
		WithArgs(1).
		WillReturnRows(rowsDishes)

	ctx := new(context.Context)

	dish := new(models.Dish)
	dish, err = dishRepo.GetDish(*ctx, 1)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	require.EqualValues(t, dish.ID, 1)
}

func TestDishRepo_DeleteDish(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	dishRepo := &dishRepo{
		DB: db,
	}

	mock.
		ExpectExec("delete from dishes").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := new(context.Context)

	err = dishRepo.DeleteDish(*ctx, 1)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestDishRepo_DeleteDishError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	dishRepo := &dishRepo{
		DB: db,
	}

	mock.
		ExpectExec("delete from restaurants").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	logger.InitLogger()
	err = dishRepo.DeleteDish(ctx, 1)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
}

func TestDishRepo_checkExistingDish(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	dishRepo := &dishRepo{
		DB: db,
	}

	rowsDishes := sqlmock.NewRows([]string{"did", "section", "name", "price", "weight", "description", "image"})
	expectDishes := []*models.Dish{
		{ID: 1, Section: 1, Name: "dish1", Price: 350, Weight: 150, Description: "yum", Image: "img.sql"},
	}
	for _, item := range expectDishes {
		rowsDishes = rowsDishes.AddRow(item.ID, item.Section, item.Name, item.Price, item.Weight, item.Description, item.Image)
	}
	dishData := models.CheckDishExists{RestaurantId: 1, Name: "name"}

	mock.
		ExpectQuery("select did, name from dishes").
		WithArgs(1).
		WillReturnRows(rowsDishes)

	ctx := new(context.Context)

	err = dishRepo.checkExistingDish(*ctx, dishData)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestDishRepo_AddDish(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	dishRepo := &dishRepo{
		DB: db,
	}

	rowsDishes := sqlmock.NewRows([]string{"did", "name"})
	expectDishes := []*models.Dish{
		{ID: 1, Name: "dish1"},
	}
	for _, item := range expectDishes {
		rowsDishes = rowsDishes.AddRow(item.ID, item.Name)
	}

	rowsDid := sqlmock.NewRows([]string{"did"})
	expectDid := []*models.Dish{
		{ID: 1},
	}
	for _, item := range expectDid {
		rowsDid = rowsDid.AddRow(item.ID)
	}

	restName := sqlmock.NewRows([]string{"name"})
	expectRestName := []*models.RestaurantInfo{
		{Title: "rest1"},
	}
	for _, item := range expectRestName {
		restName = restName.AddRow(item.Title)
	}

	dish := models.Dish{
		Name:        "newName",
		Restaurant:  1,
		Section:     1,
		ID:          1,
		Description: "new description",
		Price:       120,
		Weight:      100,
	}

	mock.
		ExpectQuery("select did, name from dishes").
		WithArgs(1).
		WillReturnRows(rowsDishes)
	mock.
		ExpectQuery("insert into dishes").
		WithArgs(1, 1, "newName", 120, 100, "new description", config.DefaultDishImage).
		WillReturnRows(rowsDid)
	mock.
		ExpectQuery("select name from restaurants").
		WithArgs(1).
		WillReturnRows(restName)
	mock.
		ExpectExec("update dishes set").
		WithArgs("rest1", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := new(context.Context)
	var id int
	id, err = dishRepo.AddDish(*ctx, dish)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	require.EqualValues(t, id, 1)
}

func TestDishRepo_AddDish_InsertError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	dishRepo := &dishRepo{
		DB: db,
	}

	rowsDishes := sqlmock.NewRows([]string{"did", "name"})
	expectDishes := []*models.Dish{
		{ID: 1, Name: "dish1"},
	}
	for _, item := range expectDishes {
		rowsDishes = rowsDishes.AddRow(item.ID, item.Name)
	}

	rowsDid := sqlmock.NewRows([]string{"did"})
	expectDid := []*models.Dish{
		{ID: 1},
	}
	for _, item := range expectDid {
		rowsDid = rowsDid.AddRow(item.ID)
	}

	dish := models.Dish{
		Name:        "newName",
		Restaurant:  1,
		Section:     1,
		ID:          1,
		Description: "new description",
		Price:       120,
		Weight:      100,
	}

	mock.
		ExpectQuery("select did, name from dishes").
		WithArgs(1).
		WillReturnRows(rowsDishes)
	mock.
		ExpectQuery("insert into dishes").
		WithArgs(1, 1, "newName", 120, 100, "new description").
		WillReturnRows(rowsDid)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	logger.InitLogger()
	_, err = dishRepo.AddDish(ctx, dish)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
}

func TestDishRepo_UpdateDishImage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	dishRepo := &dishRepo{
		DB: db,
	}

	mock.
		ExpectExec("UPDATE dishes SET image").
		WithArgs("img.jpg", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := new(context.Context)
	err = dishRepo.UpdateDishImage(*ctx, 1, "img.jpg")
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
}

func TestDishRepo_UpdateDishImageError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	dishRepo := &dishRepo{
		DB: db,
	}

	mock.
		ExpectExec("UPDATE dishes SET image").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	logger.InitLogger()
	err = dishRepo.UpdateDishImage(ctx, 1, "img.jpg")
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
}
