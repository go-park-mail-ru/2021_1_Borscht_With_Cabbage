package repository

import (
	"context"
	"database/sql"
	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/user"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"

	"io"
	"mime/multipart"
	"os"
)

type UserRepo struct {
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) user.UserRepo {
	return &UserRepo{
		DB: db,
	}
}

func (u *UserRepo) checkExistingUser(ctx context.Context, email, number string) error {
	var userInDB int
	err := u.DB.QueryRow("select uid from users where email = $1", email).Scan(&userInDB)
	if err != sql.ErrNoRows {
		logger.RepoLevel().InlineInfoLog(ctx, "User with this email already exists")
		return errors.BadRequestError("User with this email already exists")
	}

	err = u.DB.QueryRow("select uid from users where phone = $1", number).Scan(&userInDB)
	if err != sql.ErrNoRows {
		logger.RepoLevel().InlineInfoLog(ctx, "User with this number already exists")
		return errors.BadRequestError("User with this number already exists")
	}

	return nil
}

func (u *UserRepo) checkUserWithThisData(ctx context.Context, email, number string, currentUserId int) error {
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

func (u *UserRepo) Create(ctx context.Context, newUser models.User) (int, error) {
	err := u.checkExistingUser(ctx, newUser.Email, newUser.Phone)
	if err != nil {
		return 0, err
	}

	var uid int
	err = u.DB.QueryRow("insert into users (name, phone, email, password, photo) values ($1, $2, $3, $4, $5) returning uid",
		newUser.Name, newUser.Phone, newUser.Email, newUser.Password, config.DefaultAvatar).Scan(&uid)
	if err != nil {
		insertError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, insertError)
		return 0, insertError
	}

	return uid, nil
}

func (u *UserRepo) CheckUserExists(ctx context.Context, userToCheck models.UserAuth) (*models.User, error) {
	user := new(models.User)

	err := u.DB.QueryRow("select uid, name, phone, email, password, photo from users where (phone=$1 or email=$1) and password=$2",
		userToCheck.Login, userToCheck.Password).
		Scan(&user.Uid, &user.Name, &user.Phone, &user.Email, &user.Password, &user.Avatar)
	if err == sql.ErrNoRows {
		return nil, errors.AuthorizationError("user not found")
	}
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return nil, failError
	}

	return user, nil
}

func (u *UserRepo) GetByUid(ctx context.Context, uid int) (models.User, error) {
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

func (u *UserRepo) Update(ctx context.Context, newUser models.UserData, uid int) error {
	err := u.checkUserWithThisData(ctx, newUser.Phone, newUser.Email, uid)
	if err != nil {
		return err
	}

	_, err = u.DB.Exec("UPDATE users SET phone = $1, email = $2, name = $3, photo = $4 where uid = $5",
		newUser.Phone, newUser.Email, newUser.Name, newUser.Avatar, uid)
	if err != nil {
		return errors.AuthorizationError("curUser not found")
	}

	return nil
}

func (u *UserRepo) UploadAvatar(ctx context.Context, image *multipart.FileHeader, filename string) error {
	// Читаем файл из пришедшего запроса
	src, err := image.Open()
	if err != nil {
		custError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, custError)
		return custError
	}
	defer src.Close()

	// создаем файл у себя
	dst, err := os.Create(filename)
	if err != nil {
		custError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, custError)
		return custError
	}
	defer dst.Close()

	// копируем один в другой
	if _, err = io.Copy(dst, src); err != nil {
		custError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, custError)
		return custError
	}

	return nil
}
