package repository

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/user"
	_errors "github.com/borscht/backend/utils"

	"io"
	"mime/multipart"
	"os"
)

type userRepo struct {
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) user.UserRepo {
	return &userRepo{
		DB: db,
	}
}

func (u *userRepo) checkExistingUser(ctx context.Context, email, number string) error {
	var userInDB int
	err := u.DB.QueryRow("select uid from users where email = $1", email).Scan(&userInDB)
	if err != sql.ErrNoRows {
		return _errors.BadRequest(ctx, "User with this email already exists")
	}

	err = u.DB.QueryRow("SELECT uid FROM users WHERE phone = $1", number).Scan(&userInDB)
	if err != sql.ErrNoRows {
		return _errors.BadRequest(ctx, "User with this number already exists")
	}

	return nil
}

func (u *userRepo) checkUserWithThisData(ctx context.Context, email, number string, currentUserId int) error {
	var userInDB int
	err := u.DB.QueryRow("select uid from users where email = $1", email).Scan(&userInDB)
	if err != sql.ErrNoRows && userInDB != currentUserId {
		return _errors.BadRequest(ctx, "User with this email already exists")
	}

	err = u.DB.QueryRow("SELECT uid FROM users WHERE phone = $1", number).Scan(&userInDB)
	if err != sql.ErrNoRows && userInDB != currentUserId {
		return _errors.BadRequest(ctx, "User with this number already exists")
	}

	return nil
}

func (u *userRepo) Create(ctx context.Context, newUser models.User) (int, error) {
	err := u.checkExistingUser(ctx, newUser.Email, newUser.Phone)
	if err != nil {
		return 0, _errors.FailServer(ctx, err.Error())
	}

	var uid int
	err = u.DB.QueryRow("insert into users (name, phone, email, password, photo) values ($1, $2, $3, $4, $5) returning uid",
		newUser.Name, newUser.Phone, newUser.Email, newUser.Password, config.DefaultAvatar).Scan(&uid)
	if err != nil {
		return 0, _errors.FailServer(ctx, err.Error())
	}

	return uid, nil
}

func (u *userRepo) CheckUserExists(ctx context.Context, userToCheck models.UserAuth) (models.User, error) {
	user := new(models.User)

	err := u.DB.QueryRow("select uid, name, photo from users where (phone=$1 or email=$1) and password=$2",
		userToCheck.Login, userToCheck.Password).Scan(&user.Uid, &user.Name, &user.Avatar)
	if err == sql.ErrNoRows {

		return models.User{}, _errors.NewCustomError(ctx, http.StatusBadRequest, "user not found")
	}
	if err != nil {
		return models.User{}, _errors.FailServer(ctx, err.Error())
	}

	return *user, nil
}

func (u *userRepo) GetByUid(ctx context.Context, uid int) (models.User, error) {
	DBuser, err := u.DB.Query("select name, phone, email, photo from users where uid=$1", uid)
	if err != nil {
		return models.User{}, _errors.Authorization(ctx, "user not found")
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
			return models.User{}, _errors.FailServer(ctx, err.Error())
		}
	}
	return *user, nil
}

func (u *userRepo) Update(ctx context.Context, newUser models.UserData, uid int) error {
	err := u.checkUserWithThisData(ctx, newUser.Phone, newUser.Email, uid)
	if err != nil {
		return _errors.FailServer(ctx, err.Error())
	}

	_, err = u.DB.Exec("UPDATE users SET phone = $1, email = $2, name = $3, photo = $4 where uid = $5",
		newUser.Phone, newUser.Email, newUser.Name, newUser.Avatar, uid)
	if err != nil {
		return _errors.Authorization(ctx, "curUser not found")
	}

	return nil
}

func (u *userRepo) UploadAvatar(ctx context.Context, image *multipart.FileHeader, filename string) error {
	// Читаем файл из пришедшего запроса
	src, err := image.Open()
	if err != nil {
		return _errors.FailServer(ctx, err.Error())
	}
	defer src.Close()

	// создаем файл у себя
	dst, err := os.Create(filename)
	if err != nil {
		return _errors.FailServer(ctx, err.Error())
	}
	defer dst.Close()

	// копируем один в другой
	if _, err = io.Copy(dst, src); err != nil {
		return _errors.FailServer(ctx, err.Error())
	}

	return nil
}
