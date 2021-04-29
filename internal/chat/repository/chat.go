package repository

import (
	"context"
	"database/sql"

	"github.com/borscht/backend/internal/chat"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
)

type chatRepo struct {
	DB *sql.DB
}

func NewChattRepo(db *sql.DB) chat.ChatRepo {
	return &chatRepo{
		DB: db,
	}
}

func (ch chatRepo) GetAllChatsUser(ctx context.Context, uid int) ([]models.BriefInfoChat, error) {
	query :=
		`
		SELECT r.rid, r.name, r.avatar, m1.content
		FROM messages m1 LEFT 
		JOIN messages m2
		ON (m1.sentToRestaurant = m2.sentToRestaurant AND m1.mid < m2.mid) 
		JOIN restaurants r 
		ON (m1.sentToRestaurant = r.rid)
		WHERE m2.mid IS NULL AND m1.sentFromUser = $1 order by m1.mid DESC;
	`

	return ch.getAllChats(ctx, query, uid)
}

func (ch chatRepo) GetAllChatsRestaurant(ctx context.Context, rid int) ([]models.BriefInfoChat, error) {
	query :=
		`
		SELECT u.uid, u.name, u.photo, m1.content
		FROM messages m1 LEFT 
		JOIN messages m2
		ON (m1.sentToUser = m2.sentToUser AND m1.mid < m2.mid) 
		JOIN users u 
		ON (m1.sentToUser = u.uid)
		WHERE m2.mid IS NULL AND m1.sentFromRestaurant = $1 order by m1.mid DESC;
	`

	return ch.getAllChats(ctx, query, rid)
}

func (ch chatRepo) getAllChats(ctx context.Context, query string, id int) (
	[]models.BriefInfoChat, error) {

	messagesDB, err := ch.DB.Query(query, id)
	if err == sql.ErrNoRows {
		return []models.BriefInfoChat{}, nil
	}
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return nil, failError
	}

	messages := make([]models.BriefInfoChat, 0)
	for messagesDB.Next() {
		var message models.BriefInfoChat
		messagesDB.Scan(
			&message.Uid,
			&message.Name,
			&message.Avatar,
			&message.LastMessage,
		)

		logger.RepoLevel().InlineDebugLog(ctx, message)
		messages = append(messages, message)
	}

	return messages, nil
}

func (ch chatRepo) GetAllMessagesFromUser(ctx context.Context, uid, rid int) (
	[]models.InfoMessage, error) {

	query :=
		`
		SELECT mid, sentWhen, content FROM messages
		where (sentFromUser = $1 AND sentToRestaurant = $2)                                                
		order by mid DESC;
	`

	return ch.getAllMessages(ctx, query, uid, rid)
}

func (ch chatRepo) GetAllMessagesFromRestaurant(ctx context.Context, rid, uid int) (
	[]models.InfoMessage, error) {

	query :=
		`
		SELECT mid, sentWhen, content FROM messages
		where (sentFromRestaurant = $1 AND sentToUser = $2)                                                
		order by mid DESC;
	`

	return ch.getAllMessages(ctx, query, rid, uid)
}

func (ch chatRepo) getAllMessages(ctx context.Context, query string, id1, id2 int) (
	[]models.InfoMessage, error) {

	messagesDB, err := ch.DB.Query(query, id1, id2)
	if err == sql.ErrNoRows {
		return []models.InfoMessage{}, nil
	}
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return nil, failError
	}

	messages := make([]models.InfoMessage, 0)
	for messagesDB.Next() {
		var message models.InfoMessage
		messagesDB.Scan(
			&message.Id,
			&message.Date,
			&message.Text,
		)

		logger.RepoLevel().InlineDebugLog(ctx, message)
		messages = append(messages, message)
	}

	return messages, nil
}

func (ch chatRepo) GetUser(ctx context.Context, uid int) (*models.InfoOpponent, error) {
	query :=
		`
		SELECT uid, name, photo
		FROM users 
		WHERE (uid = $1)
	`
	return ch.getInfoOpponent(ctx, query, uid)
}

func (ch chatRepo) GetRestaurant(ctx context.Context, rid int) (*models.InfoOpponent, error) {
	query :=
		`
		SELECT rid, name, avatar
		FROM restaurants 
		WHERE (rid = $1)
	`
	return ch.getInfoOpponent(ctx, query, rid)
}

func (ch chatRepo) getInfoOpponent(ctx context.Context, query string, id int) (
	*models.InfoOpponent, error) {

	opponent := new(models.InfoOpponent)
	err := ch.DB.QueryRow(query, id).
		Scan(&opponent.Uid, &opponent.Name, &opponent.Avatar)

	if err == sql.ErrNoRows {
		return &models.InfoOpponent{}, nil
	}
	if err != nil {
		custError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, custError)
		return nil, custError
	}

	return opponent, nil
}
