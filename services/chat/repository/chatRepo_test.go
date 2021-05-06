package repository

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/borscht/backend/config"
	"github.com/borscht/backend/services/chat/models"
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

func TestChatRepo_GetAllChats(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	defer db.Close()
	chatRepoMock := NewChatRepository(db)

	rowsChats := sqlmock.NewRows([]string{"mid", "id", "role", "text", "date"})
	rowsChats.AddRow(1, 1, config.RoleUser, "hi", "01.01.21")

	mock.ExpectQuery("SELECT").
		WithArgs(1, config.RoleUser).
		WillReturnRows(rowsChats)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)
	logger.InitLogger()

	infoChat, err := chatRepoMock.GetAllChats(ctx, models.User{Id: 1, Role: config.RoleUser})
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	require.EqualValues(t, infoChat[0].Message.Text, "hi")
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
		WithArgs(1, config.RoleUser).
		WillReturnError(sql.ErrNoRows)

	ctx := new(context.Context)
	_, err = chatRepoMock.GetAllChats(*ctx, models.User{Id: 1, Role: config.RoleUser})
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestChatRepo_GetAllChats_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	defer db.Close()
	chatRepoMock := NewChatRepository(db)

	mock.ExpectQuery("SELECT").
		WithArgs(1, config.RoleUser).
		WillReturnError(sql.ErrConnDone)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)
	logger.InitLogger()

	_, err = chatRepoMock.GetAllChats(ctx, models.User{Id: 1, Role: config.RoleUser})
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestChatRepo_GetAllMessages(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	defer db.Close()
	chatRepoMock := NewChatRepository(db)

	rowsMessage := sqlmock.NewRows([]string{"mid", "senderId", "senderRole", "recipientId", "recipientRole", "content", "sentWhen"})
	rowsMessage.AddRow(1, 1, config.RoleUser, 1, config.RoleAdmin, "hi", "21.01.2020")

	mock.ExpectQuery("SELECT").
		WithArgs(1, config.RoleUser, 1, config.RoleAdmin).
		WillReturnRows(rowsMessage)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)
	logger.InitLogger()

	user1 := models.User{Id: 1, Role: config.RoleUser}
	user2 := models.User{Id: 1, Role: config.RoleAdmin}
	messages, err := chatRepoMock.GetAllMessages(ctx, user1, user2)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	require.EqualValues(t, messages[0].Message.Text, "hi")
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
		WithArgs(1, config.RoleUser, 1, config.RoleAdmin).
		WillReturnError(sql.ErrNoRows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)
	logger.InitLogger()

	user1 := models.User{Id: 1, Role: config.RoleUser}
	user2 := models.User{Id: 1, Role: config.RoleAdmin}
	_, err = chatRepoMock.GetAllMessages(ctx, user1, user2)
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
		WithArgs(1, config.RoleUser, 1, config.RoleAdmin).
		WillReturnError(sql.ErrConnDone)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)
	logger.InitLogger()

	user1 := models.User{Id: 1, Role: config.RoleUser}
	user2 := models.User{Id: 1, Role: config.RoleAdmin}
	_, err = chatRepoMock.GetAllMessages(ctx, user1, user2)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestChatRepo_SaveMessage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	defer db.Close()
	chatRepoMock := NewChatRepository(db)

	rowsMessage := sqlmock.NewRows([]string{"mid"})
	rowsMessage.AddRow(1)

	messageInfo := models.Chat{
		Message: models.Message{
			Text: "ho",
			Date: "01.01.21",
		},
		Sender: models.User{
			Id:   1,
			Role: config.RoleUser,
		},
		Recipient: models.User{
			Id:   1,
			Role: config.RoleAdmin,
		},
	}

	mock.ExpectQuery("INSERT").
		WithArgs(1, config.RoleUser, 1, config.RoleAdmin, messageInfo.Message.Text, messageInfo.Message.Date).
		WillReturnRows(rowsMessage)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)
	logger.InitLogger()

	mid, err := chatRepoMock.SaveMessage(ctx, messageInfo)
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
