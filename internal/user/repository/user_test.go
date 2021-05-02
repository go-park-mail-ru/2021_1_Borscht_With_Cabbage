package repository

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/utils/logger"
	"testing"
)

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
