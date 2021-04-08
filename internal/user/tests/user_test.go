package tests

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/user/repository"
	"github.com/stretchr/testify/require"
	"testing"
)

type UserItem struct {
	Uid      int    `json:"uid"`
	Name     string `json:"name"`
	Phone    string `json:"number"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Avatar   string
}

type uidStruct struct {
	Uid int `json:"uid"`
}

func TestUserCreate(t *testing.T) {
	/* ----- positive test ----- */
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	userRepo := &repository.UserRepo{
		DB: db,
	}

	rows := sqlmock.NewRows([]string{"uid"})
	expect := []*uidStruct{
		{1},
	}
	for _, item := range expect {
		rows = rows.AddRow(item.Uid)
	}

	mock.
		ExpectQuery("select uid from users where email =").
		WithArgs("kate@mail.ru").
		WillReturnError(sql.ErrNoRows)

	mock.
		ExpectQuery("select uid from users where phone =").
		WithArgs("81111111111").
		WillReturnError(sql.ErrNoRows)

	mock.
		ExpectQuery("insert into users").
		WithArgs("Kate", "81111111111", "kate@mail.ru", "111111", "http://127.0.0.1:5000/default/avatar/stas.jpg").
		WillReturnRows(rows)

	ctx := new(context.Context)

	user := models.User{
		Email:    "kate@mail.ru",
		Name:     "Kate",
		Phone:    "81111111111",
		Password: "111111",
		Avatar:   "",
	}
	uid, err := userRepo.Create(*ctx, user)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	require.EqualValues(t, uid, 1)
}

func TestCheckUserExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	userRepo := &repository.UserRepo{
		DB: db,
	}

	rows := sqlmock.NewRows([]string{"uid", "name", "phone", "email", "password", "photo"})
	expect := []*UserItem{
		{1, "Kate", "81111111111", "kate@mail.ru", "111111", "http://127.0.0.1:5000/default/avatar/stas.jpg"},
	}
	for _, item := range expect {
		rows = rows.AddRow(item.Uid, item.Name, item.Phone, item.Email, item.Password, item.Avatar)
	}

	mock.
		ExpectQuery("select uid, name, phone, email, password, photo").
		WithArgs("kate@mail.ru", "111111").
		WillReturnRows(rows)

	ctx := new(context.Context)
	user := models.UserAuth{
		Login:    "kate@mail.ru",
		Password: "111111",
	}

	foundUser := new(models.User)
	foundUser, err = userRepo.CheckUserExists(*ctx, user)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	require.EqualValues(t, foundUser.Uid, 1)
	require.EqualValues(t, foundUser.Name, "Kate")
	require.EqualValues(t, foundUser.Email, "kate@mail.ru")
	require.EqualValues(t, foundUser.Phone, "81111111111")
	require.EqualValues(t, foundUser.Password, "111111")
}

func TestGetByUid(t *testing.T) {

}

func TestEditProfile(t *testing.T) {

}
