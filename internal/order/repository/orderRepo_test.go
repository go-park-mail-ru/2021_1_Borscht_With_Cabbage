package repository

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/borscht/backend/internal/models"
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
		restaurantName = restaurantName.AddRow(item.BID)
	}

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
		WithArgs(1).
		WillReturnRows(basketID)

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
	//db, mock, err := sqlmock.New()
	//if err != nil {
	//	t.Fatalf("cant create mock: %s", err)
	//}
	//defer db.Close()
	//orderRepo := &orderRepo{
	//	DB: db,
	//}
	//
	//rows := sqlmock.NewRows([]string{"rid", "name", "deliveryCost", "avgCheck", "description", "rating", "avatar"})
	//expext := models.DishToBasket{
	//
	//}
	//for _, item := range expect {
	//	rows = rows.AddRow(item.ID, item.Title, item.DeliveryCost, item.AvgCheck, item.Description, item.Rating, item.Avatar)
	//}
	//
	//
	//mock.
	//	ExpectQuery("select basketID from").
	//	WithArgs(1).
	//	WillReturnError(nil)
	//
	//c := context.Background()
	//ctx := context.WithValue(c, "request_id", 1)
	//
	//dishToBasket := models.DishToBasket{
	//	DishID: 1,
	//	SameBasket: true,
	//	IsPlus: true,
	//}
	//err = orderRepo.AddToBasket(ctx, dishToBasket, 1)
	//if err != nil {
	//	t.Errorf("unexpected err: %s", err)
	//	return
	//}
	//
	//if err := mock.ExpectationsWereMet(); err != nil {
	//	t.Errorf("there were unfulfilled expectations: %s", err)
	//	return
	//}
}

func TestOrderRepo_DeleteFromBasket(t *testing.T) {

}

func TestOrderRepo_Create(t *testing.T) {

}

func TestOrderRepo_GetUserOrders(t *testing.T) {

}

func TestOrderRepo_GetRestaurantOrders(t *testing.T) {

}

func TestOrderRepo_GetBasket(t *testing.T) {

}
