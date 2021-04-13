package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/models"
	"github.com/stretchr/testify/require"
	"testing"
)

type RestaurantInfo struct {
	ID           int     `json:"id"`
	Title        string  `json:"title"`
	DeliveryCost int     `json:"deliveryCost"`
	AvgCheck     int     `json:"cost"`
	Description  string  `json:"description"`
	Rating       float64 `json:"rating"`
	Avatar       string  `json:"avatar"`
}

type Restaurant struct {
	Title        string  `json:"title"`
	DeliveryCost int     `json:"deliveryCost"`
	AvgCheck     int     `json:"cost"`
	Description  string  `json:"description"`
	Rating       float64 `json:"rating"`
	Avatar       string  `json:"avatar"`
}

type DishInfo struct {
	did         int    `json:"did"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Weight      int    `json:"weight"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

func TestNewRestaurantRepo(t *testing.T) {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s", config.DBUser, config.DBPass, config.DBName)
	db, err := sql.Open(config.PostgresDB, dsn)
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	NewRestaurantRepo(db)
}

func TestRestaurantRepo_GetVendor(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	restaurantRepo := &restaurantRepo{
		DB: db,
	}

	rows := sqlmock.NewRows([]string{"rid", "name", "deliveryCost", "avgCheck", "description", "rating", "avatar"})
	expect := []*RestaurantInfo{
		{1, "Rest1", 200, 1200, "new", 5, "img.jpg"},
		{2, "Rest2", 100, 1300, "new2", 5, "img2.jpg"},
	}
	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.Title, item.DeliveryCost, item.AvgCheck, item.Description, item.Rating, item.Avatar)
	}

	mock.
		ExpectQuery("select rid").
		WithArgs(1, 3).
		WillReturnRows(rows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	restaurants := make([]models.RestaurantInfo, 0)
	restaurants, err = restaurantRepo.GetVendor(ctx, 2, 1)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	require.EqualValues(t, restaurants[0].ID, 1)
	require.EqualValues(t, restaurants[1].ID, 2)
}

func TestRestaurantRepo_GetById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	restaurantRepo := &restaurantRepo{
		DB: db,
	}

	restaurant := sqlmock.NewRows([]string{"name", "deliveryCost", "avgCheck", "description", "rating", "avatar"})
	expectRestaurant := []*Restaurant{
		{"Rest1", 200, 1200, "new", 5, "img.jpg"},
	}
	for _, item := range expectRestaurant {
		restaurant = restaurant.AddRow(item.Title, item.DeliveryCost, item.AvgCheck, item.Description, item.Rating, item.Avatar)
	}

	dishes := sqlmock.NewRows([]string{"did", "name", "price", "weight", "description", "image"})
	expectDishes := []*DishInfo{
		{1, "Dish1", 200, 240, "new", "img.jpg"},
		{2, "Dish2", 100, 130, "new2", "img2.jpg"},
	}
	for _, item := range expectDishes {
		dishes = dishes.AddRow(item.did, item.Name, item.Price, item.Weight, item.Description, item.Image)
	}

	mock.
		ExpectQuery("select name, deliveryCost").
		WithArgs(1).
		WillReturnRows(restaurant)

	mock.
		ExpectQuery("select did, name,").
		WithArgs(1).
		WillReturnRows(dishes)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	restaurants := new(models.RestaurantWithDishes)
	*restaurants, err = restaurantRepo.GetById(ctx, 1)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	require.EqualValues(t, restaurants.Dishes[0].Name, "Dish1")
	require.EqualValues(t, restaurants.Dishes[1].Name, "Dish2")
}