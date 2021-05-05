package repository

import (
	"context"
	"database/sql"

	protoChat "github.com/borscht/backend/services/proto/chat"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
)

type ChatRepo interface {
	GetAllChatsUser(ctx context.Context, uid int) ([]*protoChat.BriefInfoChat, error)
	GetAllChatsRestaurant(ctx context.Context, rid int) ([]*protoChat.BriefInfoChat, error)
	GetAllMessagesFromUser(ctx context.Context, uid, rid int) ([]*protoChat.InfoMessage, error)
	GetAllMessagesFromRestaurant(ctx context.Context, rid, uid int) ([]*protoChat.InfoMessage, error)
}

type chatRepo struct {
	DB *sql.DB
}

func NewChatRepository(db *sql.DB) ChatRepo {
	return &chatRepo{
		DB: db,
	}
}

func (ch chatRepo) GetAllChatsUser(ctx context.Context, uid int) ([]*protoChat.BriefInfoChat, error) {
	query :=
		`
		SELECT res.rid, res.name, res.avatar, res.content FROM (
			SELECT m1.mid, r.rid, r.name, r.avatar, m1.content
			FROM messages m1 LEFT 
			JOIN messages m2
			ON (m1.sentToRestaurant = m2.sentToRestaurant AND m1.mid < m2.mid) 
			JOIN restaurants r 
			ON (m1.sentToRestaurant = r.rid)
			WHERE m2.mid IS NULL AND m1.sentFromUser = $1
	
			UNION
	
			SELECT m1.mid, r.rid, r.name, r.avatar, m1.content
			FROM messages m1 LEFT 
			JOIN messages m2
			ON (m1.sentFromRestaurant = m2.sentFromRestaurant AND m1.mid < m2.mid) 
			JOIN restaurants r 
			ON (m1.sentFromRestaurant = r.rid)
			WHERE m2.mid IS NULL AND m1.sentToUser = $1 
			
			ORDER BY mid DESC
		) AS res;
	`

	return ch.getAllChats(ctx, query, uid)
}

func (ch chatRepo) GetAllChatsRestaurant(ctx context.Context, rid int) ([]*protoChat.BriefInfoChat, error) {
	query :=
		`
		SELECT res.uid, res.name, res.photo, res.content FROM (
			SELECT m1.mid, u.uid, u.name, u.photo, m1.content
			FROM messages m1 LEFT 
			JOIN messages m2
			ON (m1.sentToUser = m2.sentToUser AND m1.mid < m2.mid) 
			JOIN users u 
			ON (m1.sentToUser = u.uid)
			WHERE m2.mid IS NULL AND m1.sentFromRestaurant = $1

			UNION

			SELECT m1.mid, u.uid, u.name, u.photo, m1.content
			FROM messages m1 LEFT 
			JOIN messages m2
			ON (m1.sentFromUser = m2.sentFromUser AND m1.mid < m2.mid) 
			JOIN users u 
			ON (m1.sentFromUser = u.uid)
			WHERE m2.mid IS NULL AND m1.sentToRestaurant = $1 
			
			ORDER BY mid DESC
		) AS res;
	`

	return ch.getAllChats(ctx, query, rid)
}

func (ch chatRepo) getAllChats(ctx context.Context, query string, id int) (
	[]*protoChat.BriefInfoChat, error) {

	messagesDB, err := ch.DB.Query(query, id)
	if err == sql.ErrNoRows {
		return []*protoChat.BriefInfoChat{}, nil
	}
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return nil, failError
	}

	messages := make([]*protoChat.BriefInfoChat, 0)
	for messagesDB.Next() {
		message := new(protoChat.BriefInfoChat)
		Info := new(protoChat.InfoOpponent)
		message.Info = Info
		messagesDB.Scan(
			&Info.Uid,
			&Info.Name,
			&Info.Avatar,
			&message.LastMessage,
		)

		logger.RepoLevel().InlineDebugLog(ctx, message)
		messages = append(messages, message)
	}

	return messages, nil
}

func (ch chatRepo) GetAllMessagesFromUser(ctx context.Context, uid, rid int) (
	[]*protoChat.InfoMessage, error) {

	query :=
		`
		SELECT mid, sentWhen, content FROM messages
		where (sentFromUser = $1 AND sentToRestaurant = $2)                                                
		order by mid DESC;
	`

	return ch.getAllMessages(ctx, query, uid, rid)
}

func (ch chatRepo) GetAllMessagesFromRestaurant(ctx context.Context, rid, uid int) (
	[]*protoChat.InfoMessage, error) {

	query :=
		`
		SELECT mid, sentWhen, content FROM messages
		where (sentFromRestaurant = $1 AND sentToUser = $2)                                                
		order by mid DESC;
	`

	return ch.getAllMessages(ctx, query, rid, uid)
}

func (ch chatRepo) getAllMessages(ctx context.Context, query string, id1, id2 int) (
	[]*protoChat.InfoMessage, error) {

	messagesDB, err := ch.DB.Query(query, id1, id2)
	if err == sql.ErrNoRows {
		return []*protoChat.InfoMessage{}, nil
	}
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return nil, failError
	}

	messages := make([]*protoChat.InfoMessage, 0)
	for messagesDB.Next() {
		message := new(protoChat.InfoMessage)
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
