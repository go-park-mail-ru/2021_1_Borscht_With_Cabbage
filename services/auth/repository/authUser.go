package authRepo

import (
	"context"
	"database/sql"

	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/services/auth"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
)

type authRepo struct {
	DB *sql.DB
}

func NewUserAuthRepo(db *sql.DB) auth.UserAuthRepo {
	return &authRepo{
		DB: db,
	}
}

func (a authRepo) checkUserWithThisData(ctx context.Context, number, email string, currentUserId int) error {
	var userInDB int

	err := a.DB.QueryRow("select uid from users where email = $1", email).Scan(&userInDB)
	if err != sql.ErrNoRows && userInDB != currentUserId {
		logger.RepoLevel().InlineInfoLog(ctx, "User with this email already exists")
		return errors.NewErrorWithMessage("User with this email already exists")
	}

	err = a.DB.QueryRow("select uid from users where phone = $1", number).Scan(&userInDB)
	if err != sql.ErrNoRows && userInDB != currentUserId {
		logger.RepoLevel().InlineInfoLog(ctx, "User with this number already exists")
		return errors.NewErrorWithMessage("User with this number already exists")
	}

	return nil
}

func (a authRepo) Create(ctx context.Context, newUser models.User) (int, error) {
	err := a.checkUserWithThisData(ctx, newUser.Phone, newUser.Email, 0)
	if err != nil {
		return 0, err
	}

	var uid int
	err = a.DB.QueryRow("insert into users (name, phone, email, password, photo) values ($1, $2, $3, $4, $5) returning uid",
		newUser.Name, newUser.Phone, newUser.Email, newUser.HashPassword, newUser.Avatar).Scan(&uid)
	if err != nil {
		insertError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, insertError)
		return 0, insertError
	}

	return uid, nil
}

func (a authRepo) GetByUid(ctx context.Context, uid int) (models.User, error) {
	DBuser, err := a.DB.Query("select name, phone, email, photo from users where uid=$1", uid)
	if err != nil {
		return models.User{}, errors.NewErrorWithMessage("user not authorization").SetDescription("user not found")
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

func (a authRepo) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	user := new(models.User)

	err := a.DB.QueryRow("select uid, name, phone, email, photo, password from users where (phone=$1 or email=$1)",
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

func (a authRepo) GetAddress(ctx context.Context, uid int) (*models.Address, error) {
	queri := `SELECT name, latitude, longitude FROM addresses WHERE uid = $1`

	logger.RepoLevel().InlineDebugLog(ctx, uid)
	var address models.Address
	err := a.DB.QueryRow(queri, uid).Scan(&address.Name, &address.Latitude, &address.Longitude)
	if err == sql.ErrNoRows {
		logger.RepoLevel().InlineDebugLog(ctx, "end get address not address")
		return &models.Address{}, nil
	}
	if err != nil {
		err := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, err)
		return nil, err
	}

	logger.RepoLevel().InlineDebugLog(ctx, "end get address")
	return &address, nil
}
