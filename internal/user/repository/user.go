package repository

import (
	"database/sql"
	"fmt"
	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/user"
	_errors "github.com/borscht/backend/utils"
	"net/http"

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

func (u *userRepo) checkExistingUser(email, number string, currentUserId int32) error {
	var userInDB int32
	err := u.DB.QueryRow("select uid from users where email = $1", email).Scan(&userInDB)
	if err != sql.ErrNoRows && err != nil && userInDB != currentUserId {
		return _errors.NewCustomError(http.StatusBadRequest, "User with this email already exists")
	}

	err = u.DB.QueryRow("SELECT uid FROM users WHERE phone = $1", number).Scan(&userInDB)
	if err != sql.ErrNoRows && err != nil && userInDB != currentUserId {
		return _errors.NewCustomError(http.StatusBadRequest, "User with this number already exists")
	}

	return nil
}

func (u *userRepo) Create(newUser models.User) (int32, error) {
	err := u.checkExistingUser(newUser.Email, newUser.Phone, -1)
	if err != nil {
		return 0, _errors.FailServer(err.Error())
	}
	fmt.Println(newUser)

	var uid int32
	err = u.DB.QueryRow("insert into users (name, phone, email, password, photo) values ($1, $2, $3, $4, $5) returning uid",
		newUser.Name, newUser.Phone, newUser.Email, newUser.Password, config.DefaultAvatar).Scan(&uid)
	if err != nil {
		return 0, _errors.FailServer(err.Error())
	}

	return uid, nil
}

func (u *userRepo) CheckUserExists(userToCheck models.UserAuth) (models.User, error) {
	user := new(models.User)

	err := u.DB.QueryRow("select uid, name, photo from users where (phone=$1 or email=$1) and password=$2",
		userToCheck.Login, userToCheck.Password).Scan(&user.Uid, &user.Name, &user.Avatar)
	if err == sql.ErrNoRows {

		return models.User{}, _errors.NewCustomError(http.StatusBadRequest, "user not found")
	}
	if err != nil {
		return models.User{}, _errors.FailServer(err.Error())
	}

	return *user, nil
}

func (u *userRepo) GetByUid(uid int32) (models.User, error) {
	DBuser, err := u.DB.Query("select name, phone, email, photo from users where uid=$1", uid)
	if err != nil {
		return models.User{}, _errors.Authorization("user not found")
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
			return models.User{}, _errors.FailServer(err.Error())
		}
	}
	return *user, nil
}

func (u *userRepo) Update(newUser models.UserData, uid int32) error {
	err := u.checkExistingUser(newUser.Phone, newUser.Email, uid)
	if err != nil {
		return _errors.FailServer(err.Error())
	}

	_, err = u.DB.Exec("UPDATE users SET phone = $1, email = $2, name = $3, photo = $4 where uid = $5",
		newUser.Phone, newUser.Email, newUser.Name, newUser.Avatar, uid)
	if err != nil {
		return _errors.Authorization("curUser not found")
	}

	return nil
}

func (u *userRepo) UploadAvatar(image *multipart.FileHeader, filename string) error {
	// Читаем файл из пришедшего запроса
	src, err := image.Open()
	if err != nil {
		return _errors.FailServer(err.Error())
	}
	defer src.Close()

	// создаем файл у себя
	dst, err := os.Create(filename)
	if err != nil {
		return _errors.FailServer(err.Error())
	}
	defer dst.Close()

	// копируем один в другой
	if _, err = io.Copy(dst, src); err != nil {
		return _errors.FailServer(err.Error())
	}

	return nil
}
