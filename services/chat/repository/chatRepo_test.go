package repository

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/borscht/backend/utils/logger"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewChatRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	defer db.Close()
	chatRepoMock := NewChatRepository(db)
	if chatRepoMock != nil {
		return
	}
}

func TestChatRepo_GetAllChatsUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	defer db.Close()
	chatRepoMock := NewChatRepository(db)

	rowsChats := sqlmock.NewRows([]string{"id", "name", "photo", "content"})
	rowsChats.AddRow(1, "Dasha", "image.jpg", "hi")

	mock.ExpectQuery("SELECT").
		WithArgs(1).
		WillReturnRows(rowsChats)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)
	logger.InitLogger()

	infoChat, err := chatRepoMock.GetAllChatsUser(ctx, 1)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	require.EqualValues(t, infoChat[0].Info.Uid, 1)
}

func TestChatRepo_GetAllChats_UserNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	defer db.Close()
	chatRepoMock := NewChatRepository(db)

	mock.ExpectQuery("SELECT").
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	ctx := new(context.Context)
	_, err = chatRepoMock.GetAllChatsUser(*ctx, 1)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestChatRepo_GetAllChatsUser_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	defer db.Close()
	chatRepoMock := NewChatRepository(db)

	mock.ExpectQuery("SELECT").
		WithArgs(1).
		WillReturnError(sql.ErrConnDone)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)
	logger.InitLogger()

	_, err = chatRepoMock.GetAllChatsUser(ctx, 1)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestChatRepo_GetAllChatsRestaurant(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	defer db.Close()
	chatRepoMock := NewChatRepository(db)

	rowsChats := sqlmock.NewRows([]string{"id", "name", "photo", "content"})
	rowsChats.AddRow(1, "Dasha", "image.jpg", "hi")

	mock.ExpectQuery("SELECT").
		WithArgs(1).
		WillReturnRows(rowsChats)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)
	logger.InitLogger()

	infoChat, err := chatRepoMock.GetAllChatsRestaurant(ctx, 1)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	require.EqualValues(t, infoChat[0].Info.Uid, 1)
}

func TestChatRepo_GetAllMessagesFromUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	defer db.Close()
	chatRepoMock := NewChatRepository(db)

	rowsMessage := sqlmock.NewRows([]string{"id", "sentWhen", "content"})
	rowsMessage.AddRow(1, "21.01.2020", "hi")

	mock.ExpectQuery("SELECT").
		WithArgs(1, 1).
		WillReturnRows(rowsMessage)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)
	logger.InitLogger()

	messages, err := chatRepoMock.GetAllMessagesFromUser(ctx, 1, 1)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	require.EqualValues(t, messages[0].Id, 1)
}

func TestChatRepo_GetAllMessages_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	defer db.Close()
	chatRepoMock := NewChatRepository(db)

	mock.ExpectQuery("SELECT").
		WithArgs(1, 1).
		WillReturnError(sql.ErrNoRows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)
	logger.InitLogger()

	_, err = chatRepoMock.GetAllMessagesFromUser(ctx, 1, 1)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestChatRepo_GetAllMessages_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	defer db.Close()
	chatRepoMock := NewChatRepository(db)

	mock.ExpectQuery("SELECT").
		WithArgs(1, 1).
		WillReturnError(sql.ErrConnDone)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)
	logger.InitLogger()

	_, err = chatRepoMock.GetAllMessagesFromUser(ctx, 1, 1)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestChatRepo_GetAllMessagesFromRestaurant(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	defer db.Close()
	chatRepoMock := NewChatRepository(db)

	rowsMessage := sqlmock.NewRows([]string{"id", "sentWhen", "content"})
	rowsMessage.AddRow(1, "21.01.2020", "hi")

	mock.ExpectQuery("SELECT").
		WithArgs(1, 1).
		WillReturnRows(rowsMessage)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)
	logger.InitLogger()

	messages, err := chatRepoMock.GetAllMessagesFromRestaurant(ctx, 1, 1)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	require.EqualValues(t, messages[0].Id, 1)
}
