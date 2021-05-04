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

func TestNewChatRepo(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	chatRepoMock := NewChattRepo(db)
	if chatRepoMock != nil {
		return
	}
}

func TestChatRepo_GetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	chatRepoMock := NewChattRepo(db)

	rowsUser := sqlmock.NewRows([]string{"uid", "name", "photo"})
	rowsUser.AddRow(1, "Daria", config.DefaultUserImage)

	mock.ExpectQuery("SELECT").
		WithArgs(1).
		WillReturnRows(rowsUser)

	ctx := new(context.Context)
	infoOpponent, err := chatRepoMock.GetUser(*ctx, 1)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	require.EqualValues(t, infoOpponent.Name, "Daria")
}

func TestChatRepo_GetUser_UserNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	chatRepoMock := NewChattRepo(db)

	mock.ExpectQuery("SELECT").
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	ctx := new(context.Context)
	_, err = chatRepoMock.GetUser(*ctx, 1)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestChatRepo_GetUser_User_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	chatRepoMock := NewChattRepo(db)

	mock.ExpectQuery("SELECT").
		WithArgs(1).
		WillReturnError(sql.ErrConnDone)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)
	logger.InitLogger()

	_, err = chatRepoMock.GetUser(ctx, 1)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestChatRepo_GetRestaurant(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	chatRepoMock := NewChattRepo(db)

	rowsUser := sqlmock.NewRows([]string{"uid", "name", "photo"})
	rowsUser.AddRow(1, "Daria", config.DefaultUserImage)

	mock.ExpectQuery("SELECT").
		WithArgs(1).
		WillReturnRows(rowsUser)

	ctx := new(context.Context)
	infoOpponent, err := chatRepoMock.GetRestaurant(*ctx, 1)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	require.EqualValues(t, infoOpponent.Name, "Daria")
}

func TestChatRepo_SaveMessageFromUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	chatRepoMock := NewChattRepo(db)

	rowsUser := sqlmock.NewRows([]string{"mid"})
	rowsUser.AddRow(1)

	info := models.WsMessageForRepo{
		SentFromId: 1,
		SentToId:   1,
		Content:    "hello",
		Date:       "2.01.2020 21:11",
	}

	mock.ExpectQuery("INSERT").
		WithArgs(info.SentToId, info.SentToId, info.Content, info.Date).
		WillReturnRows(rowsUser)

	ctx := new(context.Context)
	mid, err := chatRepoMock.SaveMessageFromUser(*ctx, info)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	require.EqualValues(t, mid, 1)
}

func TestChatRepo_SaveMessageFromUser_sqlError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	chatRepoMock := NewChattRepo(db)

	info := models.WsMessageForRepo{
		Content: "hello",
		Date:    "2.01.2020 21:11",
	}

	mock.ExpectQuery("INSERT").
		WithArgs(0, 0, info.Content, info.Date).
		WillReturnError(sql.ErrConnDone)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)
	logger.InitLogger()

	_, err = chatRepoMock.SaveMessageFromUser(ctx, info)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestChatRepo_SaveMessageFromRestaurant(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	chatRepoMock := NewChattRepo(db)

	rowsUser := sqlmock.NewRows([]string{"mid"})
	rowsUser.AddRow(1)

	info := models.WsMessageForRepo{
		SentFromId: 1,
		SentToId:   1,
		Content:    "hello",
		Date:       "2.01.2020 21:11",
	}

	mock.ExpectQuery("INSERT").
		WithArgs(info.SentToId, info.SentToId, info.Content, info.Date).
		WillReturnRows(rowsUser)

	ctx := new(context.Context)
	mid, err := chatRepoMock.SaveMessageFromRestaurant(*ctx, info)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	require.EqualValues(t, mid, 1)
}
func TestChatRepo_SaveMessageFromRestaurant_sqlError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	chatRepoMock := NewChattRepo(db)

	info := models.WsMessageForRepo{
		SentFromId: 1,
		SentToId:   1,
		Content:    "hello",
		Date:       "2.01.2020 21:11",
	}

	mock.ExpectQuery("INSERT").
		WithArgs(info.SentToId, info.SentToId, info.Content, info.Date).
		WillReturnError(sql.ErrConnDone)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)
	logger.InitLogger()

	_, err = chatRepoMock.SaveMessageFromRestaurant(ctx, info)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}
