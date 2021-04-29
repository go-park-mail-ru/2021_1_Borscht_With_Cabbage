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
		ON (m1.senttorestaurant = m2.senttorestaurant AND m1.mid < m2.mid) 
		JOIN restaurants r 
		ON (m1.senttorestaurant = r.rid)
		WHERE m2.mid IS NULL AND m1.sentfromuser = $1 order by m1.mid DESC;
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
