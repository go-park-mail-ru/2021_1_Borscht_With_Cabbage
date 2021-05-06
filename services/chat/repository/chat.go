package repository

import (
	"context"
	"database/sql"

	"github.com/borscht/backend/services/chat/models"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
)

type ChatRepo interface {
	GetAllChats(ctx context.Context, user models.User) ([]models.ChatInfo, error)
	GetAllMessages(ctx context.Context, user1, user2 models.User) ([]models.Chat, error)
	SaveMessage(ctx context.Context, messageInfo models.Chat) (int32, error)
}

type chatRepo struct {
	DB *sql.DB
}

func NewChatRepository(db *sql.DB) ChatRepo {
	return &chatRepo{
		DB: db,
	}
}

func (ch chatRepo) GetAllChats(ctx context.Context, user models.User) (
	[]models.ChatInfo, error) {

	query :=
		`
		SELECT DISTINCT ON(res.recipientId, res.recipientRole) 
		res.mid, res.recipientId, res.recipientRole, res.content, res.sentWhen
			FROM ((SELECT m1.mid, m1.recipientId, m1.recipientRole, m1.content, m1.sentWhen
			FROM messages m1
			LEFT JOIN messages m2
			ON (m1.recipientId = m2.recipientId AND m1.recipientRole = m2.recipientRole AND m1.mid < m2.mid)
			WHERE m2.mid IS NULL AND m1.senderId = $1 AND m1.senderRole = $2)
				UNION
				(SELECT m3.mid, m3.senderId, m3.senderRole, m3.content, m3.sentWhen
			FROM messages m3
			LEFT JOIN messages m4
			ON (m3.senderId = m4.senderId AND m3.senderRole = m4.senderRole AND m3.mid < m4.mid)
			WHERE m4.mid IS NULL AND m3.recipientId = $1 AND m3.recipientRole = $2) 
			ORDER BY mid DESC) 
			AS res;
	`

	return ch.getAllChats(ctx, query, user)
}

func (ch chatRepo) getAllChats(ctx context.Context, query string, user models.User) (
	[]models.ChatInfo, error) {

	messagesDB, err := ch.DB.Query(query, user.Id, user.Role)
	if err == sql.ErrNoRows {
		return []models.ChatInfo{}, nil
	}
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return nil, failError
	}

	chats := make([]models.ChatInfo, 0)
	for messagesDB.Next() {
		var chat models.ChatInfo
		var user models.User
		var message models.Message

		messagesDB.Scan(
			&message.Mid,
			&user.Id,
			&user.Role,
			&message.Text,
			&message.Date,
		)
		chat.User = user
		chat.Message = message

		logger.RepoLevel().DebugLog(ctx, logger.Fields{"message": chat})
		chats = append(chats, chat)
	}

	return chats, nil
}

func (ch chatRepo) GetAllMessages(ctx context.Context, user1, user2 models.User) (
	[]models.Chat, error) {

	logger.RepoLevel().DebugLog(ctx, logger.Fields{"user1": user1, "user2": user2})
	query :=
		`
		SELECT mid, senderId, senderRole, recipientId, recipientRole, content, sentWhen
		FROM messages
		WHERE (senderId = $1 AND senderRole = $2 AND recipientId = $3 AND recipientRole = $4)
			OR (senderId = $3 AND senderRole = $4 AND recipientId = $1 AND recipientRole = $2)                             
		ORDER BY mid DESC;
	`

	messagesDB, err := ch.DB.Query(query, user1.Id, user1.Role, user2.Id, user2.Role)
	if err == sql.ErrNoRows {
		return []models.Chat{}, nil
	}
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return nil, failError
	}

	chats := make([]models.Chat, 0)
	for messagesDB.Next() {
		chat := new(models.Chat)
		message := new(models.Message)
		sender := new(models.User)
		recipient := new(models.User)

		messagesDB.Scan(
			&message.Mid,
			&sender.Id,
			&sender.Role,
			&recipient.Id,
			&recipient.Role,
			&message.Text,
			&message.Date,
		)

		chat.Message = *message
		chat.Recipient = *recipient
		chat.Sender = *sender

		logger.RepoLevel().InlineDebugLog(ctx, message)
		chats = append(chats, *chat)
	}

	return chats, nil
}

func (ch chatRepo) SaveMessage(ctx context.Context, messageInfo models.Chat) (
	int32, error) {

	logger.RepoLevel().InfoLog(ctx, logger.Fields{"save message": messageInfo})
	query :=
		`
		INSERT INTO messages (senderId, senderRole, recipientId, recipientRole, content, sentWhen) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING mid
	`

	var mid int32
	err := ch.DB.QueryRow(query, messageInfo.Sender.Id, messageInfo.Sender.Role,
		messageInfo.Recipient.Id, messageInfo.Recipient.Role,
		messageInfo.Message.Text, messageInfo.Message.Date).Scan(&mid)

	if err != nil {
		insertError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, insertError)
		return 0, insertError
	}

	logger.RepoLevel().InfoLog(ctx, logger.Fields{"mid": mid})
	return mid, nil
}
