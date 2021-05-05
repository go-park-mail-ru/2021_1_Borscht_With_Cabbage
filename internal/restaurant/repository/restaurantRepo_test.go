package repository

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/borscht/backend/internal/models"
	restModel "github.com/borscht/backend/internal/restaurant"
	"github.com/stretchr/testify/require"
	"testing"
)

type RestaurantInfo struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	DeliveryCost int    `json:"deliveryCost"`
	AvgCheck     int    `json:"cost"`
	Description  string `json:"description"`
	Avatar       string `json:"avatar"`
	RatingsSum   float64
	ReviewsCount float64
	Latitude     string
	Longitude    string
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
	Section     int    `json:"section"`
	Image       string `json:"image"`
}

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

func TestRestaurantRepo_GetVendor(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	restaurantRepo := &restaurantRepo{
		DB: db,
	}

	rows := sqlmock.NewRows([]string{"rid", "name", "deliveryCost", "avgCheck", "description", "avatar", "ratingssum", "reviewscount", "latitude", "longitude"})
	expect := []*RestaurantInfo{
		{1, "Rest1", 200, 1200, "new", "img.jpg", 10, 2, "55.766516", "37.653424"},
		{2, "Rest2", 100, 1300, "new2", "img2.jpg", 8, 2, "55.735439", "37.584981"},
	}
	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.Title, item.DeliveryCost, item.AvgCheck, item.Description, item.Avatar, item.RatingsSum, item.ReviewsCount, item.Latitude, item.Longitude)
	}

	params := restModel.GetVendorParams{
		Limit:     1,
		Offset:    2,
		Address:   true,
		Latitude:  "55.768096",
		Longitude: "37.646839",
	}

	mock.
		ExpectQuery("SELECT r.rid,").
		WillReturnRows(rows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	restaurants := make([]models.RestaurantInfo, 0)
	restaurants, err = restaurantRepo.GetVendor(ctx, params)
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

	restaurant := sqlmock.NewRows([]string{"rid", "name", "deliveryCost", "avgCheck", "description", "avatar", "ratingssum", "reviewscount", "lan", "lon"})
	expectRestaurant := []*RestaurantInfo{
		{1, "Rest1", 200, 1200, "new", "img.jpg", 10, 5, "55.766516", "37.653424"},
	}
	for _, item := range expectRestaurant {
		restaurant = restaurant.AddRow(item.ID, item.Title, item.DeliveryCost, item.AvgCheck, item.Description, item.Avatar, item.RatingsSum, item.ReviewsCount, item.Latitude, item.Longitude)
	}

	dishes := sqlmock.NewRows([]string{"did", "name", "price", "weight", "description", "section", "image"})
	expectDishes := []*DishInfo{
		{1, "Dish1", 200, 240, "new", 1, "img.jpg"},
		{2, "Dish2", 100, 130, "new2", 2, "img2.jpg"},
	}
	for _, item := range expectDishes {
		dishes = dishes.AddRow(item.did, item.Name, item.Price, item.Weight, item.Description, item.Section, item.Image)
	}

	sections := sqlmock.NewRows([]string{"sid", "name"})
	expectSections := []*models.Section{
		{1, "section1", 0},
		{2, "section2", 0},
	}
	for _, item := range expectSections {
		sections = sections.AddRow(item.ID, item.Name)
	}

	mock.
		ExpectQuery("SELECT r.rid, r.name, deliveryCost, ").
		WithArgs(1).
		WillReturnRows(restaurant)

	mock.
		ExpectQuery("select sid, name from sections").
		WithArgs(1).
		WillReturnRows(sections)

	mock.
		ExpectQuery("select did, name,").
		WithArgs(1).
		WillReturnRows(dishes)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	restaurants := new(models.RestaurantWithDishes)
	restaurants, err = restaurantRepo.GetById(ctx, 1, models.Coordinates{})
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
