package repository

import (
	"context"
	"database/sql"

	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/user"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
)

type userRepo struct {
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) user.UserRepo {
	return &userRepo{
		DB: db,
	}
}

func (u userRepo) checkUserWithThisData(ctx context.Context, number, email string, currentUserId int) error {
	var userInDB int

	err := u.DB.QueryRow("select uid from users where email = $1", email).Scan(&userInDB)
	if err != sql.ErrNoRows && userInDB != currentUserId {
		logger.RepoLevel().InlineInfoLog(ctx, "User with this email already exists")
		return errors.NewErrorWithMessage("User with this email already exists")
	}

	err = u.DB.QueryRow("SELECT uid FROM users WHERE phone = $1", number).Scan(&userInDB)
	if err != sql.ErrNoRows && userInDB != currentUserId {
		logger.RepoLevel().InlineInfoLog(ctx, "User with this number already exists")
		return errors.NewErrorWithMessage("User with this number already exists")
	}

	return nil
}

func (u userRepo) UpdateData(ctx context.Context, user models.UserData) error {
	err := u.checkUserWithThisData(ctx, user.Phone, user.Email, user.ID)

	if err != nil {
		return err
	}

	_, err = u.DB.Exec("UPDATE users SET phone = $1, email = $2, name = $3 where uid = $4",
		user.Phone, user.Email, user.Name, user.ID)
	if err != nil {
		return errors.FailServerError("curUser not found")
	}

	return nil
}

func (u userRepo) UpdateAvatar(ctx context.Context, idUser int, filename string) error {
	_, err := u.DB.Exec("UPDATE users SET photo = $1 where uid = $2",
		filename, idUser)
	if err != nil {
		dbError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, dbError)
		return dbError
	}

	return nil
}
