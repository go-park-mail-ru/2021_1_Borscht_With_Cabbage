package authRepo

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/services/auth/config"
	"github.com/borscht/backend/utils/logger"
	"github.com/stretchr/testify/require"
)

type UserInfo struct {
	Name   string `json:"name"`
	Phone  string `json:"number"`
	Email  string `json:"email"`
	Avatar string
}

type uidStruct struct {
	Uid int `json:"uid"`
}

func TestNewUserAuthRepo(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	defer db.Close()
	userRepo := NewUserAuthRepo(db)
	if userRepo != nil {
		return
	}
}

func TestUserRepo_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	userRepo := &authRepo{
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
		WillReturnRows(rows)

	ctx := new(context.Context)

	user := models.User{
		Email:    "kate@mail.ru",
		Name:     "Kate",
		Phone:    "81111111111",
		Password: "111111",
		Avatar:   config.ConfigStatic.DefaultUserImage,
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

func TestUserRepo_Create_UserExist(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	userRepo := &authRepo{
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

func TestUserRepo_Create_InsertError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	userRepo := &authRepo{
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
		WillReturnError(sql.ErrNoRows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)
	logger.InitLogger()

	user := models.User{
		Email:    "kate@mail.ru",
		Name:     "Kate",
		Phone:    "81111111111",
		Password: "111111",
		Avatar:   "",
	}

	_, err = userRepo.Create(ctx, user)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestUserRepo_GetByLogin(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	userRepo := &authRepo{
		DB: db,
	}

	rows := sqlmock.NewRows([]string{"uid", "name", "phone", "email", "photo", "password"})
	expect := []*models.User{
		{1, "Kate", "kate@mail.ru", "89111111111", []byte(""), "89111111111", "http://127.0.0.1:5000/default/avatar/stas.jpg", models.Address{}},
	}
	for _, item := range expect {
		rows = rows.AddRow(item.Uid, item.Name, item.Phone, item.Email, item.Avatar, item.Password)
	}

	mock.
		ExpectQuery("select uid, name").
		WithArgs("kate@mail.ru").
		WillReturnRows(rows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	logger.InitLogger()
	foundUser := new(models.User)
	foundUser, err = userRepo.GetByLogin(ctx, "kate@mail.ru")
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
	require.EqualValues(t, foundUser.Phone, "89111111111")
}

func TestUserRepo_GetByLogin_GetUserError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	userRepo := &authRepo{
		DB: db,
	}

	rows := sqlmock.NewRows([]string{"uid", "name", "phone", "email", "photo", "password"})
	expect := []*models.User{
		{1, "Kate", "kate@mail.ru", "89111111111", []byte(""), "89111111111", "http://127.0.0.1:5000/default/avatar/stas.jpg", models.Address{}},
	}
	for _, item := range expect {
		rows = rows.AddRow(item.Uid, item.Name, item.Phone, item.Email, item.Avatar, item.Password)
	}

	mock.
		ExpectQuery("select uid, name").
		WithArgs("kate@mail.ru").
		WillReturnError(sql.ErrNoRows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	logger.InitLogger()
	_, err = userRepo.GetByLogin(ctx, "kate@mail.ru")
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

}

func TestUserRepo_GetByUid(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	userRepo := &authRepo{
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

func TestAuthRestaurantRepo_GetAddress(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	userRepo := &authRepo{
		DB: db,
	}

	uid := 1

	rows := sqlmock.NewRows([]string{"name", "latitude", "longitude"})
	rows = rows.AddRow("proskekt 2", 23.3434, 34.443)

	mock.
		ExpectQuery("SELECT name,").
		WithArgs(1).
		WillReturnRows(rows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	logger.InitLogger()

	_, err = userRepo.GetAddress(ctx, uid)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestAuthRestaurantRepo_GetAddress_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	userRepo := &authRepo{
		DB: db,
	}

	uid := 1

	mock.
		ExpectQuery("SELECT name,").
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	logger.InitLogger()

	_, err = userRepo.GetAddress(ctx, uid)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestAuthRestaurantRepo_GetAddress_DBerror(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	userRepo := &authRepo{
		DB: db,
	}

	uid := 1

	mock.
		ExpectQuery("SELECT name,").
		WithArgs(1).
		WillReturnError(sql.ErrConnDone)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	logger.InitLogger()

	_, err = userRepo.GetAddress(ctx, uid)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}
