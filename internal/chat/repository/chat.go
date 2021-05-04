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

func (ch chatRepo) SaveMessageFromUser(ctx context.Context, info models.WsMessageForRepo) (
	mid int, err error) {

	query :=
		`
	INSERT INTO messages (sentFromUser, sentToRestaurant, content, sentWhen) 
	VALUES ($1, $2, $3, $4) 
	RETURNING mid
	`

	err = ch.DB.QueryRow(query, info.SentFromId, info.SentToId,
		info.Content, info.Date).Scan(&mid)

	if err != nil {
		insertError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, insertError)
		return 0, insertError
	}

	return mid, nil
}

func (ch chatRepo) SaveMessageFromRestaurant(ctx context.Context, info models.WsMessageForRepo) (
	mid int, err error) {

	query :=
		`
	INSERT INTO messages (sentFromRestaurant, sentToUser, content, sentWhen) 
	VALUES ($1, $2, $3, $4) 
	RETURNING mid
	`

	err = ch.DB.QueryRow(query, info.SentFromId, info.SentToId,
		info.Content, info.Date).Scan(&mid)

	if err != nil {
		insertError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, insertError)
		return 0, insertError
	}

	return mid, nil
}
