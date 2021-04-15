package repository

import (
	"context"
	"database/sql"

	"github.com/borscht/backend/config"
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
		return errors.BadRequestError("User with this email already exists")
	}

	err = u.DB.QueryRow("SELECT uid FROM users WHERE phone = $1", number).Scan(&userInDB)
	if err != sql.ErrNoRows && userInDB != currentUserId {
		logger.RepoLevel().InlineInfoLog(ctx, "User with this number already exists")
		return errors.BadRequestError("User with this number already exists")
	}

	return nil
}

func (u userRepo) Create(ctx context.Context, newUser models.User) (int, error) {
	err := u.checkUserWithThisData(ctx, newUser.Phone, newUser.Email, 0)
	if err != nil {
		return 0, err
	}

	var uid int
	err = u.DB.QueryRow("insert into users (name, phone, email, password, photo) values ($1, $2, $3, $4, $5) returning uid",
		newUser.Name, newUser.Phone, newUser.Email, newUser.HashPassword, config.DefaultUserImage).Scan(&uid)
	if err != nil {
		insertError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, insertError)
		return 0, insertError
	}

	return uid, nil
}

func (u userRepo) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	user := new(models.User)

	err := u.DB.QueryRow("select uid, name, phone, email, photo, password from users where (phone=$1 or email=$1)",
		login).Scan(&user.Uid, &user.Name, &user.Phone, &user.Email, &user.Avatar, &user.HashPassword)
	if err == sql.ErrNoRows {
		return nil, errors.AuthorizationError("user not found")
	}
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return nil, failError
	}

	logger.RepoLevel().InlineDebugLog(ctx, user.HashPassword)

	return user, nil
}

func (u userRepo) GetByUid(ctx context.Context, uid int) (models.User, error) {
	DBuser, err := u.DB.Query("select name, phone, email, photo from users where uid=$1", uid)
	if err != nil {
		return models.User{}, errors.AuthorizationError("user not found")
	}
	user := new(models.User)
	for DBuser.Next() {
		err = DBuser.Scan(
			&user.Name,
			&user.Phone,
			&user.Email,
			&user.Avatar,
		)
		if err != nil {
			failError := errors.FailServerError(err.Error())
			logger.RepoLevel().ErrorLog(ctx, failError)
			return models.User{}, failError
		}
	}
	return *user, nil
}

func (u userRepo) UpdateData(ctx context.Context, user models.UserData) error {
	err := u.checkUserWithThisData(ctx, user.Phone, user.Email, user.ID)

	if err != nil {
		return err
	}

	_, err = u.DB.Exec("UPDATE users SET phone = $1, email = $2, name = $3 where uid = $4",
		user.Phone, user.Email, user.Name, user.ID)
	if err != nil {
		return errors.AuthorizationError("curUser not found")
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
