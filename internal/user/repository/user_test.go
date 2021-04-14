package repository

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

type UserItem struct {
	Uid    int    `json:"uid"`
	Name   string `json:"name"`
	Phone  string `json:"number"`
	Email  string `json:"email"`
	Avatar string
}

type UserInfo struct {
	Name   string `json:"name"`
	Phone  string `json:"number"`
	Email  string `json:"email"`
	Avatar string
}

type uidStruct struct {
	Uid int `json:"uid"`
}

func TestUserCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	userRepo := &UserRepo{
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
		WithArgs("Kate", "81111111111", "kate@mail.ru", "111111", config.DefaultUserImage).
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

func TestUserCreateNegative(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	userRepo := &UserRepo{
		DB: db,
	}

	//rows := sqlmock.NewRows([]string{"uid"})
	//rows = rows.AddRow(1)

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
		WillReturnRows(rows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	user := models.User{
		Email:    "kate@mail.ru",
		Name:     "Kate",
		Phone:    "81111111111",
		Password: "111111",
		Avatar:   "",
	}

	logger.InitLogger()
	uid, err := userRepo.Create(ctx, user)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	require.EqualValues(t, uid, 0)
}

func TestCheckUserExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	userRepo := &UserRepo{
		DB: db,
	}

	rows := sqlmock.NewRows([]string{"uid", "name", "phone", "email", "photo"})
	expect := []*UserItem{
		{1, "Kate", "81111111111", "kate@mail.ru", "http://127.0.0.1:5000/default/avatar/stas.jpg"},
	}
	for _, item := range expect {
		rows = rows.AddRow(item.Uid, item.Name, item.Phone, item.Email, item.Avatar)
	}

	mock.
		ExpectQuery("select uid, name").
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
}

func TestCheckUserExistsNegative(t *testing.T) {
	// TODO
}

func TestGetByUid(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	userRepo := &UserRepo{
		DB: db,
	}

	rows := sqlmock.NewRows([]string{"name", "phone", "email", "photo"})
	expect := []*UserInfo{
		{"Kate", "81111111111", "kate@mail.ru", "http://127.0.0.1:5000/default/avatar/stas.jpg"},
	}
	for _, item := range expect {
		rows = rows.AddRow(item.Name, item.Phone, item.Email, item.Avatar)
	}

	mock.
		ExpectQuery("select name, phone, email, photo from users where uid=").
		WithArgs(1).
		WillReturnRows(rows)

	ctx := new(context.Context)
	uid := 1

	foundUser := new(models.User)
	*foundUser, err = userRepo.GetByUid(*ctx, uid)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestGetByUidNegative(t *testing.T) {
	// TODO
}

func TestEditProfile(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	userRepo := &UserRepo{
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

func TestEditProfileNegative(t *testing.T) {
	// TODO
}

func TestUploadAvatar(t *testing.T) {
	// TODO
}