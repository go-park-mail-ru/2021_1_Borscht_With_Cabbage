package repository

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/utils/logger"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewUserRepo(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	defer db.Close()
	userRepo := NewUserRepo(db)
	if userRepo != nil {
		return
	}
}

func TestUserRepo_UpdateData(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	userRepo := &userRepo{
		DB: db,
	}

	mock.
		ExpectQuery("select uid from users where email =").
		WithArgs("kate@mail.ru").
		WillReturnError(sql.ErrNoRows)

	mock.
		ExpectQuery("SELECT uid FROM users WHERE phone =").
		WithArgs("81111111111").
		WillReturnError(sql.ErrNoRows)

	mock.
		ExpectExec("UPDATE users SET").
		WithArgs("81111111111", "kate@mail.ru", "Kate", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	userId := models.User{
		Uid:  1,
		Name: "Daria",
	}

	user := models.UserData{
		Phone:  "81111111111",
		Email:  "kate@mail.ru",
		Name:   "Kate",
		Avatar: "http://127.0.0.1:5000/default/avatar/stas.jpg",
	}
	user.ID = 1

	c := context.Background()
	ctx := context.WithValue(c, "User", userId)

	err = userRepo.UpdateData(ctx, user)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestUserRepo_UpdateAvatar(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	userRepo := &userRepo{
		DB: db,
	}

	mock.
		ExpectExec("UPDATE users SET").
		WithArgs("img.jpg", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := new(context.Context)

	err = userRepo.UpdateAvatar(*ctx, 1, "img.jpg")
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestUserRepo_UpdateAvatar_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	userRepo := &userRepo{
		DB: db,
	}

	mock.
		ExpectExec("UPDATE users SET").
		WithArgs("img.jpg", 1).
		WillReturnError(sql.ErrNoRows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)
	logger.InitLogger()

	err = userRepo.UpdateAvatar(ctx, 1, "img.jpg")
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestUserRepo_GetAddress(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	userRepo := &userRepo{
		DB: db,
	}

	address := sqlmock.NewRows([]string{"name", "longitude", "latitude"})
	address.AddRow("address1", 1234.0, 4321.0)

	mock.
		ExpectQuery("SELECT name, latitude").
		WithArgs(1).
		WillReturnRows(address)

	c := context.Background()

	addressResult, errr := userRepo.GetAddress(c, 1)
	if errr != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	require.EqualValues(t, addressResult.Name, "address1")
}

func TestUserRepo_GetAddress_GetError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	userRepo := &userRepo{
		DB: db,
	}

	address := sqlmock.NewRows([]string{"name", "longitude", "latitude"})
	address.AddRow("address1", "1234", "4321")

	mock.
		ExpectQuery("SELECT name, latitude").
		WithArgs(1).
		WillReturnError(sql.ErrConnDone)

	c := context.Background()

	_, errr := userRepo.GetAddress(c, 1)
	if errr == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestUserRepo_GetAddress_NoAddress(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	userRepo := &userRepo{
		DB: db,
	}

	mock.
		ExpectQuery("SELECT name, latitude").
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	c := context.Background()

	addressResult, _ := userRepo.GetAddress(c, 1)
	if addressResult == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

}

func TestUserRepo_AddAddress(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	userRepo := &userRepo{
		DB: db,
	}

	addressInput := models.Address{
		Name:      "address1",
		Longitude: 4321,
		Latitude:  1234,
	}

	mock.
		ExpectExec("INSERT INTO addresses").
		WithArgs(1, "address1", 1234.0, 4321.0).
		WillReturnResult(sqlmock.NewResult(1, 1))

	c := context.Background()

	err = userRepo.AddAddress(c, 1, addressInput)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestUserRepo_AddAddress_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	userRepo := &userRepo{
		DB: db,
	}

	addressInput := models.Address{}

	mock.
		ExpectExec("INSERT INTO addresses").
		WillReturnError(sql.ErrNoRows)

	c := context.Background()

	err = userRepo.AddAddress(c, 1, addressInput)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestUserRepo_DeleteAddress(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	userRepo := &userRepo{
		DB: db,
	}

	mock.
		ExpectExec("DELETE FROM addresses").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	c := context.Background()

	err = userRepo.DeleteAddress(c, 1)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestUserRepo_DeleteAddressError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	userRepo := &userRepo{
		DB: db,
	}

	mock.
		ExpectExec("DELETE FROM addresses").
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	c := context.Background()

	err = userRepo.DeleteAddress(c, 1)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}
