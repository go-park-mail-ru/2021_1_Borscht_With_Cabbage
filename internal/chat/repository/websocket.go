package repository

import (
	"context"
	"database/sql"

	"github.com/borscht/backend/internal/chat"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
)

type websocketRepo struct {
	DB *sql.DB
}

func NewWebsocketRepo(db *sql.DB) chat.WebSocketRepo {
	return &websocketRepo{
		DB: db,
	}
}

func (w websocketRepo) SaveMessageFromUser(ctx context.Context, info models.WsMessageForRepo) (
	mid int, err error) {

	query :=
		`
	INSERT INTO messages (sentFromUser, sentToRestaurant, content, sentWhen) 
	VALUES ($1, $2, $3, $4) 
	RETURNING mid
	`

	err = w.DB.QueryRow(query, info.SentFromId, info.SentToId,
		info.Content, info.Date).Scan(&mid)

	if err != nil {
		insertError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, insertError)
		return 0, insertError
	}

	return mid, nil
}

func (w websocketRepo) SaveMessageFromRestaurant(ctx context.Context, info models.WsMessageForRepo) (
	mid int, err error) {

	query :=
		`
	INSERT INTO messages (sentFromRestaurant, sentToUser, content, sentWhen) 
	VALUES ($1, $2, $3, $4) 
	RETURNING mid
	`

	err = w.DB.QueryRow(query, info.SentFromId, info.SentToId,
		info.Content, info.Date).Scan(&mid)

	if err != nil {
		insertError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, insertError)
		return 0, insertError
	}

	return mid, nil
}
